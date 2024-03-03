FROM golang:1.20.7 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o auth

FROM scratch
COPY --from=build /app/auth /
ENTRYPOINT ["/auth", "serve"]
