FROM golang:alpine3.18 AS build

# Important:
# Because this is a CGO enabled package, you are required to set it as 1.
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /app/gender-engine

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/gender-engine .
COPY --from=build /app/data.db .

ENV GIN_MODE=release
ENV RATE_LIMIT=50
ENV RATE_BURST=500

EXPOSE 8080
CMD ["./gender-engine"]
