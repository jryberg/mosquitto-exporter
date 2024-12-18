Mosquitto Exporter
=============================
Prometheus exporter for the [Mosquitto MQTT message broker](https://mosquitto.org/).

There is a docker image available:
```
docker run \
  -p 9234:9234 jryberg/mosquitto-exporter:v0.7.4 \
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
   v0.7.4 (082cef5d97e73778b4452be24a8afbaf9065c703), go1.23.3

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
   --client-id value, -i value     Client id to be used to connect to the Mosquitto message broker [$MQTT_CLIENT_ID]
   --reset-metrics, -r             Reset metrics when loosing connection to broker (default: true) [$RESET_METRICS]
   --help, -h                      show help
   --version, -v                   print the version
```
