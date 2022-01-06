package api

type ProgrammingLanguage string

const (
	Java   ProgrammingLanguage = "java"
	Go     ProgrammingLanguage = "go"
	Python ProgrammingLanguage = "python"
	Ruby   ProgrammingLanguage = "ruby"
	DotNet ProgrammingLanguage = "dotnet"
)

var (
	supportedLangs = []ProgrammingLanguage{Java, Go, Python, Ruby, DotNet}
)

func AvailableLanguages() []ProgrammingLanguage {
	return supportedLangs
}

func IsSupportedLanguage(lang string) bool {
	return containsLang(ProgrammingLanguage(lang), AvailableLanguages())
}

func containsLang(l ProgrammingLanguage, langs []ProgrammingLanguage) bool {
	for _, current := range langs {
		if l == current {
			return true
		}
	}

	return false
}
