package translation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const api_template = "https://translate.yandex.net/api/v1.5/tr.json/translate?lang=%s&format=%s&text=%s&key=%s"

type YandexTranslator struct {
	apiKey string
}

func (t YandexTranslator) TranslateEnToVi(text string) (string, error) {
	url := fmt.Sprintf(api_template, "en-vi", "plain", url.QueryEscape(text), t.apiKey)

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error calling api: %v", err)
	}

	translateRes := YandexTranslateRes{}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(resBytes, &translateRes)
	if err != nil || len(translateRes.Text) == 0 {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return translateRes.Text[0], nil
}

type YandexTranslateRes struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}
