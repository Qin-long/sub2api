package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
)

const (
	OpenAIAuthModeChatGPTTokenSharing = "chatgpt_token_sharing"

	openAISIWCCredentialAuthModeKey     = "auth_mode"
	openAISIWCCredentialAuthFlowKey     = "auth_flow"
	openAISIWCCredentialGrantedScopeKey = "granted_scope"

	openAISIWCPreviewHeader = "x-openai-chatpass-test"
	openAISIWCPreviewValue  = "codex-direct"
)

// IsOpenAISIWCTokenSharing reports whether this OAuth account is a
// Sign in with ChatGPT token-sharing credential.
//
// The explicit auth_mode/auth_flow marker is preferred. As a compatibility
// fallback we also accept credentials whose granted scope contains both
// resource.invoke and the direct ChatGPT token-sharing permission.
func (a *Account) IsOpenAISIWCTokenSharing() bool {
	if a == nil || !a.IsOpenAIOAuth() {
		return false
	}
	// Token-sharing registrations are dynamically issued as oaiapp_* clients.
	// Requiring that registration prevents a normal Codex OAuth account from
	// being accidentally routed to the public token-sharing endpoint by a stale
	// or manually edited auth_mode flag.
	if !openai.IsTokenSharingClientID(a.GetCredential("client_id")) {
		return false
	}

	for _, key := range []string{openAISIWCCredentialAuthModeKey, openAISIWCCredentialAuthFlowKey} {
		switch strings.ToLower(strings.TrimSpace(a.GetCredential(key))) {
		case OpenAIAuthModeChatGPTTokenSharing, "chatgpt-token-sharing", "siwc":
			return true
		}
	}

	scope := strings.TrimSpace(a.GetCredential(openAISIWCCredentialGrantedScopeKey))
	if scope == "" {
		scope = strings.TrimSpace(a.GetCredential("scope"))
	}
	return IsOpenAISIWCDirectScope(scope)
}

// IsOpenAISIWCDirectScope verifies the token-sharing grant returned by OpenAI.
func IsOpenAISIWCDirectScope(scope string) bool {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return false
	}
	granted := make(map[string]struct{})
	for _, item := range strings.Fields(scope) {
		granted[item] = struct{}{}
	}
	_, invoke := granted["resource.invoke"]
	_, direct := granted["chatgpt.tokens.use.direct"]
	if !direct {
		_, direct = granted["chatpass.enable.request.direct"]
	}
	return invoke && direct
}

// normalizeOpenAISIWCPayload applies the public Responses restrictions used by
// OpenAI's ChatGPT token-sharing flow. The grant is stateless, so storage and
// server-side replay controls are removed.
func normalizeOpenAISIWCPayload(body []byte) ([]byte, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("empty Responses request body")
	}

	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, fmt.Errorf("decode SIWC Responses request: %w", err)
	}

	request["store"] = false
	for _, field := range []string{
		"context_management",
		"metadata",
		"max_output_tokens",
		"temperature",
		"top_p",
		"prompt_cache_retention",
	} {
		delete(request, field)
	}

	normalized, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode SIWC Responses request: %w", err)
	}
	return normalized, nil
}
