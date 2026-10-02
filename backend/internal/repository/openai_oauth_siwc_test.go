package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/stretchr/testify/require"
)

func TestOpenAIOAuthServiceRefreshTokenWithClientIDUsesSIWCEndpoint(t *testing.T) {
	var gotResource string
	var gotScope string
	var gotClientID string
	var gotRefreshToken string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		gotResource = r.Form.Get("resource")
		gotScope = r.Form.Get("scope")
		gotClientID = r.Form.Get("client_id")
		gotRefreshToken = r.Form.Get("refresh_token")

		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{
			"access_token":"siwc-access-2",
			"refresh_token":"siwc-refresh-2",
			"token_type":"Bearer",
			"expires_in":3600,
			"scope":"openid profile resource.invoke chatgpt.tokens.use.direct offline_access"
		}`))
	}))
	defer server.Close()

	client := &openaiOAuthService{
		tokenURL:              "http://127.0.0.1:1/should-not-be-used",
		tokenSharingTokenURL: server.URL,
	}

	token, err := client.RefreshTokenWithClientID(
		context.Background(),
		"siwc-refresh-1",
		"",
		"oaiapp_test_registration",
	)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, "siwc-access-2", token.AccessToken)
	require.Equal(t, "siwc-refresh-2", token.RefreshToken)
	require.Equal(t, openai.TokenSharingResource, gotResource)
	require.Empty(t, gotScope)
	require.Equal(t, "oaiapp_test_registration", gotClientID)
	require.Equal(t, "siwc-refresh-1", gotRefreshToken)
}
