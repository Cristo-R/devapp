// Generating slug tool prototypes comes from the samoyed service
// see https://gitlab.shoplazza.site/shoplaza/backend/product/samoyed/-/blob/develop/goat/slug/slug.go
package utils

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	replaceableCodes = make(map[rune]string)
	strippableCodes  = make(map[rune]struct{})

	MaxLength            = 255
	regexpMultipleDashes = regexp.MustCompile("-+")
)

func GenerateSlug(s string) string {
	handle := strings.TrimSpace(s)
	handle = strings.ToLower(handle)
	handle = doSubstitute(handle)

	handle = strings.TrimSpace(handle)
	handle = strings.ReplaceAll(handle, " ", "-")
	handle = strings.Trim(handle, "-")
	handle = regexpMultipleDashes.ReplaceAllString(handle, "-")
	handle = smartTruncate(handle)

	return handle
}

func doSubstitute(handle string) string {
	var buf bytes.Buffer

	for _, s := range strings.Split(handle, "") {
		ru, _ := utf8.DecodeRune([]byte(s))

		if d, ok := replaceableCodes[ru]; ok {
			buf.WriteString(d)
		} else if _, ok := strippableCodes[ru]; ok {
			continue
		} else {
			buf.WriteString(s)
		}
	}
	return buf.String()
}

func smartTruncate(text string) string {
	if len(text) < MaxLength {
		return text
	}

	var truncated string
	words := strings.SplitAfter(text, "-")
	// If MaxLength is smaller than length of the first word return word
	// truncated after MaxLength.
	if len(words[0]) > MaxLength {
		return words[0][:MaxLength]
	}
	for _, word := range words {
		if len(truncated)+len(word)-1 <= MaxLength {
			truncated += word
		} else {
			break
		}
	}
	return strings.Trim(truncated, "-")
}

var defaultSub = map[string]string{
	"×":  "x",
	"‐":  " ",
	"‑":  " ",
	"‒":  " ",
	"–":  " ",
	"—":  " ",
	"―":  " ",
	"~":  " ",
	"\n": " ",
}

// various kinds of space characters
var spaceSub = map[string]string{
	"\u0020":       " ",
	"\xc2\xa0":     " ",
	"\xe2\x80\x80": " ",
	"\xe2\x80\x81": " ",
	"\xe2\x80\x82": " ",
	"\xe2\x80\x83": " ",
	"\xe2\x80\x84": " ",
	"\xe2\x80\x85": " ",
	"\xe2\x80\x86": " ",
	"\xe2\x80\x87": " ",
	"\xe2\x80\x88": " ",
	"\xe2\x80\x89": " ",
	"\xe2\x80\x8a": " ",
	"\xe2\x81\x9f": " ",
	"\xe3\x80\x80": " ",
}

// Latin
var latinSub = map[string]string{
	"À":  "A",
	"Á":  "A",
	"Â":  "A",
	"Ã":  "A",
	"Ä":  "A",
	"Å":  "A",
	"Æ":  "Ae",
	"Ç":  "C",
	"È":  "E",
	"É":  "E",
	"Ê":  "E",
	"Ë":  "E",
	"Ì":  "I",
	"Í":  "I",
	"Î":  "I",
	"Ï":  "I",
	"Ð":  "D",
	"Ñ":  "N",
	"Ò":  "O",
	"Ó":  "O",
	"Ô":  "O",
	"Õ":  "O",
	"Ö":  "O",
	"Ø":  "O",
	"Ù":  "U",
	"Ú":  "U",
	"Û":  "U",
	"Ü":  "U",
	"Ý":  "Y",
	"Þ":  "Th",
	"ß":  "ss",
	"à":  "a",
	"á":  "a",
	"â":  "a",
	"ã":  "a",
	"ä":  "a",
	"å":  "a",
	"æ":  "ae",
	"ç":  "c",
	"è":  "e",
	"é":  "e",
	"ê":  "e",
	"ẽ":  "e",
	"ë":  "e",
	"ì":  "i",
	"í":  "i",
	"î":  "i",
	"ï":  "i",
	"ð":  "d",
	"ñ":  "n",
	"ò":  "o",
	"ó":  "o",
	"ô":  "o",
	"õ":  "o",
	"ö":  "o",
	"ø":  "o",
	"ù":  "u",
	"ú":  "u",
	"û":  "u",
	"ü":  "u",
	"ý":  "y",
	"þ":  "th",
	"ÿ":  "y",
	"Ā":  "A",
	"Ă":  "A",
	"Ą":  "A",
	"Ć":  "C",
	"Ĉ":  "C",
	"Ċ":  "C",
	"Č":  "C",
	"Ď":  "D",
	"Đ":  "D",
	"Ē":  "E",
	"Ĕ":  "E",
	"Ė":  "E",
	"Ę":  "E",
	"Ě":  "E",
	"Ĝ":  "G",
	"Ğ":  "G",
	"Ġ":  "G",
	"Ģ":  "G",
	"Ĥ":  "H",
	"Ħ":  "H",
	"Ĩ":  "I",
	"Ī":  "I",
	"Ĭ":  "I",
	"Į":  "I",
	"İ":  "I",
	"Ĳ":  "Ij",
	"Ĵ":  "J",
	"Ķ":  "K",
	"Ĺ":  "L",
	"Ļ":  "L",
	"Ľ":  "L",
	"Ŀ":  "L",
	"Ł":  "L",
	"Ń":  "N",
	"Ņ":  "N",
	"Ň":  "N",
	"Ŋ":  "Ng",
	"Ō":  "O",
	"Ŏ":  "O",
	"Ő":  "O",
	"Œ":  "OE",
	"Ŕ":  "R",
	"Ŗ":  "R",
	"Ř":  "R",
	"Ś":  "S",
	"Ŝ":  "S",
	"Ş":  "S",
	"Š":  "S",
	"Ţ":  "T",
	"Ť":  "T",
	"Ŧ":  "T",
	"Ũ":  "U",
	"Ū":  "U",
	"Ŭ":  "U",
	"Ů":  "U",
	"Ű":  "U",
	"Ų":  "U",
	"Ŵ":  "W",
	"Ŷ":  "Y",
	"Ÿ":  "Y",
	"Ź":  "Z",
	"Ż":  "Z",
	"Ž":  "Z",
	"ā":  "a",
	"ă":  "a",
	"å": "a",
	"ą̊": "a",
	"ą":  "a",
	"ć":  "c",
	"ĉ":  "c",
	"ċ":  "c",
	"č":  "c",
	"ď":  "d",
	"đ":  "d",
	"ē":  "e",
	"ĕ":  "e",
	"ė":  "e",
	"ę":  "e",
	"ě":  "e",
	"ĝ":  "g",
	"ğ":  "g",
	"ġ":  "g",
	"ģ":  "g",
	"ĥ":  "h",
	"ħ":  "h",
	"ĩ":  "i",
	"ī":  "i",
	"ĭ":  "i",
	"į":  "i",
	"ı":  "i",
	"ĳ":  "ij",
	"ĵ":  "j",
	"ķ":  "k",
	"ĸ":  "k",
	"ĺ":  "l",
	"ļ":  "l",
	"ľ":  "l",
	"ŀ":  "l",
	"ł":  "l",
	"ń":  "n",
	"ņ":  "n",
	"ň":  "n",
	"ŉ":  "n",
	"ŋ":  "ng",
	"ō":  "o",
	"ŏ":  "o",
	"ő":  "o",
	"ǫ":  "o",
	"ǭ":  "o",
	"œ":  "oe",
	"ŕ":  "r",
	"ŗ":  "r",
	"ř":  "r",
	"ś":  "s",
	"ŝ":  "s",
	"ş":  "s",
	"š":  "s",
	"ţ":  "t",
	"ť":  "t",
	"ŧ":  "t",
	"ũ":  "u",
	"ū":  "u",
	"ŭ":  "u",
	"ů":  "u",
	"ű":  "u",
	"ų":  "u",
	"ŵ":  "w",
	"ŷ":  "y",
	"ž":  "z",
	"ź":  "z",
	"ż":  "z",
	"€":  "eu",
}

func init() {
	for k, v := range defaultSub {
		c, _ := utf8.DecodeRune([]byte(k))
		replaceableCodes[c] = v
	}

	replaceToSpace := []rune{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 33, 35, 36, 37, 38,
		42, 43, 44, 46, 47, 58, 59, 60, 61, 62, 63, 64, 92, 94,
		96, 123, 124, 125, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 10084,
		137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 732, 65039,
		152, 153, 154, 155, 156, 157, 158, 159, 161, 162, 163, 164, 165, 166, 167, 8221, 12290, 65292, 12289,
		168, 171, 172, 173, 175, 176, 177, 178, 179, 180, 182, 183, 184, 12304, 12305, 12298, 8219,
		185, 187, 188, 189, 190, 191, 215, 247, 8203, 8204, 8205, 8239, 65279, 8216, 8223, 65281, 65311,
		65374, 8230, 12301, 12300, 65307, 65306, 12299, 8220, 8222, 8217, 65288, 65289,
	}
	for _, k := range replaceToSpace {
		replaceableCodes[k] = " "
	}

	strip := []rune{34, 39, 40, 41, 91, 93}
	for _, code := range strip {
		strippableCodes[code] = struct{}{}
	}

	for k, v := range latinSub {
		c, _ := utf8.DecodeRune([]byte(k))
		replaceableCodes[c] = v
	}

	for k, v := range spaceSub {
		c, _ := utf8.DecodeRune([]byte(k))
		replaceableCodes[c] = v
	}
}
