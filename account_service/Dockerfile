FROM golang:1.21.1-alpine AS builder

WORKDIR /account_service/go

ADD app/  /account_service/go

RUN  go clean --modcache && go build -mod=readonly -o app/ cmd/server/app.go

FROM scratch

COPY --from=builder  /account_service/go/app bin/
COPY --from=builder  /account_service/go/configs configs/
COPY --from=builder  /account_service/go/cert cert/

EXPOSE 8080:8080

CMD ["bin/app"]