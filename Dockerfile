FROM golang:1.14 as builder

ENV GOPATH "/go"
ENV GOOS "linux"
ENV CGO_ENABLED "0"
WORKDIR /go/src
COPY . .
RUN cd /go/src/cmd/acnh && go install

FROM busybox:latest
WORKDIR /go
COPY --from=builder /usr/share/zoneinfo/America/Los_Angeles /etc/localtime
COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=builder /go/bin/acnh .
COPY --from=builder /go/src/cmd/acnh/acnh.json .
COPY --from=builder /go/src/cmd/acnh/templates/ ./templates/
COPY --from=builder /go/src/cmd/acnh/css/ ./css/
COPY --from=builder /go/src/cmd/acnh/js/ ./js/
ENTRYPOINT ["/go/acnh"]