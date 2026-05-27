# --- Estágio de Compilação ---
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Instalação de certificados e pacotes de build essenciais
RUN apk add --no-cache git

# Copia dos arquivos do projeto
COPY go.mod main.go ./

# Baixa as dependências e gera o arquivo go.sum dinamicamente no container
RUN go mod tidy

# Compilação otimizada: CGO desabilitado, binário estático e remoção de símbolos de debug (-w -s)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/http-server-projeto-korp main.go

# --- Estágio de Execução ---
FROM alpine:latest

# Instalação de certificados CA para segurança de rede
RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copia do binário do estágio anterior
COPY --from=builder /app/http-server-projeto-korp /app/http-server-projeto-korp

# Permissões seguras para o executável
RUN chown -R appuser:appgroup /app

# Execução como usuário sem privilégios root
USER appuser

# Exposição da porta de comunicação interna
EXPOSE 8080

CMD ["/app/http-server-projeto-korp"]
