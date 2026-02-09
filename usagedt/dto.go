package usagedt

import (
	"fmt"
	"strings"
)

type Usages []*Usage

func (us Usages) String() string {
	if len(us) == 0 {
		return "Usages[]"
	}
	var parts []string
	for i, u := range us {
		parts = append(parts, fmt.Sprintf("[%d] %s", i, u))
	}
	return fmt.Sprintf("Usages{\n  %s\n}", strings.Join(parts, "\n  "))
}

type Result struct {
	Usages `json:"result"`
}

type Usage struct {
	Usage  string  `json:"usage"`
	Reason string  `json:"reason"`
	Weight float64 `json:"weight"`
}

func (u *Usage) String() string {
	if u == nil {
		return "Usage(nil)"
	}
	return fmt.Sprintf(
		"Usage{Usage: %q, Reason: %q, Weight: %.2f}",
		u.Usage,
		u.Reason,
		u.Weight,
	)
}

type Message struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type MessageImageContent struct {
	Type     string   `json:"type"`
	ImageURL ImageURL `json:"image_url"`
}

type MessageTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
