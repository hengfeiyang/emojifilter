package emojifilter

import (
	"testing"
)

func Test1(t *testing.T) {
	str := "1234567890my👮VANS album😎😜😞pudding by orange🍊eggs🥚Zoeva eyes💋shadow😈canmake"
	noEmojiStr := "1234567890myVANS albumpudding by orangeeggsZoeva eyesshadowcanmake"
	if noEmojiStr != Filter(str) {
		t.Errorf("emoji filter error: \n%s\n%s\n", noEmojiStr, Filter(str))
	}
}
