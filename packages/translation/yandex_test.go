package translation

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

var translator Translator

func init() {
	cwd, _ := os.Getwd()
	fmt.Println("test", cwd)
	godotenv.Load("../../.env.test")
	translator = YandexTranslator{apiKey: os.Getenv("YANDEX_API_KEY")}
}

func TestTranslateWord(t *testing.T) {
	text := "absolutely"
	expectedResult := "tuyệt đối"
	fmt.Printf("translate \"%s\" to vietnamese, expect \"%s\"\n", text, expectedResult)

	result, err := translator.TranslateEnToVi(text)
	if err != nil {
		t.Error(err)
	}

	if strings.ToLower(result) != strings.ToLower(expectedResult) {
		t.Errorf("received = \"%s\", expect \"%s\"", result, expectedResult)
	}
}

func TestTranslateSentence(t *testing.T) {
	text := "hello, My name is Peakee"
	expectedResult := "xin chào, tên tôi là Peakee"
	fmt.Printf("translate \"%s\" to vietnamese, expect \"%s\"\n", text, expectedResult)

	result, err := translator.TranslateEnToVi(text)
	if err != nil {
		t.Error(err)
	}

	if strings.ToLower(result) != strings.ToLower(expectedResult) {
		t.Errorf("received = \"%s\", expect \"%s\"", result, expectedResult)
	}
}
