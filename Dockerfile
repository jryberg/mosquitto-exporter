FROM       debian:trixie-slim
LABEL      source_repository="https://github.com/jryberg/mosquitto-exporter"
COPY       mosquitto-exporter /mosquitto_exporter
RUN        apt-get update && \
           apt-get install -y --no-install-recommends --no-install-suggests ca-certificates && \
           rm -rf /var/lib/apt/lists/*
EXPOSE     9234
ENTRYPOINT ["/mosquitto_exporter"]
