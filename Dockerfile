FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
  git \
  golang-go \
  nodejs nodejs-legacy npm

ENV GOPATH=/go
ADD . /go/src/github.com/SpComb/qmsk-dmx

WORKDIR /go/src/github.com/SpComb/qmsk-dmx/web
RUN npm install
RUN ./node_modules/typescript/bin/tsc

WORKDIR /go/src/github.com/SpComb/qmsk-dmx
RUN go get -d ./cmd/...
RUN go install -v ./cmd/...

ENV ARTNET_DISCOVERY=2.255.255.255
CMD ["/go/bin/qmsk-dmx", \
  "--log=info", \
  "--http-listen=:8000", "--http-static=web/", \
  "--heads-library=library", "config/" \
]

EXPOSE 8000
