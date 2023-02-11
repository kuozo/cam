FROM caddy:2.6.3-builder AS builder

RUN caddy-builder \
    github.com/kuozo/cam@master 

FROM caddy:2.6.3-builder

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
