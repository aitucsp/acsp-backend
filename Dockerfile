FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

COPY ./ /github.com/aitucsp/acsp-backend
WORKDIR /github.com/aitucsp/acsp-backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go
# ./.bin/api ./cmd/api/main.go
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /github.com/aitucsp/acsp-backend/.bin/app .
COPY --from=builder /github.com/aitucsp/acsp-backend/.env ./.env

EXPOSE 8080

CMD ["./app"]