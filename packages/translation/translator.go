package translation

type Languages string

const (
	EN_VI Languages = "en-vi"
	VI_EN Languages = "vi-en"
)

type Translator interface {
	Translate(text string, langs Languages) (string, error)
}
