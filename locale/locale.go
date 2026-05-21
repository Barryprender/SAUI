package locale

import "strings"

type Lang string

const (
	EN Lang = "en"
	ES Lang = "es"
)

func (l Lang) T(en, es string) string {
	if l == ES {
		return es
	}
	return en
}

func (l Lang) Prefix() string {
	if l == ES {
		return "/es"
	}
	return ""
}

// BasePath returns the path without any language prefix.
// Used to build hreflang alternate URLs.
func (l Lang) BasePath(currentPath string) string {
	if l == ES {
		return strings.TrimPrefix(currentPath, "/es")
	}
	return currentPath
}

func (l Lang) AlternatePath(currentPath string) string {
	if l == ES {
		return strings.TrimPrefix(currentPath, "/es")
	}
	return "/es" + currentPath
}

func (l Lang) AlternateHreflang() string {
	if l == ES {
		return "en"
	}
	return "es"
}

func (l Lang) AlternateLabel() string {
	if l == ES {
		return "EN"
	}
	return "ES"
}

func (l Lang) OGLocale() string {
	if l == ES {
		return "es_ES"
	}
	return "en_US"
}
