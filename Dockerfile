ARG GOLANG_BUILDER_IMAGE
ARG DEBIAN_IMAGE

FROM $GOLANG_BUILDER_IMAGE as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app .

FROM $DEBIAN_IMAGE

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

WORKDIR /usr/src/app

COPY --from=builder /usr/local/bin/app /app

CMD ["/app"]
