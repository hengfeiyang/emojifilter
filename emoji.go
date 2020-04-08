// Package emojifilter provider filter emoji from text returns without emoji
// emoji data file: https://unicode.org/Public/13.0.0/ucd/emoji/emoji-data.txt
package emojifilter

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-baa/common/util"
)

type emoji struct {
	data map[rune]bool
}

var _defaultFilter = newEmoji()

// Filter filter emoji from str returns without emoji
func Filter(str string) string {
	rs := []rune(str)
	ns := []rune{}
	for _, c := range rs {
		if _defaultFilter.data[c] == false {
			ns = append(ns, c)
		}
	}
	return string(ns)
}

func newEmoji() *emoji {
	file := os.TempDir() + "/emoji-data.txt"
	url := "https://unicode.org/Public/13.0.0/ucd/emoji/emoji-data.txt"

	t := new(emoji)
	// download emoji data file
	if err := t.downloadData(file, url); err != nil {
		log.Fatal(err)
	}

	// fetch emoji charaters
	emojis, err := t.fetchEmojis(file)
	if err != nil {
		log.Fatal(err)
	}
	t.data = emojis

	return t
}

func (t *emoji) downloadData(file, url string) error {
	if util.IsExist(file) {
		return nil
	}
	_, err := util.HTTPDownload(url, file)
	if err != nil {
		return err
	}
	return nil
}

func (t *emoji) fetchEmojis(file string) (map[rune]bool, error) {
	data, err := util.ReadFile(file)
	if err != nil {
		return nil, err
	}

	chars := make(map[rune]bool)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "#") == false || strings.Contains(line, ";") == false {
			continue
		}
		// 0023          ; Emoji                # E0.0   [1] (#️)       number sign
		// 002A          ; Emoji                # E0.0   [1] (*️)       asterisk
		// 0030..0039    ; Emoji                # E0.0  [10] (0️..9️)    digit zero..digit nine
		// 1F600         ; Emoji                # E1.0   [1] (😀)       grinning face
		// 1F9BA..1F9BF  ; Extended_Pictographic# E12.0  [6] (🦺..🦿)    safety vest..mechanical leg
		// 1F337..1F34A  ; Emoji                # E0.6  [20] (🌷..🍊)    tulip..tangerine
		line = line[:strings.IndexRune(line, ';')]
		columns := strings.Split(line, "..")
		if len(columns) == 1 {
			columns = append(columns, columns[0])
		}
		if len(columns) != 2 {
			continue
		}
		columns[0] = strings.TrimSpace(columns[0])
		columns[1] = strings.TrimSpace(columns[1])
		if columns[0] == "" || columns[1] == "" {
			continue
		}
		// skip number sign and digit zero..digit nine
		if columns[0] == "0023" && columns[1] == "0023" {
			continue
		}
		if columns[0] == "002A" && columns[1] == "002A" {
			continue
		}
		if columns[0] == "0030" && columns[1] == "0039" {
			continue
		}
		i1, _ := strconv.ParseInt("0x000"+columns[0], 0, 64)
		i2, _ := strconv.ParseInt("0x000"+columns[1], 0, 64)
		for i := i1; i <= i2; i++ {
			chars[t.decodeEmojiCode(fmt.Sprintf("%X", i))] = true
		}
	}

	return chars, nil
}

func (t *emoji) decodeEmojiCode(unicode string) rune {
	// Hex to Int
	i, err := strconv.ParseInt("0x000"+unicode, 0, 64)
	if err != nil {
		return 0
	}
	rs := []rune(string(i))
	return rs[0]
}
