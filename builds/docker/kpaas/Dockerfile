FROM alpine:latest

COPY builds/release/deploy /app/deploy
COPY builds/release/service /app/service
COPY run/docker/entrypoint.sh /app/entrypoint.sh

EXPOSE 8080 8081

ENV LOG_LEVEL info
ENV SERVICE_ID 0

WORKDIR /app
ENTRYPOINT ["/app/entrypoint.sh"]
