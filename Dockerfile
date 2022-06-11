FROM golang:alpine
RUN apk add --no-cache git
WORKDIR /trueid-shorten-link
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download
EXPOSE 8080
CMD go run cmd/shorten-link/main.go