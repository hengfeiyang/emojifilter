# emojifilter
filter emoji from text by golang

# demo

```
package main

import (
	"fmt"

	"github.com/safeie/emojifilter"
)

func main() {
    str := "1234567890my👮VANS album😎😜😞pudding by orange🍊eggs🥚Zoeva eyes💋shadow😈canmake"
    // 1234567890myVANS albumpudding by orangeeggsZoeva eyesshadowcanmake
	noEmojiStr := emojifilter.Filter(str)
	fmt.Println(noEmojiStr)
}
```
