package lang

import (
	"bufio"
	"fmt"
	"github.com/maxazimi/v2ray-gio/ui/lang/locale"
	"regexp"
	"strings"
)

var (
	langStrings map[string]map[string]string
	languages   = map[string]string{
		"en": "English",
		"fa": "فارسی",
	}
	defaultLanguage = "en"
	rex             = regexp.MustCompile(`(?m)("(?:\\.|[^"\\])*")\s*=\s*("(?:\\.|[^"\\])*")`) // "key"="value"
)

func init() {
	readIntoMap := func(localeStrings string) map[string]string {
		m := make(map[string]string)
		scanner := bufio.NewScanner(strings.NewReader(localeStrings))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "/") {
				continue
			}

			matches := rex.FindAllStringSubmatch(line, -1)
			if len(matches) == 0 {
				continue
			}

			kv := matches[0]
			key := trimQuotes(kv[1])
			value := trimQuotes(kv[2])

			m[key] = value
		}
		return m
	}

	langStrings = make(map[string]map[string]string)
	langStrings["en"] = readIntoMap(locale.EN)
	langStrings["fa"] = readIntoMap(locale.FA)
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func has(language string) bool {
	for key := range languages {
		if key == language {
			return true
		}
	}
	return false
}

func GetLanguages() (map[string]string, string) {
	return languages, defaultLanguage
}

func Set(lang string) {
	if has(lang) {
		if lang != defaultLanguage {
			defaultLanguage = lang
		}
	}
}

func Str(key string) string {
	langMap := langStrings[defaultLanguage]
	str, ok := langMap[key]
	if ok {
		return str
	}
	return key
}

func StrF(key string, a ...any) string {
	str := Str(key)
	if str == "" {
		return str
	}
	return fmt.Sprintf(str, a...)
}
