FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/wateringhole ./cmd/wateringhole/
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /bin/wateringhole /usr/local/bin/wateringhole
ENV PORT="9000" DATA_DIR="/data"
EXPOSE 9000
HEALTHCHECK --interval=30s --timeout=5s CMD curl -sf http://localhost:9000/health || exit 1
ENTRYPOINT ["wateringhole"]
