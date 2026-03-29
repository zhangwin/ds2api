package openai

import (
	"regexp"
)

var emptyJSONFencePattern = regexp.MustCompile("(?is)```json\\s*```")
var leakedToolCallArrayPattern = regexp.MustCompile(`(?is)\[\{\s*"function"\s*:\s*\{[\s\S]*?\}\s*,\s*"id"\s*:\s*"call[^"]*"\s*,\s*"type"\s*:\s*"function"\s*}\]`)
var leakedToolResultBlobPattern = regexp.MustCompile(`(?is)<\s*\|\s*tool\s*\|\s*>\s*\{[\s\S]*?"tool_call_id"\s*:\s*"call[^"]*"\s*}`)

// leakedMetaMarkerPattern matches DeepSeek special tokens in BOTH forms:
//   - ASCII underscore: <｜end_of_sentence｜>
//   - U+2581 variant:   <｜end▁of▁sentence｜>  (used in some DeepSeek outputs)
var leakedMetaMarkerPattern = regexp.MustCompile(`(?i)<[｜\|]\s*(?:assistant|tool|end[_▁]of[_▁]sentence|end[_▁]of[_▁]thinking)\s*[｜\|]>`)

// leakedAgentXMLPattern catches agent-style XML tags that leak through when
// the sieve fails to capture them (e.g. incomplete blocks at stream end).
var leakedAgentXMLPattern = regexp.MustCompile(`(?is)</?(?:attempt_completion|ask_followup_question|new_task|result)>`)

func sanitizeLeakedOutput(text string) string {
	if text == "" {
		return text
	}
	out := emptyJSONFencePattern.ReplaceAllString(text, "")
	out = leakedToolCallArrayPattern.ReplaceAllString(out, "")
	out = leakedToolResultBlobPattern.ReplaceAllString(out, "")
	out = leakedMetaMarkerPattern.ReplaceAllString(out, "")
	out = leakedAgentXMLPattern.ReplaceAllString(out, "")
	return out
}
