FROM       debian:bookworm-slim
LABEL      source_repository="https://github.com/jryberg/mosquitto-exporter"
COPY       mosquitto-exporter /mosquitto_exporter
RUN        apt update && \
             apt install -y ca-certificates
EXPOSE     9234
ENTRYPOINT ["/mosquitto_exporter"]
