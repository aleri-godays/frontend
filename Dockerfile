#backend build
FROM golang:1.13-stretch as gobuilder

#run as unpriviledged user
RUN addgroup --gid 990 app && adduser --disabled-password --uid 994 --gid 990 --gecos '' app

RUN mkdir -p /build
WORKDIR /build

#use the go modules proxy for faster downloads
ENV GOPROXY=https://proxy.golang.org

#copy go.mod and go.sum and download all modules to cache this docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download

#copy the rest and build the binary
COPY . .
RUN make build-linux


#production image
FROM alpine:3.10
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN addgroup -g 990 app && \
    adduser -D -u 990 -G app app

COPY --chown=994:990 --from=gobuilder /build/release/linux-amd64/frontend /app

USER 994:990
ENTRYPOINT ["./app"]
CMD ["run"]
