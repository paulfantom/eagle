FROM gcr.io/distroless/static

ARG BINARY=eagle-amd64
COPY _output/$BINARY /usr/local/bin/eagle

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/eagle"]
