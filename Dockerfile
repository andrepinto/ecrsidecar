FROM alpine

RUN apk -U add ca-certificates

COPY ecrsidecar /bin/ecrsidecar

VOLUME /opt/config/ecrsidecar

ENTRYPOINT ["/bin/ecrsidecar"]