# Stage 1: Build the CSV to SQLite converter
FROM golang:1.24-alpine3.21 AS builder
ENV CGO_ENABLED=1
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev


FROM builder AS csv-to-sqlite-builder
RUN apk add --no-cache curl
WORKDIR /app/csv-to-sqlite
COPY csv-to-sqlite/ .
RUN mkdir -p data && \
    curl -L -o data/wgnd_2_0_name-gender-code.csv \
      "https://dataverse.harvard.edu/api/access/datafile/4750348?format=original" && \
    curl -L -o data/wgnd_2_0_name-gender.csv \
      "https://dataverse.harvard.edu/api/access/datafile/4750351?format=original"
RUN go mod tidy
RUN go build -o /csv-to-sqlite

# Run the CSV to SQLite converter
RUN /csv-to-sqlite

# Stage 2: Build the API
FROM builder AS api-builder
WORKDIR /app/api
COPY api/ .
ENV CGO_ENABLED=1
RUN go mod tidy
RUN go build -o /api

# Stage 3: Create the final image
FROM alpine:3.21
LABEL org.opencontainers.image.source=https://github.com/cryling/gender-engine
WORKDIR /root/
COPY --from=api-builder /api .
COPY --from=csv-to-sqlite-builder /app/csv-to-sqlite/data/data.db ./data.db

RUN apk add --no-cache libc6-compat

ENV RATE_LIMIT_ENABLED=true
ENV RATE_LIMIT=50
ENV RATE_BURST=500

EXPOSE 8080
CMD ["./api"]
