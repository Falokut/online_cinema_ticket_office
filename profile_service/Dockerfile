FROM golang:1.21.1-alpine3.18 AS builder

WORKDIR /user_service/go/src/

ADD app/  /user_service/go/src/

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/server/app.go

FROM alpine:3.14

COPY --from=builder  /user_service/go/src/app /
COPY --from=builder  /user_service/go/src/config.yml /
COPY --from=builder  /user_service/go/src/.env /

CMD ["/app"]