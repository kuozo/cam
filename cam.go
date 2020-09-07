package cam

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Cam{})
	httpcaddyfile.RegisterHandlerDirective("cam", parseCaddyfile)
}

// Cam .
type Cam struct {
	AuthEndpoint string   `json:"auth_endpoint,omitempty"`
	PrefixURL    []string `json:"prefix_url,omitempty"`
	AllowURL     []string `json:"allow_url,omitempty"`
	logger       *zap.Logger
}

// CaddyModule caddy module
func (Cam) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.cam",
		New: func() caddy.Module { return new(Cam) },
	}
}

// Provision provision
func (c *Cam) Provision(ctx caddy.Context) error {
	c.logger = ctx.Logger(c)
	return nil
}

//Validate validates that the module has a usable config.
func (c Cam) Validate() error {
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var c Cam
	err := c.UnmarshalCaddyfile(h.Dispenser)
	return c, err
}

// UnmarshalCaddyfile unmarshal caddy file
func (c *Cam) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {

		parameter := d.Val()
		args := d.RemainingArgs()
		switch parameter {
		case "prefix_url":
			if len(args) != 1 {
				return d.Err("invalid prefix url")
			}
			c.PrefixURL = c.splitPrefix(args[0])
		case "auth_endpoint":
			if len(args) != 1 {
				return d.Err("invalid verify url")
			}
			AuthEndpoint := args[0]
			if !strings.HasPrefix(AuthEndpoint, "http://") {
				return d.Err("auth endpoint like http://** format")
			}
			c.AuthEndpoint = AuthEndpoint
		case "allow_url":
			if len(args) != 1 {
				return d.Err("invalid allow url")
			}
			c.AllowURL = c.splitPrefix(args[0])
		default:
			d.Err("Unknow cam parameter: " + parameter)
		}

	}
	return nil
}

func (Cam) splitPrefix(data string) []string {
	return strings.Split(data, ",")
}

func (c Cam) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	// make sure the url filter
	url := r.URL.Path

	if include(url, c.AllowURL) {
		return next.ServeHTTP(w, r)
	}
	if !include(url, c.PrefixURL) {
		return next.ServeHTTP(w, r)
	}
	token := r.Header.Get("token")
	if token == "" {
		makeErrResp(w, 401, "token must value")
		return nil
	}
	ar := verifyToken(c.AuthEndpoint, token, url)
	if ar.Code != 200 {
		makeErrResp(w, ar.Code, ar.Message)
		return nil
	}
	r.Header.Add("x-user-id", strconv.Itoa(ar.Data.ID))
	r.Header.Add("x-user-type", strconv.Itoa(ar.Data.IsSuper))
	r.Header.Add("x-user-name", ar.Data.Name)
	return next.ServeHTTP(w, r)
}

var (
	_ caddy.Provisioner           = (*Cam)(nil)
	_ caddy.Validator             = (*Cam)(nil)
	_ caddyhttp.MiddlewareHandler = (*Cam)(nil)
	_ caddyfile.Unmarshaler       = (*Cam)(nil)
)
