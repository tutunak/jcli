package branch

import (
	"math/rand"
	"regexp"
	"strings"
	"unicode"
)

var (
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
	multipleHyphens = regexp.MustCompile(`-+`)
)

type Generator struct {
	randFunc func() int
}

func NewGenerator() *Generator {
	return &Generator{
		randFunc: func() int {
			return rand.Intn(1000000)
		},
	}
}

func NewGeneratorWithRand(randFunc func() int) *Generator {
	return &Generator{
		randFunc: randFunc,
	}
}

func (g *Generator) Generate(issueKey, summary string) string {
	normalized := normalizeSummary(summary)
	randomNum := g.randFunc()
	return strings.ToLower(issueKey) + "-" + normalized + "-" + formatNumber(randomNum)
}

func normalizeSummary(summary string) string {
	// Convert to lowercase
	s := strings.ToLower(summary)

	// Replace common unicode characters with ASCII equivalents
	s = replaceUnicode(s)

	// Replace non-alphanumeric characters with hyphens
	s = nonAlphanumeric.ReplaceAllString(s, "-")

	// Collapse multiple hyphens
	s = multipleHyphens.ReplaceAllString(s, "-")

	// Trim leading and trailing hyphens
	s = strings.Trim(s, "-")

	// Limit length to prevent overly long branch names
	if len(s) > 50 {
		s = s[:50]
		// Avoid cutting in the middle of a word
		if lastHyphen := strings.LastIndex(s, "-"); lastHyphen > 30 {
			s = s[:lastHyphen]
		}
		s = strings.Trim(s, "-")
	}

	return s
}

func replaceUnicode(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			result.WriteRune(r)
		case r >= '0' && r <= '9':
			result.WriteRune(r)
		case unicode.IsLetter(r):
			// Try to convert accented characters
			switch r {
			case 'à', 'á', 'â', 'ã', 'ä', 'å':
				result.WriteRune('a')
			case 'è', 'é', 'ê', 'ë':
				result.WriteRune('e')
			case 'ì', 'í', 'î', 'ï':
				result.WriteRune('i')
			case 'ò', 'ó', 'ô', 'õ', 'ö':
				result.WriteRune('o')
			case 'ù', 'ú', 'û', 'ü':
				result.WriteRune('u')
			case 'ñ':
				result.WriteRune('n')
			case 'ç':
				result.WriteRune('c')
			default:
				result.WriteRune('-')
			}
		default:
			result.WriteRune('-')
		}
	}
	return result.String()
}

func formatNumber(n int) string {
	return strings.TrimPrefix(strings.TrimPrefix(
		strings.TrimPrefix(strings.TrimPrefix(
			strings.TrimPrefix(padNumber(n), "0"), "0"), "0"), "0"), "0")
}

func padNumber(n int) string {
	s := "000000"
	num := n % 1000000
	result := s + itoa(num)
	return result[len(result)-6:]
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
