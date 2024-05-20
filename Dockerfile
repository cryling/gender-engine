# Stage 1: Build the CSV to SQLite converter
FROM golang:alpine3.19 as builder
ENV CGO_ENABLED=1
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev


FROM builder as csv-to-sqlite-builder
WORKDIR /app/csv-to-sqlite
COPY csv-to-sqlite/ .
ARG CSV_FILE_PATH=data/wgnd_2_0_name-gender-code.csv
COPY $CSV_FILE_PATH /app/csv-to-sqlite/data/data.csv
RUN go mod tidy
RUN go build -o /csv-to-sqlite

# Run the CSV to SQLite converter
RUN /csv-to-sqlite

# Stage 2: Build the Gin app
FROM builder as gin-builder
WORKDIR /app/api
COPY api/ .
ENV CGO_ENABLED=1
RUN go mod tidy
RUN go build -o /api

# Stage 3: Create the final image
FROM alpine:latest
WORKDIR /root/
COPY --from=gin-builder /api .
COPY --from=csv-to-sqlite-builder /app/csv-to-sqlite/data/data.db ./data.db

RUN apk add --no-cache libc6-compat

ENV GIN_MODE=release
ENV RATE_LIMIT_ENABLED=true
ENV RATE_LIMIT=50
ENV RATE_BURST=500

EXPOSE 8080
CMD ["./api"]
