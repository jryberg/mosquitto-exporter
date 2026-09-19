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
