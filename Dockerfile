FROM       debian:trixie-slim@sha256:d7e12182ce18b85b93007c1dedf31f2d29e01ccf3182cc4017c709b6259bc132
LABEL      source_repository="https://github.com/jryberg/mosquitto-exporter"
ARG        TARGETPLATFORM
COPY       $TARGETPLATFORM/mosquitto-exporter /mosquitto_exporter
RUN        apt-get update && \
           apt-get install -y --no-install-recommends --no-install-suggests ca-certificates && \
           rm -rf /var/lib/apt/lists/*
EXPOSE     9234
ENTRYPOINT ["/mosquitto_exporter"]
