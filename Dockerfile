# Estágio de construção
FROM golang:1.23.2-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos do projeto
COPY . .

# Baixa as dependências
RUN go mod download

# Compila o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o exilium-blog-backend ./cmd/api

# Estágio final
FROM alpine:latest

# Instala ca-certificates para HTTPS e postgresql para pg_isready
RUN apk --no-cache add ca-certificates postgresql

# Define o diretório de trabalho
WORKDIR /root/

# Copia o binário compilado
COPY --from=builder /app/exilium-blog-backend .

# Expõe a porta do servidor
EXPOSE 8080

# Comando para rodar o servidor
CMD ["./exilium-blog-backend"]