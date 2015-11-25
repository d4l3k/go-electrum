package wallet

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

var languages = []string{"chinese_simplified", "english", "japanse", "portuguese", "spanish"}
var ErrInvalidLanguage = errors.New("invalid seed language")

func wordlistDir() string {
	gopath := os.Getenv("GOPATH")
	return gopath + "/github.com/d4l3k/go-electrum/wordlist"
}

func wordsLanguage(lang string) ([]string, error) {
	exists := false
	for _, validLang := range languages {
		if validLang == lang {
			exists = true
			break
		}
	}
	if !exists {
		return nil, ErrInvalidLanguage
	}
	path := fmt.Sprintf("%s/%s.txt", wordlistDir(), lang)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(raw)), "\n"), nil
}

// GenerateSeed returns a random seed based off a word list in the specified language. Languages are "
func GenerateSeed(lang string) (string, error) {
	words, err := wordsLanguage(lang)
	if err != nil {
		return "", nil
	}
	var seed []string
	for i := 0; i < 13; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		if err != nil {
			return "", err
		}
		seed = append(seed, words[n.Int64()])
	}
	return strings.Join(seed, " "), nil
}
