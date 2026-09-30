# Multi-stage Dockerfile for Production
# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instalar certificados SSL de CA para conexiones HTTPS seguras (LinkedIn & Gemini APIs)
RUN apk add --no-cache ca-certificates

# Copiar archivos de modulos Go
COPY go.mod go.sum ./
RUN go mod download

# Copiar el codigo fuente
COPY . .

# Compilar binario estatico optimizado sin CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/git-to-feed ./cmd/api

# Final Stage: Imagen minimalista para ejecucion
FROM alpine:3.19

WORKDIR /app

# Copiar ca-certificates del builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar binario ejecutable
COPY --from=builder /app/git-to-feed /app/git-to-feed

# Copiar la carpeta de migraciones SQL
COPY --from=builder /app/migrations /app/migrations

# Exponer el puerto por defecto de la aplicacion
EXPOSE 8080

# Usuario sin privilegios por seguridad
RUN adduser -D appuser
USER appuser

# Comando de ejecucion
ENTRYPOINT ["/app/git-to-feed"]