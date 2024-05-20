# Stage 1: Build the CSV to SQLite converter
FROM golang:1.22.2 as builder
WORKDIR /app/csv-to-sqlite
COPY csv-to-sqlite/ .
RUN go mod tidy
RUN go build -o /csv-to-sqlite

# Run the CSV to SQLite converter
RUN /csv-to-sqlite

# Stage 2: Build the Gin app
FROM golang:1.22.2 as gin-builder
WORKDIR /app/api
COPY api/ .
RUN go mod tidy
RUN go build -o /api

# Stage 3: Create the final image
FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=gin-builder /api .
COPY --from=builder /app/csv-to-sqlite/data/data.db ./data.db

ENV GIN_MODE=release
ENV RATE_LIMIT=50
ENV RATE_BURST=500

EXPOSE 8080
CMD ["./api"]
