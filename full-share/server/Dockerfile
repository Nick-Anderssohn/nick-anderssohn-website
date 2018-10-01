FROM golang:1.10.3 as builder
RUN go get github.com/lib/pq
RUN go get github.com/satori/go.uuid
RUN mkdir -p /go/src/nick-anderssohn-website/full-share/server
WORKDIR /go/src/nick-anderssohn-website/full-share/server
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .

FROM alpine:3.8
RUN mkdir download
COPY --from=builder /go/src/nick-anderssohn-website/full-share/server/server .
COPY --from=builder /go/src/nick-anderssohn-website/full-share/server/db_config.json .
COPY --from=builder /go/src/nick-anderssohn-website/full-share/server/config.json .
EXPOSE 8080
CMD ["/server"]