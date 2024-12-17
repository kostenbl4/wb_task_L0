FROM golang:1.23.4-alpine AS builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install

COPY . .

RUN go build -tags musl --ldflags "-extldflags -static" -o main ./cmd/api

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine AS runner
COPY .env .
COPY index.html .

ADD ./cmd/migrations/*.sql /migrations/
ADD entrypoint.sh /migrations/

COPY --from=builder /go/src/main /
COPY --from=builder /go/bin/goose /goose

RUN chmod +x /migrations/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "/migrations/entrypoint.sh && /main"]