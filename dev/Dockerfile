FROM golang:1.25.5-bookworm
WORKDIR /app

RUN apt-get update && apt-get install -y make gcc g++ curl

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64 && \
    chmod +x tailwindcss-linux-arm64 && \
    mv tailwindcss-linux-arm64 /usr/local/bin/tailwindcss

COPY dev/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

COPY . .

EXPOSE 4000
ENTRYPOINT ["entrypoint.sh"]