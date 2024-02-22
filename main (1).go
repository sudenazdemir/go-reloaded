package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	inputText, err := readFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	modifiedText := modifyText(inputText)

	err = writeFile(outputFile, modifiedText)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}

	fmt.Println("Modifications applied successfully!")
}

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var text string
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	return text, scanner.Err()
}

func writeFile(filename string, text string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}

func modifyText(text string) string {
	words := strings.Fields(text)
	modifiedWords := make([]string, len(words))

	for i := 0; i < len(words); i++ {
		switch {
		case words[i] == "(up)":
			if i > 0 {
				modifiedWords[i-1] = strings.ToUpper(words[i-1])
			}
		case words[i] == "(low)":
			if i > 0 {
				modifiedWords[i-1] = strings.ToLower(words[i-1])
			}
		case words[i] == "(cap)":
			if i > 0 {
				modifiedWords[i-1] = strings.Title(strings.ToLower(words[i-1]))
			}
		case words[i] == "(hex)":
			if i > 0 {
				hexValue, _ := strconv.ParseInt(words[i-1], 16, 64)
				modifiedWords[i-1] = strconv.FormatInt(hexValue, 10)
			}
		case words[i] == "(bin)":
			if i > 0 {
				binValue, _ := strconv.ParseInt(words[i-1], 2, 64)
				modifiedWords[i-1] = strconv.FormatInt(binValue, 10)
			}
		case i > 0 && words[i] == "(up,":
			re := regexp.MustCompile(`\d+`)
			match := re.FindString(words[i+1])
			num, _ := strconv.Atoi(match)
			for j := 0; j <= num; j++ {
				modifiedWords[i-j] = strings.ToUpper(words[i-j])
			}
		case i > 0 && words[i] == "(cap,":
			re := regexp.MustCompile(`\d+`)
			match := re.FindString(words[i+1])
			num, _ := strconv.Atoi(match)
			for j := 0; j <= num; j++ {
				modifiedWords[i-j] = strings.Title(words[i-j])
			}

		case i > 0 && words[i] == "(low,":
			re := regexp.MustCompile(`\d+`)
			match := re.FindString(words[i+1])
			num, _ := strconv.Atoi(match)
			for j := 0; j <= num; j++ {
				modifiedWords[i-j] = strings.ToLower(words[i-j])
			}

		default:

			modifiedWords[i] = words[i]
		}

		if words[i] == "a" && startVowel(words[i+1]) {
			words[i] = "an"
			modifiedWords[i] = words[i]
		} else if words[i] == "A" && startVowel(words[i+1]) {
			words[i] = "An"
			modifiedWords[i] = words[i]
		}

	}
	cleanedText := strings.Join(modifiedWords, " ")
	cleanedText = regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(cleanedText, "")
	// Birden fazla noktalama işareti yan yana ise bir önceki kelime ile aradaki boşluğu kaldırıp bitişik yazacak
	re := regexp.MustCompile(`(\w)\s+([.,?!;:])\s+`)
	cleanedText = re.ReplaceAllString(cleanedText, "$1$2 ")

	re = regexp.MustCompile(`([.,?!;:])\s+(\w)`)
	cleanedText = re.ReplaceAllString(cleanedText, "$1 $2")

	// Eğer tek noktalama işareti görürse bir önceki kelimeye bitişik yazacak
	re = regexp.MustCompile(`(\w)\s+([.,?!;])`)
	cleanedText = re.ReplaceAllString(cleanedText, "$1$2")
	// Eğer tek noktalama işareti görürse sonraki kelime ile arasında boşluk koyacak
	re = regexp.MustCompile(`([,])(\w)`)
	cleanedText = re.ReplaceAllString(cleanedText, "$1 $2")

	// Eğer tek tırnak işareti görürse tırnaktan sonra gelen kelime ile arasındaki boşluğu kaldıracak
	re = regexp.MustCompile(`'\s+(\w)`)
	cleanedText = re.ReplaceAllString(cleanedText, "'$1")

	// Eğer belli bir kısım tek tırnaklar arasına alınmışsa ilk tırnak ile sonraki kelime arasındaki boşluğu kaldıracak
	re = regexp.MustCompile(`'\s+(\w+)\s+'`)
	cleanedText = re.ReplaceAllString(cleanedText, "'$1'")

	// Son tırnak ile son kelime arasındaki boşluğu da kaldıracak
	re = regexp.MustCompile(`(\w+)\s+'`)
	cleanedText = re.ReplaceAllString(cleanedText, "$1'")

	re = regexp.MustCompile(`([.?!])\s+(\')`)
	cleanedText = re.ReplaceAllString(cleanedText, "$1$2")

	re = regexp.MustCompile(`\s+`)
	cleanedText = re.ReplaceAllString(cleanedText, " ")

	return cleanedText
}

func startVowel(s string) bool {
	vowel := "aeiouAEIOU"
	first := s[0]
	var result bool
	for i := 0; i < len(vowel); i++ {
		if first == vowel[i] {
			result = true
			break
		}
	}
	return result
}
