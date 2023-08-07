FROM golang:1.20.7 as build

ENV APP_NAME auth

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -v -o $APP_NAME

FROM alpine:3.18

ENV APP_NAME auth

RUN mkdir /app

COPY --from=build /app/$APP_NAME /usr/local/bin/$APP_NAME

CMD $APP_NAME serve
