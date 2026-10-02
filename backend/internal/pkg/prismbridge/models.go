package prismbridge

import (
	"fmt"
	"strings"
)

var modelAliases = map[string]string{
	"prism-sol":     "gpt-5.6-sol",
	"prism-sol-6.1": "gpt-6.1-sol",
	"prism-luna":    "gpt-6-luna",
	"prism-terra":   "gpt-5.6-terra",
}

var supportedModels = map[string]bool{
	"gpt-6.1-sol":   true,
	"gpt-6-luna":    true,
	"gpt-5.6-sol":   true,
	"gpt-5.6-terra": true,
}

func NormalizeModel(model string) (string, error) {
	model = strings.TrimSpace(model)
	if alias, ok := modelAliases[model]; ok {
		model = alias
	}
	if !supportedModels[model] {
		return "", fmt.Errorf("unsupported Prism model: %s", model)
	}
	return model, nil
}
