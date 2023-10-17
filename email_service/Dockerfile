FROM golang:1.21.1-alpine AS builder

WORKDIR /email_service/go

ADD app/  /email_service/go

RUN  go clean --modcache && go build -mod=readonly -o app/ cmd/server/app.go

FROM scratch

COPY --from=builder  /email_service/go/app bin/
COPY --from=builder  /email_service/go/configs configs/
COPY --from=builder  /email_service/go/cert cert/
COPY --from=builder  /email_service/go/templates templates/

EXPOSE 8080

CMD ["bin/app"]