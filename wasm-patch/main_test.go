package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var (
	// malicious JWT token wfor log4shell wth the following claim:
	// "sub": "${jndi:ldap://localhost:1389/probably_not_vulnerable}"
	malicious = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJUZXRyYXRlIiwiaWF0IjoxNjQyNTE1MzI2LCJleHAiOjE2NzQwNTEzMjYsImF1ZCI6Ind3dy5leGFtcGxlLmNvbSIsInN1YiI6IiR7am5kaTpsZGFwOi8vbG9jYWxob3N0OjEzODkvcHJvYmFibHlfbm90X3Z1bG5lcmFibGV9IiwiR2l2ZW5OYW1lIjoiSm9obm55IiwiU3VybmFtZSI6IlJvY2tldCIsIkVtYWlsIjoianJvY2tldEBleGFtcGxlLmNvbSIsIlJvbGUiOlsiTWFuYWdlciIsIlByb2plY3QgQWRtaW5pc3RyYXRvciJdfQ.9j2utCHdvZcNCfbYzrqkouF7-uDsx5uEhZzZDycVto8"
)

func TestHttpFilter_OnHttpRequestHeaders(t *testing.T) {
	opt := proxytest.NewEmulatorOption().WithVMContext(&vmContext{})
	host, reset := proxytest.NewHostEmulator(opt)
	defer reset()

	// Call OnPluginStart -> the metric is initialized.
	status := host.StartPlugin()
	// Check the status returned by OnPluginStart is OK.
	require.Equal(t, types.OnPluginStartStatusOK, status)

	// Create http context.
	contextID := host.InitializeHttpContext()

	// Call OnHttpRequestHeaders no user
	action := host.CallOnRequestHeaders(contextID, [][2]string{}, false)
	require.Equal(t, types.ActionContinue, action)

	action = host.CallOnRequestHeaders(contextID, [][2]string{
		{"Authorization", "Bearer " + malicious},
	}, false)
	localResponse := host.GetSentLocalResponse(contextID)
	require.Equal(t, types.ActionPause, action)
	require.NotNil(t, localResponse)
	require.Equal(t, uint32(403), localResponse.StatusCode)

	host.CompleteHttpContext(contextID)

	// Check Envoy logs.
	logs := host.GetInfoLogs()
	require.Contains(t, logs, "no authorization header found")
	require.Contains(t, logs, "access granted for: anonymous")
	require.Contains(t, logs, "access denied for: ${jndi:ldap://localhost:1389/probably_not_vulnerable}")
}
