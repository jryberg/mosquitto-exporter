FROM       debian:trixie-slim@sha256:a99cfc517144bc59b1978475ec53b46ecabec7e43635402ee5b77cc54cd1b20a
LABEL      source_repository="https://github.com/jryberg/mosquitto-exporter"
ARG        TARGETPLATFORM
COPY       $TARGETPLATFORM/mosquitto-exporter /mosquitto_exporter
RUN        apt-get update && \
           apt-get install -y --no-install-recommends --no-install-suggests ca-certificates && \
           rm -rf /var/lib/apt/lists/*
EXPOSE     9234
ENTRYPOINT ["/mosquitto_exporter"]
