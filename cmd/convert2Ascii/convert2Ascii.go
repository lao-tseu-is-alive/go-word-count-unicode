package main

import (
	"fmt"
	anyascii "github.com/anyascii/go"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"unicode"
)

const (
	VERSION = "0.1.0"
	APP     = "convert2Ascii"
)

func UnicodeToASCII(input string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	// Normalize and convert the input string to ASCII by removing diacritics
	result, _, _ := transform.String(t, input)
	result = anyascii.Transliterate(result)
	return result
}

func checkConversion(unicodeTestString string, groundTruth string) {
	if UnicodeToASCII(unicodeTestString) == groundTruth {
		fmt.Printf("SUCCESS : '%s'\t is equivalent to:\t '%s'\n", unicodeTestString, groundTruth)
	} else {
		fmt.Printf("FAILURE : '%s'\t did not convert to:\t '%s'\n", unicodeTestString, groundTruth)
	}
}

func main() {
	l := log.New(os.Stdout, APP, log.Ldate|log.Ltime|log.Lshortfile)
	if len(os.Args) > 1 {
		unicodeInput := os.Args[1]
		fmt.Printf("'%s'\t is equivalent in ASCII to:\t '%s'\n", unicodeInput, UnicodeToASCII(unicodeInput))
	} else {
		l.Println("Expecting first argument to be an Unicode String to convert to Ascii")
		asciiString := "Cette fille aime lire"
		l.Printf("Usage:  %s \"Ⓒⓔⓣⓣⓔ ⓕⓘⓛⓛⓔ ⓐⓘⓜⓔ ⓛⓘⓡⓔ\"", APP)
		checkConversion("Ⓒⓔⓣⓣⓔ ⓕⓘⓛⓛⓔ ⓐⓘⓜⓔ ⓛⓘⓡⓔ", asciiString)
		checkConversion("𝓒𝓮𝓽𝓽𝓮 𝓯𝓲𝓵𝓵𝓮 𝓪𝓲𝓶𝓮 𝓵𝓲𝓻𝓮", asciiString)
		checkConversion("ℂ𝕖𝕥𝕥𝕖 𝕗𝕚𝕝𝕝𝕖 𝕒𝕚𝕞𝕖 𝕝𝕚𝕣𝕖", asciiString)
		checkConversion("C̶e̶t̶t̶e̶ ̶f̶i̶l̶l̶e̶ ̶a̶i̶m̶e̶ ̶l̶i̶r̶e̶", asciiString)
	}
}
