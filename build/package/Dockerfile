FROM golang:1.12-alpine as builder

RUN apk add --update --no-cache git gcc musl-dev make
RUN go get github.com/Pixep/bitbadger/cmd/bitbadger

WORKDIR /go/src/github.com/Pixep/bitbadger
RUN make tests

WORKDIR /go/src/github.com/Pixep/bitbadger/cmd/bitbadger
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /opt/bin/bitbadger .

FROM centurylink/ca-certs

COPY --from=builder /opt/bin/bitbadger /opt/bin/bitbadger

ENTRYPOINT ["/opt/bin/bitbadger"]