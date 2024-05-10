FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src/app

COPY . ./

RUN go build -o ./bin/app ./cmd/app/main.go

FROM alpine

WORKDIR /app

RUN apk add make

COPY --from=builder /usr/local/src/app/bin/app .
COPY --from=builder /usr/local/src/app/cases ./cases
COPY --from=builder /usr/local/src/app/Makefile .

CMD make test_cases