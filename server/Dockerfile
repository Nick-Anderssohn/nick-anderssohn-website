FROM golang:1.10.3 as builder
RUN go get github.com/lib/pq
RUN go get github.com/satori/go.uuid
RUN mkdir -p /go/src/full-share/server
WORKDIR /go/src/full-share/server
COPY . .
RUN go build

FROM alpine
RUN mkdir download
COPY --from=builder /go/src/full-share/server/server .
EXPOSE 8080
CMD ["/server"]