FROM gcr.io/distroless/static

WORKDIR /

COPY autobot-linux-amd64 /main
COPY autobot-hangouts-linux-amd64 /autobot-hangouts

ENTRYPOINT ["/main"]