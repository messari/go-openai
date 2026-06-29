package openai_test

import (
	"encoding/json"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

// TestChatMessagePartCacheControlMarshal verifies that a CacheControl breakpoint
// on a MultiContent text part serializes into the OpenAI/LiteLLM wire format:
//
//	{"role":"system","content":[{"type":"text","text":"...","cache_control":{"type":"ephemeral"}}]}
func TestChatMessagePartCacheControlMarshal(t *testing.T) {
	msg := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		MultiContent: []openai.ChatMessagePart{
			{
				Type:         openai.ChatMessagePartTypeText,
				Text:         "stable system prompt",
				CacheControl: &openai.CacheControl{Type: "ephemeral"},
			},
		},
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	for _, want := range []string{
		`"content":[`,
		`"type":"text"`,
		`"text":"stable system prompt"`,
		`"cache_control":{"type":"ephemeral"}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("marshaled message missing %q\n got: %s", want, got)
		}
	}
}

// TestChatMessagePartCacheControlOmitted verifies the field is omitted when unset,
// so existing callers are unaffected.
func TestChatMessagePartCacheControlOmitted(t *testing.T) {
	msg := openai.ChatCompletionMessage{
		Role:         openai.ChatMessageRoleSystem,
		MultiContent: []openai.ChatMessagePart{{Type: openai.ChatMessagePartTypeText, Text: "hi"}},
	}
	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "cache_control") {
		t.Errorf("cache_control should be omitted when unset, got: %s", string(b))
	}
}
