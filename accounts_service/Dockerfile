FROM golang:1.21.3-alpine AS builder

WORKDIR /accounts_service/go

ADD app/  /accounts_service/go

RUN  go clean --modcache && go build -mod=readonly -o app/ cmd/server/app.go

FROM scratch

COPY --from=builder  /accounts_service/go/app bin/
COPY --from=builder  /accounts_service/go/cert cert/

EXPOSE 8080:8080

CMD ["bin/app"]