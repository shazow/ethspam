FROM golang:1.13-alpine AS builder
RUN apk update && apk add make
ADD . /ethspam
RUN cd /ethspam && make
FROM alpine:latest AS production
COPY --from=builder /ethspam/ethspam /usr/local/bin/
ENTRYPOINT ["ethspam"]
