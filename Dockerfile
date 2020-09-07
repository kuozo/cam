FROM caddy:2.1.1-builder AS builder

RUN caddy-builder \
    github.com/kuozo/cam@master 

FROM caddy:2.1.1

COPY --from=builder /usr/bin/caddy /usr/bin/caddy