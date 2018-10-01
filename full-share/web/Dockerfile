# Run with a volume to build dart code
FROM google/dart:2.0

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