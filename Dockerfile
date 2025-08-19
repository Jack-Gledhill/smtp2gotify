FROM golang:1.24-bookworm AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

FROM debian:12.7-slim AS runner

ARG version="dev"
ARG revision="dev"

ENV MODE="production"
ENV GIT_REF=$version

LABEL org.opencontainers.image.authors="Jack Gledhill"
LABEL org.opencontainers.image.description="An adapter written in Go that receives alerts over SMTP and relays them to a Gotify server."
LABEL org.opencontainers.image.documentation="https://github.com/Jack-Gledhill/smtp2gotify"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.revision=$revision
LABEL org.opencontainers.image.source="https://github.com/Jack-Gledhill/smtp2gotify"
LABEL org.opencontainers.image.title="SMTP2Gotify"
LABEL org.opencontainers.image.url="https://github.com/Jack-Gledhill/smtp2gotify"
LABEL org.opencontainers.image.version=$version

EXPOSE 25

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app
COPY --from=builder /app/ /app

ENTRYPOINT [ "./main" ]