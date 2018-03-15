FROM golang:1.9.4 as go-build

RUN curl -L -o /tmp/dep-linux-amd64 https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && install -m 0755 /tmp/dep-linux-amd64 /usr/local/bin/dep

WORKDIR /go/src/github.com/qmsk/dmx

COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . ./
RUN go install -v ./cmd/...




FROM node:9.8.0 as web-build

WORKDIR /go/src/github.com/qmsk/dmx/web

COPY web/package.json ./
RUN npm install

COPY web ./
RUN ./node_modules/typescript/bin/tsc


# must match with go-build base image
FROM debian:stretch

RUN mkdir -p /opt/qmsk-dmx /opt/qmsk-dmx/bin

COPY --from=go-build /go/bin/qmsk-dmx /opt/qmsk-dmx/bin/
COPY --from=web-build /go/src/github.com/qmsk/dmx/web/ /opt/qmsk-dmx/web
COPY library/ /opt/qmsk-dmx/library

WORKDIR /opt/qmsk-dmx
VOLUME /etc/qmsk-dmx
ENV ARTNET_DISCOVERY=2.255.255.255
CMD ["/opt/qmsk-dmx/bin/qmsk-dmx", \
  "--log=info", \
  "--http-listen=:8000", \
  "--http-static=/opt/qmsk-dmx/web/", \
  "--heads-library=/opt/qmsk-dmx/library/", \
  "/etc/qmsk-dmx" \
]

EXPOSE 8000
