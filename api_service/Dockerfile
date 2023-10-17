FROM golang:1.21.1-alpine3.18 AS builder

WORKDIR /api_service/go

ADD app/  /api_service/go
ADD app/configs /api_service/configs
ADD app/cert /api_service/cert

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/server/app.go

FROM alpine:3.18

COPY --from=builder  /api_service/go/app service/
COPY --from=builder  /api_service/configs configs/
COPY --from=builder  /api_service/cert cert/

EXPOSE 8080

CMD ["service/app"]