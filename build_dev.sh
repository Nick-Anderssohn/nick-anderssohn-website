#! /bin/bash

cd full-share/web
#docker image build -t dart-builder .
#cwd=$(pwd)
#docker run --name full-share-dart-builder -v $cwd:/builder dart-builder /bin/bash -c "pub get; webdev build --no-release"
#docker rm full-share-dart-builder
pub get
webdev build --no-release