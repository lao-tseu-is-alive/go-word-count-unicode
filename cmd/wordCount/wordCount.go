package main

import (
	"encoding/json"
	"fmt"
	anyascii "github.com/anyascii/go"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	VERSION = "0.1.0"
	APP     = "wordCount"
)

/*
better use a type to store word

	type WordCounter struct {
		err error
		msg string
		List map[string]int
	}
*/
func getRuneType(c rune) string {
	var runeType []string
	if unicode.IsControl(c) {
		runeType = append(runeType, "control")
	}
	if unicode.IsDigit(c) {
		runeType = append(runeType, "digit")
	}
	if unicode.IsGraphic(c) {
		runeType = append(runeType, "graphic")
	}
	if unicode.IsLetter(c) {
		runeType = append(runeType, "letter")
	}
	if unicode.IsLower(c) {
		runeType = append(runeType, "lower case")
	}
	if unicode.IsMark(c) {
		runeType = append(runeType, "mark")
	}
	if unicode.IsNumber(c) {
		runeType = append(runeType, "number")
	}
	if unicode.IsPrint(c) {
		runeType = append(runeType, "printable")
	}
	if !unicode.IsPrint(c) {
		runeType = append(runeType, "not printable")
	}
	if unicode.IsPunct(c) {
		runeType = append(runeType, "punct")
	}
	if unicode.IsSpace(c) {
		runeType = append(runeType, "space")
	}
	if unicode.IsSymbol(c) {
		runeType = append(runeType, "symbol")
	}
	if unicode.IsTitle(c) {
		runeType = append(runeType, "title case")
	}
	if unicode.IsUpper(c) {
		runeType = append(runeType, "upper case")
	}
	return strings.Join(runeType, ",")
}

func CountWords(buf []byte, log *log.Logger, minWordLength int, toLower bool, removeAccent bool) map[string]int {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	wordDic := make(map[string]int, 100)
	lineCount := 1
	wordCount := 0
	runeCount := 0
	runeLetterCount := 0
	size := 0
	col := 0
	currentWord := ""
	for start := 0; start < len(buf); start += size {
		var r rune
		r, size = utf8.DecodeRune(buf[start:])
		if r == utf8.RuneError {
			b := buf[start]
			fmt.Printf("[%d:%d](%v)\t%8v\t%#6x\t%#U\n",
				lineCount, col, start, b, b, b)
			log.Printf("ERROR 💥 invalid utf8 encoding at offset %d", start)
		}
		runeCount += 1
		col += size
		if r == '\n' {
			fmt.Printf("### eol found at line: %d\tcol: %d\t currentWord : %s\n", lineCount, runeCount, currentWord)
			col = 0
			lineCount += 1
			if utf8.RuneCountInString(currentWord) > minWordLength {
				wordCount += 1
				if toLower == true {
					currentWord = strings.ToLower(currentWord)
				}
				if removeAccent == true {
					unAccentedWord, _, _ := transform.String(t, currentWord)
					wordDic[unAccentedWord]++
				} else {
					wordDic[currentWord]++
				}
			}
			currentWord = ""
			fmt.Printf("###### starting line %d\t ######\n", lineCount)
		} else if unicode.IsLetter(r) {
			runeLetterCount += 1
			currentWord += anyascii.TransliterateRune(r)
		} else if unicode.IsSpace(r) {
			fmt.Printf("### space found at line: %d\tcol: %d\t currentWord : %s\n", lineCount, runeCount, currentWord)
			if utf8.RuneCountInString(currentWord) > minWordLength {
				wordCount += 1
				if toLower == true {
					currentWord = strings.ToLower(currentWord)
				}
				if removeAccent == true {
					unAccentedWord, _, _ := transform.String(t, currentWord)
					wordDic[unAccentedWord]++
				} else {
					wordDic[currentWord]++
				}
			}
			currentWord = ""
		} else {
			fmt.Printf("### discarded : [%d:%d]\t%8v\t%#6x\t%#U\t%#U\t%s\t['%s']\t(%s)\n",
				lineCount, runeCount, r, r, r,
				unicode.SimpleFold(r),
				runenames.Name(r), anyascii.TransliterateRune(r), getRuneType(r))
		}
	}
	fmt.Printf("### end found at line: %d\tcol: %d\t currentWord : %s\n", lineCount, runeCount, currentWord)
	// need to store the last word
	if utf8.RuneCountInString(currentWord) > minWordLength {
		wordCount += 1
		if toLower == true {
			currentWord = strings.ToLower(currentWord)
		}
		if removeAccent == true {
			unAccentedWord, _, _ := transform.String(t, currentWord)
			wordDic[unAccentedWord]++
		} else {
			wordDic[currentWord]++
		}
	}
	fmt.Printf("Num lines : %d,\tNum words: %d,\tNum Runes: %d,\tNum Letters:%d", lineCount, wordCount, runeCount, runeLetterCount)
	return wordDic
}

func analyseBuffer(buf []byte, log *log.Logger) {
	const header = "#[line:col] byte offset  decimal\thex\tUnicode\t\tSimple\tName\tAscii\ttype"
	size := 0
	lineCounter := 1
	col := 0
	fmt.Println(header)
	for start := 0; start < len(buf); start += size {
		var r rune
		r, size = utf8.DecodeRune(buf[start:])
		if r == utf8.RuneError {
			b := buf[start]
			fmt.Printf("[%d:%d](%v)\t%8v\t%#6x\t%#U\n",
				lineCounter, col, start, b, b, b)
			log.Printf("ERROR 💥 invalid utf8 encoding at offset %d", start)
		}
		if r == '\n' {
			col = 0
			lineCounter += 1
			fmt.Printf("### starting line %d\t ###\n", lineCounter)
			fmt.Println(header)
		} else {
			fmt.Printf("[%d:%d](%v)\t%8v\t%#6x\t%#U\t%#U\t%s\t['%s']\t(%s)\n",
				lineCounter, col, start, r, r, r,
				unicode.SimpleFold(r),
				runenames.Name(r), anyascii.TransliterateRune(r), getRuneType(r))
			col += size
		}
	}
}

func main() {
	l := log.New(os.Stdout, APP, log.Ldate|log.Ltime|log.Lshortfile)
	if len(os.Args) > 1 {
		filename := os.Args[1]
		l.Println(filename)
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("unable to read file %s. Error: %v", filename, err)
		}
		res := CountWords(content, l, 1, true, true)
		// fmt.Println(res)
		resText, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Println("error at json.MarshalIndent", err)
		}
		fmt.Printf("\n\nListe finale de %d mots   : \n %s\n", len(res), resText)

	} else {
		l.Fatal("Expecting first argument to be the text filename ")
	}
}
