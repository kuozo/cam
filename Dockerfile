FROM caddy:2.7.5-builder AS builder

RUN caddy-builder \
    github.com/kuozo/cam@master 

FROM caddy:2.7.5-builder

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
