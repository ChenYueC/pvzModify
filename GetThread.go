package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	ProcessAllAccess = windows.PROCESS_QUERY_INFORMATION |
		windows.PROCESS_VM_OPERATION |
		windows.PROCESS_VM_READ |
		windows.PROCESS_VM_WRITE
)

func processHandle() windows.Handle {
	handle, err := getProcessHandle("PlantsVsZombies.exe")
	if err != nil {
		fmt.Println("获取进程句柄失败:", err)
		return 0
	}
	return handle
}

func getProcessHandle(processName string) (windows.Handle, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, nil
	}
	defer func(handle windows.Handle) {
		err := windows.CloseHandle(handle)
		if err != nil {

		}
	}(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snapshot, &entry); err != nil {
		return 0, nil
	}

	for {
		exeFileName := windows.UTF16ToString(entry.ExeFile[:])
		if exeFileName == processName {
			handle, err := windows.OpenProcess(ProcessAllAccess, false, entry.ProcessID)
			if err != nil {
				return 0, fmt.Errorf("无法打开进程: %w", err)
			}
			return handle, nil
		}
		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}

	return 0, nil
}

func readMemory(processHandle windows.Handle, address uintptr, buffer []byte) error {
	var bytesRead uintptr
	err := windows.ReadProcessMemory(processHandle, address, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		return nil
	}
	return nil
}

func writeMemory(processHandle windows.Handle, address uintptr, buffer []byte) error {
	var bytesWritten uintptr
	err := windows.WriteProcessMemory(processHandle, address, &buffer[0], uintptr(len(buffer)), &bytesWritten)
	if err != nil {
		return nil
	}
	return nil
}

func increaseSunValue() {
	processHandle := processHandle()
	var baseBuffer [4]byte
	if err := readMemory(processHandle, baseAddress, baseBuffer[:]); err != nil {
		return
	}
	realBaseAddress := uintptr(*(*uint32)(unsafe.Pointer(&baseBuffer[0])))

	var offsetBuffer1 [4]byte
	if err := readMemory(processHandle, realBaseAddress+offset, offsetBuffer1[:]); err != nil {
		return
	}
	finalAddress := uintptr(*(*uint32)(unsafe.Pointer(&offsetBuffer1[0]))) + SunOffset

	newSunValue := 9990
	newSunBuffer := (*[4]byte)(unsafe.Pointer(&newSunValue))

	if err := writeMemory(processHandle, finalAddress, newSunBuffer[:]); err != nil {
		return
	}

}
