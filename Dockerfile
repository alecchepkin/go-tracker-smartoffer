FROM golang:alpine as builder

RUN apk add -q --update \
   && apk add -q \
           bash \
           git \
           curl \
&& rm -rf /var/cache/apk/*

RUN mkdir /build
ADD . /build/src/tracker-smartoffer
WORKDIR /build/src/tracker-smartoffer
ENV GOPATH /build
RUN go get && go build -o tracker-smartoffer .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /build/src/tracker-smartoffer/tracker-smartoffer /app/

WORKDIR /app

CMD ["./tracker-smartoffer"]