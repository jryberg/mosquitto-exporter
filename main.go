package main

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"crypto/tls"
	"crypto/x509"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

const (
	appName = "mosquitto-exporter"
)

var (
	ignoreKeyMetrics = map[string]string{
		"$SYS/broker/timestamp":        "The timestamp at which this particular build of the broker was made. Static.",
		"$SYS/broker/version":          "The version of the broker. Static.",
		"$SYS/broker/clients/active":   "deprecated in favour of $SYS/broker/clients/connected",
		"$SYS/broker/clients/inactive": "deprecated in favour of $SYS/broker/clients/disconnected",
	}
	counterKeyMetrics = map[string]string{
		"$SYS/broker/bytes/received":            "The total number of bytes received since the broker started.",
		"$SYS/broker/bytes/sent":                "The total number of bytes sent since the broker started.",
		"$SYS/broker/messages/received":         "The total number of messages of any type received since the broker started.",
		"$SYS/broker/messages/sent":             "The total number of messages of any type sent since the broker started.",
		"$SYS/broker/publish/bytes/received":    "The total number of PUBLISH bytes received since the broker started.",
		"$SYS/broker/publish/bytes/sent":        "The total number of PUBLISH bytes sent since the broker started.",
		"$SYS/broker/publish/messages/received": "The total number of PUBLISH messages received since the broker started.",
		"$SYS/broker/publish/messages/sent":     "The total number of PUBLISH messages sent since the broker started.",
		"$SYS/broker/publish/messages/dropped":  "The total number of PUBLISH messages that have been dropped due to inflight/queuing limits.",
		"$SYS/broker/uptime":                    "The total number of seconds since the broker started.",
		"$SYS/broker/clients/maximum":           "The maximum number of clients connected simultaneously since the broker started",
		"$SYS/broker/clients/expired":           "The number of disconnected persistent clients that have been expired and removed through the persistent_client_expiration option.",
		"$SYS/broker/connections/socket/count":  "The total number of socket connections made to the broker since it started, whether or not the MQTT connection succeeded.",
	}
	counterMetrics = map[string]*MosquittoCounter{}
	gaugeMetrics   = map[string]prometheus.Gauge{}
)

func main() {
	cmd := &cli.Command{
		Name:    appName,
		Version: versionString(),
		Authors: []any{
			"Johan Ryberg <johan@securit.se>",
			"Arturo Reuschenbach Puncernau <a.reuschenbach.puncernau@sap.com>",
			"Fabian Ruff <fabian.ruff@sap.com>",
		},
		Usage:  "Prometheus exporter for broker metrics",
		Action: runServer,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "endpoint",
				Aliases: []string{"e"},
				Usage:   "Endpoint for the Mosquitto message broker",
				Sources: cli.EnvVars("BROKER_ENDPOINT"),
				Value:   "tcp://127.0.0.1:1883",
			},
			&cli.StringFlag{
				Name:    "bind-address",
				Aliases: []string{"b"},
				Usage:   "Listen address for metrics HTTP endpoint",
				Value:   "0.0.0.0:9234",
				Sources: cli.EnvVars("BIND_ADDRESS"),
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Username for the Mosquitto message broker",
				Value:   "",
				Sources: cli.EnvVars("MQTT_USER"),
			},
			&cli.StringFlag{
				Name:    "pass",
				Aliases: []string{"p"},
				Usage:   "Password for the User on the Mosquitto message broker",
				Value:   "",
				Sources: cli.EnvVars("MQTT_PASS"),
			},
			&cli.StringFlag{
				Name:    "cert",
				Aliases: []string{"c"},
				Usage:   "Location of a TLS certificate .pem file for the Mosquitto message broker",
				Value:   "",
				Sources: cli.EnvVars("MQTT_CERT"),
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "Location of a TLS private key .pem file for the Mosquitto message broker",
				Value:   "",
				Sources: cli.EnvVars("MQTT_KEY"),
			},
			&cli.StringFlag{
				Name:    "ca",
				Aliases: []string{"a"},
				Usage:   "Location of a CA certificate .pem file used to verify the Mosquitto message broker",
				Value:   "",
				Sources: cli.EnvVars("MQTT_CA"),
			},
			&cli.BoolFlag{
				Name:    "insecure-skip-verify",
				Usage:   "Skip verification of the Mosquitto message broker TLS certificate (insecure)",
				Value:   false,
				Sources: cli.EnvVars("MQTT_INSECURE_SKIP_VERIFY"),
			},
			&cli.StringFlag{
				Name:    "client-id",
				Aliases: []string{"i"},
				Usage:   "Client id to be used to connect to the Mosquitto message broker",
				Value:   "",
				Sources: cli.EnvVars("MQTT_CLIENT_ID"),
			},
			&cli.BoolFlag{
				Name:    "reset-metrics",
				Aliases: []string{"r"},
				Usage:   "Reset metrics when loosing connection to broker",
				Value:   true,
				Sources: cli.EnvVars("RESET_METRICS"),
			},
		},
	}

	cmd.Run(context.Background(), os.Args)
}

func resetMetrics() {
	for topic := range counterMetrics {
		if counterMetrics[topic] != nil {
			counterMetrics[topic].Set(0)
		}
	}
	for topic := range gaugeMetrics {
		if gaugeMetrics[topic] != nil {
			gaugeMetrics[topic].Set(0)
		}
	}
}

func runServer(ctx context.Context, cmd *cli.Command) error {
	log.Infof("Starting %s %s", appName, versionString())

	opts := mqtt.NewClientOptions()
	opts.SetCleanSession(true)
	opts.AddBroker(cmd.String("endpoint"))

	if cmd.String("client-id") != "" {
		opts.SetClientID(cmd.String("client-id"))
	}

	// if you have a username you'll need a password with it
	if cmd.String("user") != "" {
		opts.SetUsername(cmd.String("user"))
		if cmd.String("pass") != "" {
			opts.SetPassword(cmd.String("pass"))
		}
	}
	// if you have a client certificate you want a key aswell
	hasKeyPair := cmd.String("cert") != "" && cmd.String("key") != ""
	if !hasKeyPair && (cmd.String("cert") != "" || cmd.String("key") != "") {
		log.Println("Warning: For a client certificate both certificate and private key are needed. Skipping client certificate.")
	}
	if hasKeyPair || cmd.String("ca") != "" || cmd.Bool("insecure-skip-verify") {
		tlsConfig := &tls.Config{}
		if hasKeyPair {
			keyPair, err := tls.LoadX509KeyPair(cmd.String("cert"), cmd.String("key"))
			fatalfOnError(err, "Failed to load certificate/keypair: %s", err)
			tlsConfig.Certificates = []tls.Certificate{keyPair}
		}
		if cmd.String("ca") != "" {
			caPEM, err := os.ReadFile(cmd.String("ca"))
			fatalfOnError(err, "Failed to read CA certificate: %s", err)
			caPool := x509.NewCertPool()
			if !caPool.AppendCertsFromPEM(caPEM) {
				log.Fatalf("No certificates found in CA file %s", cmd.String("ca"))
			}
			tlsConfig.RootCAs = caPool
		}
		if cmd.Bool("insecure-skip-verify") {
			log.Warn("Warning: TLS certificate verification of the broker is disabled")
			tlsConfig.InsecureSkipVerify = true
		}
		opts.SetTLSConfig(tlsConfig)
		if !strings.HasPrefix(cmd.String("endpoint"), "ssl://") &&
			!strings.HasPrefix(cmd.String("endpoint"), "tls://") {
			log.Println("Warning: To use TLS the endpoint URL will have to begin with 'ssl://' or 'tls://'")
		}
	}

	opts.OnConnect = func(client mqtt.Client) {
		log.Infof("Connected to %s", cmd.String("endpoint"))
		// subscribe on every (re)connect
		token := client.Subscribe("$SYS/#", 0, func(_ mqtt.Client, msg mqtt.Message) {
			processUpdate(msg.Topic(), string(msg.Payload()))
		})
		if !token.WaitTimeout(10 * time.Second) {
			log.Println("Error: Timeout subscribing to topic $SYS/#")
		}
		if err := token.Error(); err != nil {
			log.Errorf("Failed to subscribe to topic $SYS/#: %s", err)
		}
		// paho does not report a rejected subscription (SUBACK 0x80) as an error
		if st, ok := token.(*mqtt.SubscribeToken); ok {
			for topic, code := range st.Result() {
				if code == 0x80 {
					log.Errorf("Subscription to topic %s was denied by the broker, check ACL for user %s", topic, cmd.String("user"))
				}
			}
		}
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		if cmd.Bool("reset-metrics") {
			log.Warnf("Error: Connection to %s lost: %s, resetting counters", cmd.String("endpoint"), err)
			resetMetrics()
		} else {
			log.Warnf("Error: Connection to %s lost: %s", cmd.String("endpoint"), err)
		}
	}
	client := mqtt.NewClient(opts)

	// try to connect forever
	for {
		token := client.Connect()
		if token.WaitTimeout(5 * time.Second) {
			if token.Error() == nil {
				break
			}
			log.Errorf("Error: Failed to connect to broker: %s", token.Error())
		} else {
			log.Errorf("Timeout connecting to endpoint %s", cmd.String("endpoint"))
		}
		time.Sleep(5 * time.Second)
	}
	log.Infof("Connected to %s", cmd.String("endpoint"))

	// init the router and server
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", serveVersion)
	log.Infof("Listening on %s...", cmd.String("bind-address"))
	err := http.ListenAndServe(cmd.String("bind-address"), nil)
	fatalfOnError(err, "Failed to bind on %s: ", cmd.String("bind-address"))
	return nil
}

// $SYS/broker/bytes/received
func processUpdate(topic, payload string) {
	//log.Printf("Got broker update with topic %s and data %s", topic, payload)
	if _, ok := ignoreKeyMetrics[topic]; !ok {
		if _, ok := counterKeyMetrics[topic]; ok {
			// log.Printf("Processing counter metric %s with data %s", topic, payload)
			processCounterMetric(topic, payload)
		} else {
			//log.Printf("Processing gauge metric %s with data %s", topic, payload)
			processGaugeMetric(topic, payload)
		}
	}
}

func processCounterMetric(topic, payload string) {
	if counterMetrics[topic] != nil {
		value := parseValue(payload)
		counterMetrics[topic].Set(value)
	} else {
		// create a mosquitto counter pointer
		mCounter := NewMosquittoCounter(prometheus.NewDesc(
			parseTopic(topic),
			topic,
			[]string{},
			prometheus.Labels{},
		))

		// save it
		counterMetrics[topic] = mCounter
		// register the metric
		prometheus.MustRegister(mCounter)
		// add the first value
		value := parseValue(payload)
		counterMetrics[topic].Set(value)
	}
}

func processGaugeMetric(topic, payload string) {
	if gaugeMetrics[topic] != nil {
		value := parseValue(payload)
		gaugeMetrics[topic].Set(value)
	} else {
		gaugeMetrics[topic] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: parseTopic(topic),
			Help: topic,
		})
		// register the metric
		prometheus.MustRegister(gaugeMetrics[topic])
		// add the first value
		value := parseValue(payload)
		gaugeMetrics[topic].Set(value)
	}
}

func parseTopic(topic string) string {
	name := strings.Replace(topic, "$SYS/", "", 1)
	name = strings.Replace(name, "/", "_", -1)
	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	name = strings.Replace(name, ".", "_", -1)
	return name
}

func parseValue(payload string) float64 {
	// fmt.Printf("Payload %s \n", payload)
	var validValue = regexp.MustCompile(`-?\d{1,}[.]\d{1,}|\d{1,}`)
	// get the first value of the string
	strArray := validValue.FindAllString(payload, 1)
	if len(strArray) > 0 {
		// parse to float
		value, err := strconv.ParseFloat(strArray[0], 64)
		if err == nil {
			return value
		}
	}
	return 0
}

func fatalfOnError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}
