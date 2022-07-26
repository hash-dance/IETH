FROM  golang:1.13-alpine AS builder
WORKDIR /ieth
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux  go build -mod=vendor -o app main.go


FROM alpine:latest
#RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /ieth/conf/ssl ./conf/ssl
COPY --from=builder /ieth/app .

EXPOSE 80

CMD ["./app"]