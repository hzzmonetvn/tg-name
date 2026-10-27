package misc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all miscellaneous utility commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "weather",
		Aliases:     []string{"wttr"},
		Category:    "Misc",
		Description: "Get live weather forecast for a location",
		Usage:       ".weather [city]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleWeather,
	})

	d.Register(&types.Command{
		Name:        "cash",
		Aliases:     []string{"currency"},
		Category:    "Misc",
		Description: "Convert currencies using live exchange rates",
		Usage:       ".cash [amount] [from] [to] (e.g. .cash 100 USD VND)",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleCurrency,
	})

	d.Register(&types.Command{
		Name:        "tr",
		Aliases:     []string{"translate"},
		Category:    "Misc",
		Description: "Translate text to specified language (default: vi or en)",
		Usage:       ".tr [lang] [text or reply]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleTranslate,
	})

	d.Register(&types.Command{
		Name:        "wiki",
		Aliases:     []string{"wikipedia"},
		Category:    "Misc",
		Description: "Search Wikipedia summaries",
		Usage:       ".wiki [search query]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleWiki,
	})

	d.Register(&types.Command{
		Name:        "google",
		Aliases:     []string{"g"},
		Category:    "Misc",
		Description: "Generate Google Search URL",
		Usage:       ".google [query]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleGoogle,
	})
}

func handleWeather(ctx context.Context, b *types.Bot, msg *types.Message) error {
	city := strings.TrimSpace(msg.Input)
	if city == "" {
		city = b.Config.DefaultCity
	}
	if city == "" {
		city = "Hanoi"
	}

	_ = b.Respond(ctx, msg, fmt.Sprintf("🌤️ <i>Fetching weather for %s...</i>", city), 0)

	// Modern wttr.in JSON format (j1)
	reqURL := fmt.Sprintf("https://wttr.in/%s?format=j1", url.QueryEscape(city))
	req, _ := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	req.Header.Set("User-Agent", "curl/7.81.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Weather request failed: %v", err), 0)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Could not find weather for: %s", city), 5*time.Second)
	}

	var data struct {
		CurrentCondition []struct {
			TempC       string `json:"temp_C"`
			FeelsLikeC  string `json:"FeelsLikeC"`
			Humidity    string `json:"humidity"`
			Windspeed   string `json:"windspeedKmph"`
			WindDir     string `json:"winddir16Point"`
			WeatherDesc []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
		} `json:"current_condition"`
		NearestArea []struct {
			AreaName []struct {
				Value string `json:"value"`
			} `json:"areaName"`
			Country []struct {
				Value string `json:"value"`
			} `json:"country"`
		} `json:"nearest_area"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data.CurrentCondition) == 0 {
		return b.Respond(ctx, msg, "❌ Failed to parse weather information.", 5*time.Second)
	}

	curr := data.CurrentCondition[0]
	area := city
	country := ""
	if len(data.NearestArea) > 0 && len(data.NearestArea[0].AreaName) > 0 {
		area = data.NearestArea[0].AreaName[0].Value
		if len(data.NearestArea[0].Country) > 0 {
			country = data.NearestArea[0].Country[0].Value
		}
	}

	desc := "Clear"
	if len(curr.WeatherDesc) > 0 {
		desc = curr.WeatherDesc[0].Value
	}

	output := fmt.Sprintf(`🌤️ <b>Weather in %s, %s</b>
━━━━━━━━━━━━━━━━━
🌡️ <b>Temp:</b> <code>%s°C</code> (Feels: <code>%s°C</code>)
☁️ <b>Condition:</b> <code>%s</code>
💦 <b>Humidity:</b> <code>%s%%</code>
💨 <b>Wind:</b> <code>%s km/h (%s)</code>
━━━━━━━━━━━━━━━━━`,
		area, country,
		curr.TempC, curr.FeelsLikeC,
		desc,
		curr.Humidity,
		curr.Windspeed, curr.WindDir,
	)

	return b.Respond(ctx, msg, output, 0)
}

func handleCurrency(ctx context.Context, b *types.Bot, msg *types.Message) error {
	args := msg.Args
	if len(args) < 3 {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.cash [amount] [FROM] [TO]</code>\n<b>Example:</b> <code>.cash 100 USD EUR</code>", 5*time.Second)
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return b.Respond(ctx, msg, "❌ Invalid numeric amount.", 5*time.Second)
	}

	from := strings.ToUpper(args[1])
	to := strings.ToUpper(args[2])

	apiURL := fmt.Sprintf("https://api.frankfurter.dev/v1/latest?amount=%.2f&from=%s&to=%s", amount, from, to)
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Currency API request failed: %v", err), 0)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Conversion from %s to %s not supported.", from, to), 5*time.Second)
	}

	var data struct {
		Amount float64            `json:"amount"`
		Base   string             `json:"base"`
		Date   string             `json:"date"`
		Rates  map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return b.Respond(ctx, msg, "❌ Failed to parse exchange rate.", 0)
	}

	rate, ok := data.Rates[to]
	if !ok {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Rate for %s not found.", to), 0)
	}

	output := fmt.Sprintf("💱 <b>Currency Conversion</b>\n<code>%.2f %s</code> = <b>%.2f %s</b>\n📅 <i>Rate Date: %s</i>",
		amount, from, rate, to, data.Date)

	return b.Respond(ctx, msg, output, 0)
}

func handleTranslate(ctx context.Context, b *types.Bot, msg *types.Message) error {
	var targetLang string = "vi"
	var textToTranslate string

	if len(msg.Args) > 0 && len(msg.Args[0]) == 2 {
		targetLang = strings.ToLower(msg.Args[0])
		if len(msg.Args) > 1 {
			textToTranslate = strings.Join(msg.Args[1:], " ")
		}
	} else {
		textToTranslate = msg.Input
	}

	if textToTranslate == "" && msg.ReplyToText != "" {
		textToTranslate = msg.ReplyToText
	}

	if textToTranslate == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.tr [lang_code] [text or reply]</code>\n<b>Example:</b> <code>.tr en Xin chào thế giới</code>", 5*time.Second)
	}

	_ = b.Respond(ctx, msg, fmt.Sprintf("🌐 <i>Translating to %s...</i>", strings.ToUpper(targetLang)), 0)

	// Layer 1: MyMemory Translation API (reliable, no 429 datacenter blocking)
	myMemoryURL := fmt.Sprintf("https://api.mymemory.translated.net/get?q=%s&langpair=auto|%s",
		url.QueryEscape(textToTranslate), targetLang)

	req, _ := http.NewRequestWithContext(ctx, "GET", myMemoryURL, nil)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var memData struct {
			ResponseData struct {
				TranslatedText string `json:"translatedText"`
			} `json:"responseData"`
			ResponseStatus int `json:"responseStatus"`
		}
		if jsonErr := json.NewDecoder(resp.Body).Decode(&memData); jsonErr == nil && memData.ResponseData.TranslatedText != "" {
			output := fmt.Sprintf("🌐 <b>Translation (➔ %s):</b>\n\n%s",
				strings.ToUpper(targetLang), memData.ResponseData.TranslatedText)
			return b.Respond(ctx, msg, output, 0)
		}
	}
	if resp != nil {
		resp.Body.Close()
	}

	// Layer 2: Google Translate GTX endpoint fallback
	apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=t&q=%s",
		targetLang, url.QueryEscape(textToTranslate))

	req, _ = http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	resp, err = client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Translation failed: %v", err), 0)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var rawData []interface{}
	if err := json.Unmarshal(bodyBytes, &rawData); err != nil || len(rawData) == 0 {
		return b.Respond(ctx, msg, "❌ Failed to parse translation.", 0)
	}

	firstPart, ok := rawData[0].([]interface{})
	if !ok || len(firstPart) == 0 {
		return b.Respond(ctx, msg, "❌ Translation parsing error.", 0)
	}

	var translatedResult strings.Builder
	for _, segment := range firstPart {
		if segSlice, ok := segment.([]interface{}); ok && len(segSlice) > 0 {
			if strVal, ok := segSlice[0].(string); ok {
				translatedResult.WriteString(strVal)
			}
		}
	}

	detectedLang := "auto"
	if len(rawData) > 2 {
		if langStr, ok := rawData[2].(string); ok {
			detectedLang = langStr
		}
	}

	output := fmt.Sprintf("🌐 <b>Translation (%s ➔ %s):</b>\n\n%s",
		strings.ToUpper(detectedLang), strings.ToUpper(targetLang), translatedResult.String())

	return b.Respond(ctx, msg, output, 0)
}

func handleWiki(ctx context.Context, b *types.Bot, msg *types.Message) error {
	query := strings.TrimSpace(msg.Input)
	if query == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.wiki [search query]</code>", 5*time.Second)
	}

	apiURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(query))
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req.Header.Set("User-Agent", "TG-Userboot/1.0 (https://github.com/hzzmonetvn)")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ No Wikipedia page found for: %s", query), 5*time.Second)
	}
	defer resp.Body.Close()

	var data struct {
		Title       string `json:"title"`
		Extract     string `json:"extract"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Extract == "" {
		return b.Respond(ctx, msg, "❌ Could not extract article summary.", 5*time.Second)
	}

	output := fmt.Sprintf("📚 <b>%s</b>\n\n%s\n\n🔗 <a href=\"%s\">Read on Wikipedia</a>",
		data.Title, data.Extract, data.ContentURLs.Desktop.Page)

	return b.Respond(ctx, msg, output, 0)
}

func handleGoogle(ctx context.Context, b *types.Bot, msg *types.Message) error {
	query := strings.TrimSpace(msg.Input)
	if query == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.google [search query]</code>", 5*time.Second)
	}

	gURL := fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(query))
	output := fmt.Sprintf("🔍 <b>Google Search:</b> <a href=\"%s\">%s</a>", gURL, query)
	return b.Respond(ctx, msg, output, 0)
}
