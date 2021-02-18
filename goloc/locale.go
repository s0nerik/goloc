package goloc

import (
	"fmt"
	"strings"
)

type locale struct {
	Lang    string
	Country string
}

type Locale = locale

func (l *Locale) String() string {
	if len(l.Lang) == 0 {
		return ""
	}
	if len(l.Country) == 0 {
		return l.Lang
	}
	return fmt.Sprintf("%s_%s", l.Lang, l.Country)
}

// ParseLocale returns a new Locale parsed from the standard string
// representation. Input value is case-insensitive.
//
// Supported formats: "<language>", "<language>_<COUNTRY CODE (ISO Alpha-2)>".
// Examples of supported locales: "en", "EN", "En", "en_US", "ru_UA", "ru_ua", "custom_DE"
func ParseLocale(localeStr string) (result Locale, err error) {
	localeParts := strings.Split(localeStr, "_")
	if len(localeParts) < 1 {
		return locale{}, &wrongLocaleFormatError{
			localeString: localeStr,
		}
	}
	result = locale{
		Lang: strings.ToLower(localeParts[0]),
	}
	if len(localeParts) < 2 {
		return
	}
	result.Country = strings.ToUpper(localeParts[1])
	return
}

type wrongLocaleFormatError struct {
	localeString string
}

func (e *wrongLocaleFormatError) Error() string {
	return fmt.Sprintf(`%v: unsupported locale format.
Supported formats: "<language>", "<language>_<COUNTRY CODE (ISO Alpha-2)>".
Examples of supported locales: "en", "EN", "En", "en_US", "ru_UA", "ru_ua", "custom_DE"
`, e.localeString)
}
