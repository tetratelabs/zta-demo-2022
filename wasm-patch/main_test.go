package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
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

	host.CompleteHttpContext(contextID)

	// Check Envoy logs.
	// logs := host.GetInfoLogs()
	// require.Contains(t, logs, "check(nacx, GET, /public) = true")
}
