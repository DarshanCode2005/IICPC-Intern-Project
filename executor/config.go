package executor

import (
	"fmt"
	"strings"
	"time"
)

// LanguageConfig defines execution settings for a language
type LanguageConfig struct {
	Timeout time.Duration
	Args    func(code, stdin string) []string
}

// languageConfigs holds execution configurations for supported languages
var languageConfigs = map[string]LanguageConfig{
	"go": {
		Timeout: 10 * time.Second,
		Args: func(code, stdin string) []string {
			return []string{"sh", "-c", fmt.Sprintf(
				`printf '%%s' '%s' > /app/temp/input.txt && echo '%s' > /app/temp/code.go && go run /app/temp/code.go < /app/temp/input.txt`,
				sanitizePrintf(stdin),
				strings.ReplaceAll(code, "'", "'\\''"),
			)}
		},
	},
	"js": {
		Timeout: 10 * time.Second,
		Args: func(code, stdin string) []string {
			return []string{"sh", "-c", fmt.Sprintf(
				`printf '%%s' '%s' > /app/temp/input.txt && echo '%s' > /app/temp/code.js && node /app/temp/code.js < /app/temp/input.txt`,
				sanitizePrintf(stdin),
				strings.ReplaceAll(code, "'", "'\\''"),
			)}
		},
	},
	"python": {
		Timeout: 10 * time.Second,
		Args: func(code, stdin string) []string {
			return []string{"sh", "-c", fmt.Sprintf(
				`printf '%%s' '%s' > /app/temp/input.txt && echo '%s' > /app/temp/code.py && python3 /app/temp/code.py < /app/temp/input.txt`,
				sanitizePrintf(stdin),
				strings.ReplaceAll(code, "'", "'\\''"),
			)}
		},
	},
	"cpp": {
		Timeout: 10 * time.Second,
		Args: func(code, stdin string) []string {
			return []string{"sh", "-c", fmt.Sprintf(
				`printf '%%s' '%s' > /app/temp/input.txt && echo '%s' > /app/temp/code.cpp && g++ -o /app/temp/exe /app/temp/code.cpp && /app/temp/exe < /app/temp/input.txt`,
				sanitizePrintf(stdin),
				strings.ReplaceAll(code, "'", "'\\''"),
			)}
		},
	},
	"c": {
		Timeout: 10 * time.Second,
		Args: func(code, stdin string) []string {
			return []string{"sh", "-c", fmt.Sprintf(
				`printf '%%s' '%s' > /app/temp/input.txt && echo '%s' > /app/temp/code.c && gcc -o /app/temp/exe /app/temp/code.c && /app/temp/exe < /app/temp/input.txt`,
				sanitizePrintf(stdin),
				strings.ReplaceAll(code, "'", "'\\''"),
			)}
		},
	},
	"java": {
		Timeout: 10 * time.Second,
		Args: func(code, stdin string) []string {
			return []string{"sh", "-c", fmt.Sprintf(
				`printf '%%s' '%s' > /app/temp/input.txt && echo '%s' > /app/temp/Main.java && javac /app/temp/Main.java && java -cp /app/temp Main < /app/temp/input.txt`,
				sanitizePrintf(stdin),
				strings.ReplaceAll(code, "'", "'\\''"),
			)}
		},
	},
}

// GetLanguageConfig retrieves the configuration for a given language
func GetLanguageConfig(language string) (LanguageConfig, bool) {
	config, ok := languageConfigs[language]
	return config, ok
}

func sanitizePrintf(input string) string {
	escaped := strings.ReplaceAll(input, "%", "%%")
	return strings.ReplaceAll(escaped, "'", "'\\''")
}
