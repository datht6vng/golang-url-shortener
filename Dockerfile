FROM golang:1.17-alpine as builder


ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /service/app/url-shortener

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/shorten-link/main ./cmd/shorten-link


FROM alpine:lasted

ENV CONFIG_PATH="./configs/config.toml"

WORKDIR /service/app/url-shortener

RUN  apk add --no-cache bash \
        \
        && apk add --no-cache tzdata \
        \
        && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
        \
        && echo "Asia/Ho_Chi_Minh" > /etc/timezone
RUN apk --no-cache add ca-certificates && apk --no-cache add curl

COPY --from=builder /service/app/url-shortener/cmd/shorten-link/main .
COPY --from=builder /service/app/url-shortener/view ./view
COPY --from=builder /service/app/url-shortener/configs ./configs

CMD ["./main"]

