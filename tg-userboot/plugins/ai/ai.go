package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all AI commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "ai",
		Aliases:     []string{"ask", "aic", "chat", "gemini"},
		Category:    "AI",
		Description: "Query AI with custom endpoint and model configured in .env",
		Usage:       ".ai [prompt or reply to message]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleAI,
	})

	d.Register(&types.Command{
		Name:        "img",
		Aliases:     []string{"image", "genimg", "imagine"},
		Category:    "AI",
		Description: "Generate an image using custom AI endpoint and model",
		Usage:       ".img [prompt]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleImage,
	})
}

// OpenAI-compatible Chat Completions structs
type openAIChatRequest struct {
	Model       string               `json:"model"`
	Messages    []openAIChatMessage  `json:"messages"`
	Temperature float32              `json:"temperature,omitempty"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// OpenAI-compatible Image Generation structs
type openAIImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type openAIImageResponse struct {
	Data []struct {
		URL           string `json:"url"`
		B64JSON       string `json:"b64_json"`
		RevisedPrompt string `json:"revised_prompt"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Gemini direct structs
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

func handleAI(ctx context.Context, b *types.Bot, msg *types.Message) error {
	apiKey := b.Config.AIAPIKey
	baseURL := b.Config.AIBaseURL
	model := b.Config.AIModel

	if apiKey == "" {
		return b.Respond(ctx, msg, "❌ <code>AI_API_KEY</code> (or <code>GEMINI_API_KEY</code>) is not configured in .env", 5*time.Second)
	}

	prompt := strings.TrimSpace(msg.Input)
	if prompt == "" && msg.ReplyToText != "" {
		prompt = msg.ReplyToText
	} else if prompt != "" && msg.ReplyToText != "" {
		prompt = fmt.Sprintf("Context: %s\n\nQuestion: %s", msg.ReplyToText, prompt)
	}

	if prompt == "" {
		return b.Respond(ctx, msg, "❌ <i>Please provide a prompt or reply to a text message.</i>", 5*time.Second)
	}

	_ = b.Respond(ctx, msg, fmt.Sprintf("🤖 <i>Thinking with %s...</i>", model), 0)

	// Route 1: Google Gemini direct endpoint if base URL points to Google API
	if strings.Contains(baseURL, "generativelanguage.googleapis.com") {
		return callGeminiDirect(ctx, b, msg, apiKey, model, prompt)
	}

	// Route 2: OpenAI-compatible custom endpoint (/v1/chat/completions)
	chatEndpoint := strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(chatEndpoint, "/chat/completions") {
		chatEndpoint += "/chat/completions"
	}

	reqPayload := openAIChatRequest{
		Model: model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: b.Config.AISystemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chatEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 45 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ AI Request failed: %v", err), 0)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var aiResp openAIChatResponse
	if err := json.Unmarshal(bodyBytes, &aiResp); err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to parse response: %s", string(bodyBytes)), 0)
	}

	if aiResp.Error != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ AI Error: %s", aiResp.Error.Message), 0)
	}

	if len(aiResp.Choices) == 0 {
		return b.Respond(ctx, msg, "⚠️ No response choices returned by AI.", 0)
	}

	replyText := aiResp.Choices[0].Message.Content
	if len(replyText) > 4000 {
		replyText = replyText[:4000] + "\n\n<i>[Response truncated]</i>"
	}

	finalMsg := fmt.Sprintf("🤖 <b>AI (%s):</b>\n\n%s", model, replyText)
	return b.Respond(ctx, msg, finalMsg, 0)
}

func callGeminiDirect(ctx context.Context, b *types.Bot, msg *types.Message, apiKey, model, prompt string) error {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	reqPayload := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqPayload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 35 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Gemini request failed: %v", err), 0)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to parse response: %v", err), 0)
	}

	if geminiResp.Error != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Gemini API Error: %s", geminiResp.Error.Message), 0)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return b.Respond(ctx, msg, "⚠️ No response generated by Gemini.", 0)
	}

	replyText := geminiResp.Candidates[0].Content.Parts[0].Text
	if len(replyText) > 4000 {
		replyText = replyText[:4000] + "\n\n<i>[Response truncated]</i>"
	}

	finalMsg := fmt.Sprintf("🤖 <b>Gemini (%s):</b>\n\n%s", model, replyText)
	return b.Respond(ctx, msg, finalMsg, 0)
}

func handleImage(ctx context.Context, b *types.Bot, msg *types.Message) error {
	apiKey := b.Config.ImageAPIKey
	baseURL := b.Config.ImageBaseURL
	model := b.Config.ImageModel

	if apiKey == "" {
		return b.Respond(ctx, msg, "❌ <code>IMAGE_API_KEY</code> (or <code>AI_API_KEY</code>) is not configured in .env", 5*time.Second)
	}

	prompt := strings.TrimSpace(msg.Input)
	if prompt == "" && msg.ReplyToText != "" {
		prompt = msg.ReplyToText
	}

	if prompt == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.img [image description]</code>", 5*time.Second)
	}

	_ = b.Respond(ctx, msg, fmt.Sprintf("🎨 <i>Generating image with %s...</i>\n<code>%s</code>", model, prompt), 0)

	imgEndpoint := strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(imgEndpoint, "/images/generations") {
		imgEndpoint += "/images/generations"
	}

	reqPayload := openAIImageRequest{
		Model:  model,
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", imgEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Image generation request failed: %v", err), 0)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var imgResp openAIImageResponse
	if err := json.Unmarshal(bodyBytes, &imgResp); err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to parse image response: %s", string(bodyBytes)), 0)
	}

	if imgResp.Error != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Image API Error: %s", imgResp.Error.Message), 0)
	}

	if len(imgResp.Data) == 0 || imgResp.Data[0].URL == "" {
		return b.Respond(ctx, msg, "❌ No image URL returned by endpoint.", 0)
	}

	imgURL := imgResp.Data[0].URL
	revised := imgResp.Data[0].RevisedPrompt
	if revised == "" {
		revised = prompt
	}

	output := fmt.Sprintf("🎨 <b>Generated Image (%s)</b>\n\n📝 <b>Prompt:</b> <i>%s</i>\n\n🖼️ <a href=\"%s\">Click to View / Download Image</a>",
		model, revised, imgURL)

	return b.Respond(ctx, msg, output, 0)
}
