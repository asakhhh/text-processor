package textprocessor

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func Run(src io.Reader, dst io.Writer) (err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = fmt.Errorf("%v", r)
	// 	}
	// }()

	var sb strings.Builder
	b := make([]byte, 100)
	n, err := src.Read(b)
	for err == nil {
		if _, err1 := sb.Write(b[:n]); err1 != nil {
			return err1
		}
		n, err = src.Read(b)
	}
	if err != io.EOF {
		return err
	}

	fields := strings.Fields(fixPunctuation(sb.String()))
	finalFields := make([]string, 0)
	for _, field := range fields {
		stripped := field
		trailingPunct := ""
		for len(stripped) > 0 && isPunct(rune(stripped[len(stripped)-1])) {
			trailingPunct = string(stripped[len(stripped)-1]) + trailingPunct
			stripped = stripped[:len(stripped)-1]
		}

		reattach := func() {
			if trailingPunct != "" && len(finalFields) > 0 {
				finalFields[len(finalFields)-1] += trailingPunct
			}
		}

		switch stripped {
		case "(hex)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Not enough words to apply transformation \"%v\" to", stripped)
			}
			err = applyHex(finalFields)
			if err != nil {
				return err
			}
			reattach()
		case "(bin)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Not enough words to apply transformation \"%v\" to", stripped)
			}
			err = applyBin(finalFields)
			if err != nil {
				return err
			}
			reattach()
		case "(up)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Not enough words to apply transformation \"%v\" to", stripped)
			}
			applyUp(finalFields)
			reattach()
		case "(low)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Not enough words to apply transformation \"%v\" to", stripped)
			}
			applyLow(finalFields)
			reattach()
		case "(cap)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Not enough words to apply transformation \"%v\" to", stripped)
			}
			applyCap(finalFields)
			reattach()
		default:
			if len(finalFields) != 0 &&
				(finalFields[len(finalFields)-1] == "(up," ||
					finalFields[len(finalFields)-1] == "(low," ||
					finalFields[len(finalFields)-1] == "(cap,") {
				if len(stripped) > 0 && stripped[len(stripped)-1] == ')' {
					num, err := strconv.Atoi(stripped[:(len(stripped) - 1)])
					if err != nil || num < 0 {
						return fmt.Errorf("Invalid token: \"%v\"", finalFields[len(finalFields)-1]+" "+field)
					}
					if num > len(finalFields)-1 {
						return fmt.Errorf("Not enough words to apply transformation \"%v\" to", finalFields[len(finalFields)-1]+" "+field)
					}
					last := finalFields[len(finalFields)-1]
					op := last[1:(len(last) - 1)]
					finalFields = finalFields[:len(finalFields)-1]
					switch op {
					case "up":
						applyUpN(finalFields, num)
					case "low":
						applyLowN(finalFields, num)
					case "cap":
						applyCapN(finalFields, num)
					}
					reattach()
				} else {
					finalFields = append(finalFields, stripped)
				}
			} else {
				if stripped == "(up," || stripped == "(low," || stripped == "(cap," {
					finalFields = append(finalFields, stripped)
				} else {
					finalFields = append(finalFields, field)
				}
			}
		}
	}

	finalText := strings.Join(finalFields, " ")
	finalText = fixArticles(finalText)
	_, err = dst.Write([]byte(finalText))
	if err != nil {
		return fmt.Errorf("Error while writing: %v", err.Error())
	}

	return nil
}

func isPunct(r rune) bool {
	return r == '.' || r == ',' || r == '?' || r == '!' || r == ':' || r == ';'
}

func fixPunctuation(s string) string {
	var tokens []string
	var current []rune
	var inPunct bool

	for _, r := range s {
		if isPunct(r) {
			if len(current) > 0 && !inPunct {
				tokens = append(tokens, string(current))
				current = []rune{}
			}
			inPunct = true
			current = append(current, r)
		} else if unicode.IsSpace(r) {
			if len(current) > 0 {
				tokens = append(tokens, string(current))
				current = []rune{}
			}
			inPunct = false
		} else {
			if len(current) > 0 && inPunct {
				tokens = append(tokens, string(current))
				current = []rune{}
			}
			inPunct = false
			current = append(current, r)
		}
	}

	if len(current) > 0 {
		tokens = append(tokens, string(current))
	}

	var result []rune
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if len(token) > 0 && isPunct(rune(token[0])) {
			if len(result) > 0 {
				result = append(result, []rune(token)...)
				result = append(result, ' ')
			}
		} else {
			if len(result) > 0 && result[len(result)-1] != ' ' {
				result = append(result, ' ')
			}
			result = append(result, []rune(token)...)
		}
	}

	return string(result)
}

func applyHex(fields []string) error {
	n, err := strconv.ParseInt(fields[len(fields)-1], 16, 0)
	if err != nil {
		return fmt.Errorf("invalid hexadecimal num: %v", fields[len(fields)-1])
	}
	fields[len(fields)-1] = fmt.Sprintf("%v", n)
	return nil
}

func applyBin(fields []string) error {
	n, err := strconv.ParseInt(fields[len(fields)-1], 2, 0)
	if err != nil {
		return fmt.Errorf("invalid binary num: %v", fields[len(fields)-1])
	}
	fields[len(fields)-1] = fmt.Sprintf("%v", n)
	return nil
}

func applyUp(fields []string) {
	fields[len(fields)-1] = strings.ToUpper(fields[len(fields)-1])
}

func applyUpN(fields []string, n int) {
	for i := range n {
		fields[len(fields)-1-i] = strings.ToUpper(fields[len(fields)-1-i])
	}
}

func applyLow(fields []string) {
	fields[len(fields)-1] = strings.ToLower(fields[len(fields)-1])
}

func applyLowN(fields []string, n int) {
	for i := range n {
		fields[len(fields)-1-i] = strings.ToLower(fields[len(fields)-1-i])
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

func applyCap(fields []string) {
	fields[len(fields)-1] = capitalize(fields[len(fields)-1])
}

func applyCapN(fields []string, n int) {
	for i := range n {
		fields[len(fields)-1-i] = capitalize(fields[len(fields)-1-i])
	}
}

func startsWithVowelOrH(word string) bool {
	if word == "" {
		return false
	}
	r := unicode.ToLower(rune(word[0]))
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'h'
}

func fixArticles(s string) string {
	words := strings.Fields(s)

	for i := 0; i < len(words)-1; i++ {
		if strings.ToLower(words[i]) == "a" {
			next := words[i+1]
			j := 0
			for j < len(next) && unicode.IsPunct(rune(next[j])) {
				j++
			}
			if j < len(next) && startsWithVowelOrH(next[j:]) {
				if words[i] == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}

	return strings.Join(words, " ")
}
