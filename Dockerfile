FROM caddy:2.5.0-builder AS builder

RUN caddy-builder \
    github.com/kuozo/cam@master 

FROM caddy:2.5.0-builder

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
