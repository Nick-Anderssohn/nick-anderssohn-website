# Image to build dart with
FROM google/dart:2.0 as dart-builder

# install ruby and its dependencies
RUN apt-get update
RUN apt-get install --no-install-recommends -y ruby-full
RUN apt-get install --no-install-recommends -y build-essential bison openssl libreadline6\
 libreadline6-dev curl git-core zlib1g zlib1g-dev libssl-dev\
 libyaml-dev libxml2-dev autoconf libc6-dev ncurses-dev automake\
 libtool
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

# install sass
RUN gem install sass --no-user-install

# install webdev
RUN pub global activate webdev
ENV PATH="${PATH}:/root/.pub-cache/bin"

WORKDIR /
RUN mkdir builder
WORKDIR /builder
COPY web .
RUN rm -rf build || true
RUN rm -rf .dart_tool || true
RUN rm .packages || true
RUN rm pubspec.lock || true
RUN pub get
RUN webdev build

# Image to build go with
FROM golang:1.11.2 as go-builder
RUN mkdir -p /go/src/nick-anderssohn-website/ops-gateway
WORKDIR /go/src/nick-anderssohn-website/ops-gateway
RUN mkdir -p /go/bin
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY . .
RUN /go/bin/dep ensure
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .

# Actual image
FROM alpine:3.8
RUN mkdir resources
COPY --from=go-builder /go/src/nick-anderssohn-website/ops-gateway/ops-gateway .
COPY admin.json .
COPY db_config.json .
COPY --from=dart-builder /builder/build /resources
CMD ["/ops-gateway"]