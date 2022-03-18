package evasion

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-vgo/robotgo"
)

func GetNameEXE() string {
	return filepath.Base(os.Args[0])
}

func CheckNameEXE(name string) {

	println("[+] Checking the file name")
	if GetNameEXE() != name {
		println("[!] File name check failed.")
		os.Exit(1)
	}
}

func CheckMouse(sec int) {

	x, y := robotgo.GetMousePos()
	println("[+] Initial Mouse Position: ", x, y)

	println("[-] Waiting ", sec, " seconds... (Move the mouse)")
	time.Sleep(time.Duration(sec) * time.Second)

	j, z := robotgo.GetMousePos()
	fmt.Println("[+] Final Mouse Position: ", j, z)

	if x == j && y == z {
		println("[!] The mouse hasn't move")
		os.Exit(1)
	}
}

func CheckScreen() {

	width, height := robotgo.GetScreenSize()
	println("[+] Screen Size: ", width, height)

	if width < 500 && height < 500 {
		println("[!] Screen very small, suspicious")
		os.Exit(1)
	}
}

func CheckVMFilepath() bool {
	EvidenceOfSandbox := make([]string, 0)
	FilePathsToCheck := [...]string{`C:\windows\System32\Drivers\Vmmouse.sys`,
		`C:\windows\System32\Drivers\vm3dgl.dll`, `C:\windows\System32\Drivers\vmdum.dll`,
		`C:\windows\System32\Drivers\vm3dver.dll`, `C:\windows\System32\Drivers\vmtray.dll`,
		`C:\windows\System32\Drivers\vmci.sys`, `C:\windows\System32\Drivers\vmusbmouse.sys`,
		`C:\windows\System32\Drivers\vmx_svga.sys`, `C:\windows\System32\Drivers\vmxnet.sys`,
		`C:\windows\System32\Drivers\VMToolsHook.dll`, `C:\windows\System32\Drivers\vmhgfs.dll`,
		`C:\windows\System32\Drivers\vmmousever.dll`, `C:\windows\System32\Drivers\vmGuestLib.dll`,
		`C:\windows\System32\Drivers\VmGuestLibJava.dll`, `C:\windows\System32\Drivers\vmscsi.sys`,
		`C:\windows\System32\Drivers\VBoxMouse.sys`, `C:\windows\System32\Drivers\VBoxGuest.sys`,
		`C:\windows\System32\Drivers\VBoxSF.sys`, `C:\windows\System32\Drivers\VBoxVideo.sys`,
		`C:\windows\System32\vboxdisp.dll`, `C:\windows\System32\vboxhook.dll`,
		`C:\windows\System32\vboxmrxnp.dll`, `C:\windows\System32\vboxogl.dll`,
		`C:\windows\System32\vboxoglarrayspu.dll`, `C:\windows\System32\vboxoglcrutil.dll`,
		`C:\windows\System32\vboxoglerrorspu.dll`, `C:\windows\System32\vboxoglfeedbackspu.dll`,
		`C:\windows\System32\vboxoglpackspu.dll`, `C:\windows\System32\vboxoglpassthroughspu.dll`,
		`C:\windows\System32\vboxservice.exe`, `C:\windows\System32\vboxtray.exe`,
		`C:\windows\System32\VBoxControl.exe`}
	for _, FilePath := range FilePathsToCheck {
		if _, err := os.Stat(FilePath); err == nil {
			EvidenceOfSandbox = append(EvidenceOfSandbox, FilePath)
		}
	}
	if len(EvidenceOfSandbox) == 0 {
		return false
	} else {
		println("[!] VM Detected")
		return true
	}
}
