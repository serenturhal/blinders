package common

const (
	Beginner     Level = iota
	Intermediate Level = iota
	Advanced     Level = iota
)

var (
	LangVi = NewLanguage("vi", "Vietnamese")
	LangEn = NewLanguage("en", "English")
)

type (
	Level int // [Beginner, Intermediate, Advanced]
	Lang  struct {
		Code string `json:"languageCode"` // ISO-[639-1] Code of language based
		Name string `json:"languageName"` // English name of lanague
	}
)

func (l Level) String() string {
	switch l {
	case Beginner:
		return "Beginner"
	case Intermediate:
		return "Intermediate"
	case Advanced:
		return "Advanced"
	default:
		return "Unknown"
	}
}

func (l Lang) String() string {
	return l.Name
}

func NewLanguage(code string, name string) Lang {
	return Lang{
		Code: code,
		Name: name,
	}
}
