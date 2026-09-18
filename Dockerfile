FROM       debian:trixie-slim@sha256:d7e12182ce18b85b93007c1dedf31f2d29e01ccf3182cc4017c709b6259bc132
LABEL      source_repository="https://github.com/jryberg/mosquitto-exporter"
COPY       mosquitto-exporter /mosquitto_exporter
RUN        apt update && \
           apt install -y ca-certificates
EXPOSE     9234
ENTRYPOINT ["/mosquitto_exporter"]
