package proxy

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/clodevo/raven-proxy/pkg/acl"
	"github.com/clodevo/raven-proxy/pkg/config"
	"github.com/clodevo/raven-proxy/pkg/utils"
	"github.com/valyala/fasthttp"
)

var (
	defaultDialer fasthttp.TCPDialer
	fastclient    fasthttp.Client
)

func handleFastHTTP(ctx *fasthttp.RequestCtx, cfg *config.ProxyConfig) {
	if err := fastclient.DoTimeout(&ctx.Request, &ctx.Response, cfg.Timeout); err != nil {
		fmt.Printf("Client timeout: %s\n", err)
	}
}

func handleFastHTTPS(ctx *fasthttp.RequestCtx, cfg *config.ProxyConfig) {
	if len(ctx.Host()) > 0 {
		fmt.Printf("Connect to: %s\n", ctx.Host())
	}
	ctx.Hijack(func(clientConn net.Conn) {
		destConn, err := defaultDialer.DialTimeout(string(ctx.Host()), 10*time.Second)
		if err != nil {
			fmt.Printf("Dial timeout: %s\n", err)
			return
		}

		defer clientConn.Close()
		defer destConn.Close()

		go transfer(destConn, clientConn)
		transfer(clientConn, destConn)
	})
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("transfer: %s\n", err)
		}
	}()

	if _, err := io.Copy(destination, source); err != nil {
		fmt.Printf("transfer io closed: %s\n", err)
	}
}

func FastHTTPHandler(cfg *config.ProxyConfig, aclManager *acl.ACLManager) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !utils.Authenticate(ctx) {
			return
		}

		fmt.Printf("Tenant name for ACL check: %s\n", utils.TenantNameRequest)

		// if _, exists := aclManager.TenantLists[utils.TenantNameRequest]; !exists {
		aclManager.TenantLists[utils.TenantNameRequest] = aclManager.LoadTenantLists(utils.TenantNameRequest)
		// }

		if !aclManager.IsRequestAllowed(ctx, utils.TenantNameRequest) {
			utils.GetLogger().Debug("Request blocked by ACL policy")
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			ctx.Response.SetBodyString("Forbidden: The request is blocked by policy.")
			return
		}

		switch strings.ToUpper(string(ctx.Method())) {
		case fasthttp.MethodConnect:
			handleFastHTTPS(ctx, cfg)
		default:
			handleFastHTTP(ctx, cfg)
		}
	}
}
