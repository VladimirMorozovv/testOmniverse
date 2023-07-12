FROM golang:1.18 as builder
MAINTAINER Vladimir Morozov



WORKDIR /build

COPY ./go.mod ./go.sum ./

RUN go mod download
RUN go mod verify

COPY . .
RUN go build -v -o app cmd/apiserver/main.go

FROM golang:1.18

WORKDIR /application
RUN mkdir /configs

COPY --from=builder /build/app /application

EXPOSE 8080

CMD ./app \
    --config="/configs/config.yml" \
