## Builder
#FROM golang:1.16-alpine AS builder
FROM harbor.klik.doctor/library/golang:1.16-alpine AS builder

ENV GO111MODULE=on

RUN apk update && apk upgrade && apk add --no-cache git gcc libc-dev make openssl

WORKDIR /app
COPY . .

RUN go get

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/rating-reviews-svc

## Distribution
#FROM alpine:latest
FROM harbor.klik.doctor/library/alpine:latest

ENV KD_ENV=prd
ENV TZ="Asia/Jakarta"

RUN apk update && apk upgrade && apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /go/bin /go/bin

EXPOSE 5600

RUN chmod +x /go/bin/rating-reviews-svc

CMD ["/go/bin/rating-reviews-svc"]
