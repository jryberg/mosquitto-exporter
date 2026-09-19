Mosquitto Exporter
=============================
Prometheus exporter for the [Mosquitto MQTT message broker](https://mosquitto.org/).

There is a docker image available:
```
docker run \
  -p 9234:9234 jryberg/mosquitto-exporter:v0.7.7 \
  --endpoint tcp://mosquitto:1883
```

Usage:
======

```
NAME:
   mosquitto-exporter - Prometheus exporter for broker metrics

USAGE:
   mosquitto-exporter [global options] command [command options] [arguments...]

VERSION:
   v0.7.7 (7b3248b0889448cc0f026604e537d84b1abdb09e), go1.24.6

AUTHORS:
   Johan Ryberg <johan@securit.se>
   Arturo Reuschenbach Puncernau <a.reuschenbach.puncernau@sap.com>
   Fabian Ruff <fabian.ruff@sap.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --endpoint value, -e value      Endpoint for the Mosquitto message broker (default: "tcp://127.0.0.1:1883") [$BROKER_ENDPOINT]
   --bind-address value, -b value  Listen address for metrics HTTP endpoint (default: "0.0.0.0:9234") [$BIND_ADDRESS]
   --user value, -u value          Username for the Mosquitto message broker [$MQTT_USER]
   --pass value, -p value          Password for the User on the Mosquitto message broker [$MQTT_PASS]
   --cert value, -c value          Location of a TLS certificate .pem file for the Mosquitto message broker [$MQTT_CERT]
   --key value, -k value           Location of a TLS private key .pem file for the Mosquitto message broker [$MQTT_KEY]
   --ca value, -a value            Location of a CA certificate .pem file used to verify the Mosquitto message broker [$MQTT_CA]
   --insecure-skip-verify          Skip verification of the Mosquitto message broker TLS certificate (insecure) [$MQTT_INSECURE_SKIP_VERIFY]
   --client-id value, -i value     Client id to be used to connect to the Mosquitto message broker [$MQTT_CLIENT_ID]
   --reset-metrics, -r             Reset metrics when loosing connection to broker (default: true) [$RESET_METRICS]
   --help, -h                      show help
   --version, -v                   print the version
```

Grafana dashboard:
======

[`grafana/mosquitto-dashboard.json`](grafana/mosquitto-dashboard.json) is a ready-made dashboard for these metrics.

To use it, open **Dashboards → New → Import** in Grafana and upload the file, or put it in a
[provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards) folder.
It picks the Prometheus data source from a dashboard variable, so it needs no editing.

The exporter's metrics have no labels of their own, so a broker is picked by its scrape labels:

- **Job** and **Instance** select the scrape target. When several brokers are selected, the panels show their total.
- **Filters** takes any other label your Prometheus adds, for example `kubernetes_namespace = mosquitto`.

The rows are:

- **Overview:** exporter state, uptime, clients, subscriptions, message rates, retained messages and dropped publishes.
- **Clients and subscriptions:** connected and persistent sessions, new connections per second (shows reconnect loops), subscriptions and expired sessions.
- **Throughput:** messages, bytes and PUBLISH traffic with sent above zero and received below, dropped messages, average payload size and payload share of traffic.
- **Broker load averages:** the broker's own 1, 5 and 15 minute averages.
- **Message store and memory:** heap, stored and retained messages, and the outgoing packet queue.
- **Exporter process:** the exporter's own memory, CPU and scrape health (collapsed).

With `--reset-metrics` (the default) the exporter sets every metric to 0 when it loses the broker.
The dashboard marks that period with an annotation. Expect a one-off spike in the rate panels when
the connection comes back, because the counters jump from 0 to their old totals.
