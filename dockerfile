# ETAPA 1: Construcción
FROM golang:1.21-alpine AS builder

# En Linux/Docker usamos rutas simples
WORKDIR /app

# Copiamos el go.mod y el código fuente
COPY . .

# Compilamos el binario
RUN go build -o monitor .

# ETAPA 2: Ejecución
FROM alpine:latest
WORKDIR /root/

# Copiamos solo lo necesario desde la etapa de 'builder'
COPY --from=builder /app/monitor .
COPY --from=builder /app/config.json .

# Comando para ejecutar
CMD ["./monitor"]