FROM golang:1.14
WORKDIR /opt/relayer
ADD . .
RUN go build -v -o build/artemis-relay main.go


FROM parity/subkey:2.0.0

COPY --from=0 /opt/relayer/build/artemis-relay /usr/local/bin/

COPY contracts /tmp/contracts
COPY relay-config.toml /tmp

USER root
RUN apt-get update \
    && apt-get install -y --no-install-recommends apt-utils \
    && apt-get install -y ca-certificates
USER 1001
# RUN sed -i 's/.\/contracts\//\/tmp\/contracts\//g' /tmp/relay-config.toml
# RUN cat /tmp/relay-config.toml

ENV ARTEMIS_SUBSTRATE_KEY=//Relay

ENTRYPOINT ["/usr/local/bin/artemis-relay"]
