package model

import "runtime"

type Language struct {
	Language        string `json:"language"`
	LanguageVersion string `json:"language_version"`
}

func NewLanguage() *Language {
	language := &Language{
		Language:        "golang",
		LanguageVersion: runtime.Version(),
	}
	return language
}
