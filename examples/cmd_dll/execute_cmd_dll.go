package cmd_dll

import "C"
import (
	"fmt"
	"os/exec"
)

//export PopCalc
func PopCalc() {
	fmt.Println("Spawning calculator")
	cmd := exec.Command("cmd.exe", "/C", "C:\\Windows\\System32\\calc.exe")
	cmd.Run()
}

func main() {
	// Blank
}
