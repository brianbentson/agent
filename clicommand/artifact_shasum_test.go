package clicommand

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buildkite/agent/v3/logger"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.RequestURI() {
		case "/builds/buildid/artifacts/search?query=foo.%2A&state=finished":
			io.WriteString(rw, `[{"path": "foo.txt", "sha1sum": "theshastring", "sha256sum": "thesha256string"}]`)
		default:
			t.Errorf("unexpected HTTP request: %s %v", req.Method, req.URL.RequestURI())
		}
	}))
}

func TestSearchAndPrintSha1Sum(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	cfg := ArtifactShasumConfig{
		Query:            "foo.*",
		Build:            "buildid",
		AgentAccessToken: "agentaccesstoken",
		Endpoint:         server.URL,
	}
	l := logger.NewBuffer()
	stdout := new(bytes.Buffer)

	searchAndPrintShaSum(cfg, l, stdout)

	assert.Equal(t, "theshastring\n", stdout.String())

	assert.Contains(t, l.Messages, `[info] Searching for artifacts: "foo.*"`)
	assert.Contains(t, l.Messages, `[debug] Artifact "foo.txt" found`)
}

func TestSearchAndPrintSha256Sum(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	cfg := ArtifactShasumConfig{
		Query:            "foo.*",
		Build:            "buildid",
		Sha256:           true,
		AgentAccessToken: "agentaccesstoken",
		Endpoint:         server.URL,
	}
	l := logger.NewBuffer()
	stdout := new(bytes.Buffer)

	searchAndPrintShaSum(cfg, l, stdout)

	assert.Equal(t, "thesha256string\n", stdout.String())

	assert.Contains(t, l.Messages, `[info] Searching for artifacts: "foo.*"`)
	assert.Contains(t, l.Messages, `[debug] Artifact "foo.txt" found`)
}
