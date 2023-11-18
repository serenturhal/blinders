package translation

type Translator interface {
	TranslateEnToVi(text string) (string, error)
}
