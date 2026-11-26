FROM       debian:trixie-slim@sha256:8aa5a7bf5a1955bddbe7d74fcdd0b91ab30e74362a8b73af672b0481b35b44a0
LABEL      source_repository="https://github.com/jryberg/mosquitto-exporter"
ARG        TARGETPLATFORM
COPY       $TARGETPLATFORM/mosquitto-exporter /mosquitto_exporter
RUN        apt-get update && \
           apt-get install -y --no-install-recommends --no-install-suggests ca-certificates && \
           rm -rf /var/lib/apt/lists/*
EXPOSE     9234
ENTRYPOINT ["/mosquitto_exporter"]
