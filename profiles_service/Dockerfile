FROM golang:1.21.4-alpine AS builder

WORKDIR /profiles_service/go

ADD app/  /profiles_service/go

RUN  go clean --modcache && go build -mod=readonly -o app/ cmd/server/app.go

FROM scratch

COPY --from=builder  /profiles_service/go/app bin/
COPY --from=builder  /profiles_service/go/cert cert/

EXPOSE 8080:8080

CMD ["bin/app"]