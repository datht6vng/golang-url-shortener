FROM golang:alpine
RUN apk add --no-cache git
WORKDIR /golang-url-shortener
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# RUN go get .
RUN go build -o golang-url-shortener .
EXPOSE 8080
CMD ["./golang-url-shortener"]