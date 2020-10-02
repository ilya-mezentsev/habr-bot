FROM golang:1.14.6 AS build

WORKDIR /var/www/apps/habr-bot/
COPY ./main .

CMD ./main -mode tg
