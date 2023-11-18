FROM golang:1.21.4-alpine AS builder

WORKDIR /image_storage_service/go

ADD app/  /image_storage_service/go

RUN  go clean --modcache && go build -mod=readonly -o app/ cmd/server/app.go

FROM scratch

COPY --from=builder  /image_storage_service/go/app bin/

EXPOSE 8080:8080

CMD ["bin/app"]