FROM --platform=linux/arm64 arm64v8/golang:latest as builder

SHELL ["/bin/bash", "-c"]

# Устанавливаем значение переменной GOARCH внутри Docker контейнера
ENV GOARCH=arm64

WORKDIR /go/server

# Устанавливаем git
RUN apt-get update
RUN apt-get install -y git
RUN git clone https://github.com/PomoschDev/API-Gateway.git .

# Компилируем проект
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o apigateway ./cmd/entrypoint

# Создаем финальный образ
FROM arm64v8/alpine:latest

# Рабочая директория
WORKDIR /app

#Порт для прослушки
ENV PORT=443

# Копируем исполняемый файл из предыдущего образа
COPY --from=builder /go/server/entrypoint ./apigateway

# Устанавливаем время
RUN apk add tzdata && echo "Europe/Moscow" > /etc/timezone && ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
#

# Копируем файл конфигурации в контейнер
COPY ./prod.yaml .

# Открываем порты
EXPOSE ${PORT}