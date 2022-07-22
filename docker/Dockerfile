## Builder
FROM golang:1.16-alpine AS builder

ENV GO111MODULE=on

RUN apk update && apk upgrade && apk add --no-cache git gcc libc-dev make openssl

WORKDIR /app
COPY . .

RUN go get

RUN git clone https://github.com/go-swagger/go-swagger && cd go-swagger && go install ./cmd/swagger
RUN swagger generate spec -o swagger.yaml --scan-models
RUN rm -rf go-swagger

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/rating-reviews-svc

## Distribution
FROM alpine:latest

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
