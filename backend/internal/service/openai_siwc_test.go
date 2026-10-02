package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountIsOpenAISIWCTokenSharing(t *testing.T) {
	t.Run("explicit auth mode", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Credentials: map[string]any{
				"client_id": "oaiapp_test_registration",
				"auth_mode": "chatgpt_token_sharing",
			},
		}
		require.True(t, account.IsOpenAISIWCTokenSharing())
		require.False(t, account.UsesOpenAICodexProtocol())
	})

	t.Run("scope fallback", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Credentials: map[string]any{
				"client_id":     "oaiapp_test_registration",
				"granted_scope": "openid profile resource.invoke chatgpt.tokens.use.direct offline_access",
			},
		}
		require.True(t, account.IsOpenAISIWCTokenSharing())
	})

	t.Run("marker without issued client id is rejected", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Credentials: map[string]any{
				"client_id": "app_EMoamEEZ73f0CkXaXp7hrann",
				"auth_mode": "chatgpt_token_sharing",
			},
		}
		require.False(t, account.IsOpenAISIWCTokenSharing())
		require.True(t, account.UsesOpenAICodexProtocol())
	})

	t.Run("ordinary codex oauth remains codex", func(t *testing.T) {
		account := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeOAuth,
			Credentials: map[string]any{
				"client_id": "app_EMoamEEZ73f0CkXaXp7hrann",
			},
		}
		require.False(t, account.IsOpenAISIWCTokenSharing())
		require.True(t, account.UsesOpenAICodexProtocol())
	})
}

func TestNormalizeOpenAISIWCPayload(t *testing.T) {
	input := []byte(`{
		"model":"gpt-6.1-sol",
		"input":"PING",
		"stream":true,
		"store":true,
		"metadata":{"source":"test"},
		"max_output_tokens":128,
		"temperature":0.2,
		"top_p":0.9,
		"prompt_cache_retention":"24h",
		"context_management":[{"type":"compaction"}],
		"reasoning":{"effort":"high"}
	}`)

	out, err := normalizeOpenAISIWCPayload(input)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(out, &got))
	require.Equal(t, false, got["store"])
	require.NotContains(t, got, "metadata")
	require.NotContains(t, got, "max_output_tokens")
	require.NotContains(t, got, "temperature")
	require.NotContains(t, got, "top_p")
	require.NotContains(t, got, "prompt_cache_retention")
	require.NotContains(t, got, "context_management")
	require.Equal(t, "gpt-6.1-sol", got["model"])
	require.Equal(t, true, got["stream"])

	reasoning, ok := got["reasoning"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "high", reasoning["effort"])
}
