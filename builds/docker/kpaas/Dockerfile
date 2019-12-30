FROM alpine:latest

COPY builds/release/* /app/
COPY builds/docker/dashboard.conf /app/config/dashboard.conf
COPY builds/docker/entrypoint.sh /app/entrypoint.sh

EXPOSE 80 8080 8081

ENV LOG_LEVEL info
ENV SERVICE_ID 0
ENV DASHBOARD_PORT 80

WORKDIR /app
ENTRYPOINT ["/app/entrypoint.sh"]
