package prismbridge

import "testing"

func TestNormalizeModel(t *testing.T) {
	if got, err := NormalizeModel("prism-sol"); err != nil || got != "gpt-5.6-sol" {
		t.Fatalf("got=%q err=%v", got, err)
	}
	if got, err := NormalizeModel("prism-sol-6.1"); err != nil || got != "gpt-6.1-sol" {
		t.Fatalf("got=%q err=%v", got, err)
	}
	if got, err := NormalizeModel("prism-luna"); err != nil || got != "gpt-6-luna" {
		t.Fatalf("got=%q err=%v", got, err)
	}
	for _, model := range []string{"gpt-6.1-sol", "gpt-6-luna", "gpt-5.6-terra"} {
		if got, err := NormalizeModel(model); err != nil || got != model {
			t.Fatalf("model=%q got=%q err=%v", model, got, err)
		}
	}
	for _, model := range []string{"gpt-6-astra", "gpt-4o"} {
		if _, err := NormalizeModel(model); err == nil {
			t.Fatalf("unsupported model %q must fail", model)
		}
	}
}
