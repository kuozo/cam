FROM caddy:2.0.0-builder AS builder

RUN caddy-builder \
    github.com/kuozo/cam@master 

FROM caddy:2.0.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy