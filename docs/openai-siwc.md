# OpenAI SIWC / Sign in with ChatGPT token sharing

This branch adds an experimental OpenAI OAuth route for credentials issued by the
"Sign in with ChatGPT" token-sharing flow.

## What counts as an SIWC credential

The account must still use:

- `platform: openai`
- `type: oauth`

and its credentials must include either:

- `auth_mode: chatgpt_token_sharing` (preferred), or
- `auth_flow: chatgpt-token-sharing`, or
- a granted scope containing both `resource.invoke` and
  `chatgpt.tokens.use.direct` (the legacy direct scope name is also accepted).

A renewable SIWC profile also needs the issued dynamic client registration:

- `client_id: oaiapp_...`
- `refresh_token: ...`

Do not pair an SIWC refresh token with the Codex CLI client id
`app_EMoamEEZ73f0CkXaXp7hrann`. They are different OAuth registrations.

Example account fragment:

```json
{
  "platform": "openai",
  "type": "oauth",
  "credentials": {
    "access_token": "<current access token>",
    "refresh_token": "<current refresh token>",
    "client_id": "oaiapp_<issued registration>",
    "expires_at": "2026-10-02T14:00:00Z",
    "auth_mode": "chatgpt_token_sharing",
    "granted_scope": "openid profile email resource.invoke chatgpt.tokens.use.direct offline_access"
  }
}
```

## Upstream behavior

SIWC accounts are routed to:

```text
https://api.openai.com/v1/responses
```

instead of the ChatGPT Codex internal endpoint.

The gateway adds the current OSS preview header:

```text
x-openai-chatpass-test: codex-direct
```

and removes Codex-only routing/session headers.

The request body is normalized to `store=false` and strips the controls that
the direct token-sharing grant currently rejects:

- `context_management`
- `metadata`
- `max_output_tokens`
- `temperature`
- `top_p`
- `prompt_cache_retention`

## Refresh behavior

When `client_id` begins with `oaiapp_`, refresh uses:

```text
POST https://auth.openai.com/api/accounts/oauth/token
grant_type=refresh_token
client_id=<issued oaiapp client id>
refresh_token=<current refresh token>
resource=https://api.openai.com/v1
```

The returned refresh token may rotate. Sub2API's normal credential persistence
must therefore keep the newest returned refresh token.

## Current scope

This branch intentionally starts with imported/external SIWC credentials.
It does not yet implement the full dynamic-agent browser authorization flow
(`dynamic_agent_client` + PKCE + callback-issued `oaiapp_...` registration)
inside the Sub2API admin UI.

That can be added after the imported-credential path is confirmed against a real
account.
