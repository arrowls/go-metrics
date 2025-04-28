FROM alpine:3.18

RUN apk --no-cache update && \
    apk --no-cache upgrade && \
    apk add --no-cache bash curl

RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz \
    && tar -xzf migrate.linux-amd64.tar.gz -C /bin \
    && rm migrate.linux-amd64.tar.gz

RUN chmod +x /bin/migrate

WORKDIR /app

COPY ./docker/migrations/migrations.sh /app/migrations.sh
COPY ./internal/database/migrations /app/migrations

RUN chmod +x /app/migrations.sh

ENTRYPOINT ["bash", "/app/migrations.sh"]
