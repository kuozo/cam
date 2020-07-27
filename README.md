# CAM

> Caddy Auth Middleware

## How to use

1. install `xcaddy` tool

```
go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
```

2. build

```
xcaddy build --with github.com/kuozo/cam
```

3. write your `Caddyfile`

```
order {
    cam before route
}

:9527 {
    guard {

    }

    route /apis/* {
        uri strip_prefix /apis
        reverse_proxy 0.0.0.0:8080
    }
}
```