package cam

import (
	"net/http"
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
	AuthURL   string   `json:"auth_url,omitempty"`
	PrefixURL []string `json:"prefix_url,omitempty"`
	logger    *zap.Logger
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
		case "auth_url":
			if len(args) != 1 {
				return d.Err("invalid verify url")
			}
			authURL := args[0]
			if !strings.HasPrefix(authURL, "/") {
				return d.Err("auth url like /*** format")
			}
			c.AuthURL = authURL
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
	url := r.URL.String()
	if !include(url, c.PrefixURL) {
		return next.ServeHTTP(w, r)
	}
	token := r.Header.Get("token")
	if token == "" {
		makeErrResp(w, 401, "token must value")
		return nil
	}
	if !verifyToken(jointURL(r.Host, c.AuthURL), token, url) {
		makeErrResp(w, 403, "permission denied")
		return nil
	}
	return next.ServeHTTP(w, r)
}

var (
	_ caddy.Provisioner           = (*Cam)(nil)
	_ caddy.Validator             = (*Cam)(nil)
	_ caddyhttp.MiddlewareHandler = (*Cam)(nil)
	_ caddyfile.Unmarshaler       = (*Cam)(nil)
)
