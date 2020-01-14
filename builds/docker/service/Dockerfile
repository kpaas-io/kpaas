FROM alpine:latest

COPY builds/release/service /app/service

COPY charts /app/charts

EXPOSE 8080

CMD ["/app/service"]
