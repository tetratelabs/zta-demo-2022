package main

import (
	"encoding/base64"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{contextID: contextID}
}

type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// Override proxywasm.DefaultHttpContext
func (*httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	subject := getClaimValue("name")

	if strings.Contains(subject, "jndi") {
		proxywasm.LogInfof("access denied for: %s", subject)
		if err := proxywasm.SendHttpResponse(403, nil, []byte("Access Denied\n"), -1); err != nil {
			proxywasm.LogErrorf("failed to send local response: %v", err)
			proxywasm.ResumeHttpRequest()
		}
		return types.ActionPause
	}

	proxywasm.LogInfof("access granted for: %s", subject)

	return types.ActionContinue
}

// getClaimValue returns the value of the given claim.
// This method assumes the claim has a string value.
func getClaimValue(claim string) string {
	headers, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		return ""
	}

	var auth string
	for _, h := range headers {
		if strings.ToLower(h[0]) == "authorization" {
			auth = h[1]
			break
		}
	}
	if auth == "" {
		proxywasm.LogInfof("no authorization header found")
		return "anonymous"
	}

	token := auth[strings.Index(auth, " ")+1:]
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		proxywasm.LogErrorf("invalid jwt token: %v", err)
		return "anonymous"
	}

	body, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		proxywasm.LogErrorf("invalid jwt token body: %v", err)
		return "anonymous"
	}

	res, err := jsonparser.GetString(body, claim)
	if err != nil {
		proxywasm.LogErrorf("invalid jwt token body: %v", err)
		return "anonymous"
	}

	return res
}
