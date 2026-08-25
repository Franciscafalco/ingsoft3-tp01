# ---- Etapa 1: build (imagen con el SDK de Go) ----
FROM golang:1.25-alpine AS build
WORKDIR /src

# go.mod y go.sum primero: capa cacheable, solo se rehace si cambian las dependencias.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server .

# ---- Etapa 2: runtime (imagen mínima, sin el SDK) ----
FROM alpine:3.20 AS final
LABEL org.opencontainers.image.source=https://github.com/Franciscafalco/ingsoft3-tp01
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
