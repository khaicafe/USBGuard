package controller

import (
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	modkernel32           = syscall.NewLazyDLL("kernel32.dll")
	procReadProcessMemory = modkernel32.NewProc("ReadProcessMemory")
)

func readProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, buffer *byte, size uintptr, bytesRead *uint32) error {
	ret, _, err := procReadProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(buffer)),
		size,
		uintptr(unsafe.Pointer(bytesRead)),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func isDebuggerPresent() bool {
	isDebuggerPresent := modkernel32.NewProc("IsDebuggerPresent")
	ret, _, _ := isDebuggerPresent.Call()
	return ret != 0
}

func isRemoteDebuggerPresent() bool {
	checkRemoteDebugger := modkernel32.NewProc("CheckRemoteDebuggerPresent")
	var present int32
	hProcess, _ := syscall.GetCurrentProcess()
	checkRemoteDebugger.Call(uintptr(hProcess), uintptr(unsafe.Pointer(&present)))
	return present != 0
}

func isPEBDebugged() bool {
	type PROCESS_BASIC_INFORMATION struct {
		Reserved1       uintptr
		PebBaseAddress  uintptr
		Reserved2       [2]uintptr
		UniqueProcessId uintptr
		InheritedFrom   uintptr
	}

	ntdll := syscall.NewLazyDLL("ntdll.dll")
	ntQuery := ntdll.NewProc("NtQueryInformationProcess")

	var info PROCESS_BASIC_INFORMATION
	var returnLength uint32
	hProcess, _ := syscall.GetCurrentProcess()
	status, _, _ := ntQuery.Call(
		uintptr(hProcess),
		0,
		uintptr(unsafe.Pointer(&info)),
		unsafe.Sizeof(info),
		uintptr(unsafe.Pointer(&returnLength)),
	)
	if status != 0 {
		return false
	}

	var beingDebugged byte
	read := uint32(0)
	err := readProcessMemory(syscall.Handle(hProcess), info.PebBaseAddress+2, &beingDebugged, 1, &read)
	return err == nil && beingDebugged != 0
}

func isSleepSkipped() bool {
	start := time.Now()
	time.Sleep(200 * time.Millisecond)
	return time.Since(start) < 150*time.Millisecond
}

func hasBreakpointInstruction() bool {
	dummy := func() {}
	addr := uintptr(unsafe.Pointer(&dummy))
	var b [1]byte
	var read uint32
	handle, _ := syscall.GetCurrentProcess()
	err := readProcessMemory(handle, addr, &b[0], 1, &read)
	return err == nil && b[0] == 0xCC
}

func isApiHooked(dll string, funcName string) bool {
	mod := syscall.NewLazyDLL(dll)
	proc := mod.NewProc(funcName)
	addr := proc.Addr()
	var b [1]byte
	var read uint32
	handle, _ := syscall.GetCurrentProcess()
	err := readProcessMemory(handle, addr, &b[0], 1, &read)
	return err == nil && (b[0] == 0xE9 || b[0] == 0x68 || b[0] == 0xC3)
}

func isKnownDebuggerRunning() bool {
	out, err := exec.Command("tasklist").Output()
	if err != nil {
		return false
	}
	processList := strings.ToLower(string(out))
	known := []string{"x64dbg", "x32dbg", "ollydbg", "cheatengine", "ida", "scylla", "httpdebugger", "wireshark", "processhacker"}
	for _, name := range known {
		if strings.Contains(processList, name) {
			return true
		}
	}
	return false
}

func hasDebuggerWindow() bool {
	out, err := exec.Command("tasklist", "/v").Output()
	if err != nil {
		return false
	}
	s := strings.ToLower(string(out))
	return strings.Contains(s, "x64dbg") || strings.Contains(s, "cheat engine") || strings.Contains(s, "ollydbg")
}

func isCodeExecutionTampered() bool {
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
	elapsed := time.Since(start)
	return elapsed > 2*time.Second
}

func IsDebugged() bool {
	return isDebuggerPresent() ||
		isRemoteDebuggerPresent() ||
		isPEBDebugged() ||
		isSleepSkipped() ||
		hasBreakpointInstruction() ||
		isApiHooked("kernel32.dll", "IsDebuggerPresent") ||
		isApiHooked("kernel32.dll", "CheckRemoteDebuggerPresent") ||
		isKnownDebuggerRunning() ||
		hasDebuggerWindow() ||
		isCodeExecutionTampered()
}
