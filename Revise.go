package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

const (
	baseAddress = uintptr(0x6a9ec0) // 基址
	offset      = uintptr(0x768)    // 三级偏移
	SunOffset   = uintptr(0x5560)   // 阳光偏移
	cardOffset  = uintptr(0x144)    // 卡槽偏移
)

func coolingTimeClear(stopChan <-chan struct{}) {
	var baseBuffer [4]byte
	if err := readMemory(processHandle(), baseAddress, baseBuffer[:]); err != nil {
		fmt.Println("Error reading base address:", err)
		return
	}
	realBaseAddress := uintptr(*(*uint32)(unsafe.Pointer(&baseBuffer[0])))

	var offsetBuffer1 [4]byte
	if err := readMemory(processHandle(), realBaseAddress+offset, offsetBuffer1[:]); err != nil {
		fmt.Println("Error reading first offset:", err)
		return
	}
	realBaseAddress1 := uintptr(*(*uint32)(unsafe.Pointer(&offsetBuffer1[0])))

	var offsetBuffer2 [4]byte
	if err := readMemory(processHandle(), realBaseAddress1+cardOffset, offsetBuffer2[:]); err != nil {
		fmt.Println("Error reading first offset:", err)
		return
	}
	newSunValue := 1
	newSunBuffer := (*[4]byte)(unsafe.Pointer(&newSunValue))
	cards := []uintptr{0x70, 0xC0, 0x110, 0x160, 0x1B0, 0x200, 0x250, 0x2A0, 0x2F0, 0x340, 0x390, 0x3E0, 0x430, 0x480}

	for {
		select {
		case <-stopChan:
			// 收到停止信号，退出循环
			return
		default:
			for i, cards := range cards {
				finalAddress := uintptr(*(*uint32)(unsafe.Pointer(&offsetBuffer2[0]))) + cards
				if err := writeMemory(processHandle(), finalAddress, newSunBuffer[:]); err != nil {
					fmt.Printf("Error writing new sun value to final address %d: %v\n", i+1, err)
					// 如果出现错误，终止循环
					break
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func readSunValue() {
	var baseBuffer [4]byte
	if err := readMemory(processHandle(), baseAddress, baseBuffer[:]); err != nil {
		return
	}
	realBaseAddress := uintptr(*(*uint32)(unsafe.Pointer(&baseBuffer[0])))

	var offsetBuffer1 [4]byte
	if err := readMemory(processHandle(), realBaseAddress+offset, offsetBuffer1[:]); err != nil {
		return
	}
	finalAddress := uintptr(*(*uint32)(unsafe.Pointer(&offsetBuffer1[0]))) + SunOffset

	var sunValue [4]byte
	if err := readMemory(processHandle(), finalAddress, sunValue[:]); err != nil {
		return
	}

	getSunValue = int(*(*uint32)(unsafe.Pointer(&sunValue[0])))
}

func readSilverValue() {
	silverOffset := uintptr(0x82C)
	silverOffset2 := uintptr(0x208)
	var baseBuffer [4]byte
	if err := readMemory(processHandle(), baseAddress, baseBuffer[:]); err != nil {
		return
	}
	realBaseAddress := uintptr(*(*uint32)(unsafe.Pointer(&baseBuffer[0])))

	var offsetBuffer1 [4]byte
	if err := readMemory(processHandle(), realBaseAddress+silverOffset, offsetBuffer1[:]); err != nil {
		return
	}
	finalAddress := uintptr(*(*uint32)(unsafe.Pointer(&offsetBuffer1[0]))) + silverOffset2

	var sunValue [4]byte
	if err := readMemory(processHandle(), finalAddress, sunValue[:]); err != nil {
		return
	}

	getSilverValue = int(*(*uint32)(unsafe.Pointer(&sunValue[0])))
}

// 加农炮无CD
func cannonCd(selectState bool) {
	address := uintptr(0x0046103A)                             // 假设要写入的地址是 0x0046103A
	bytesToWrite := []byte{0x90, 0x90, 0x90, 0x90, 0x90, 0x90} // 用于替换的字节
	bytesToWrite1 := []byte{0x0F, 0x85, 0x9C, 0x01, 0x00, 0x00}

	// 写入内存
	if selectState {
		if err := writeMemory(processHandle(), address, bytesToWrite[:]); err != nil {
			return
		}
	} else {
		if err := writeMemory(processHandle(), address, bytesToWrite1[:]); err != nil {
			return
		}
	}
}

// 植物随机子弹
func randomBullet(selectState bool, stopChan <-chan struct{}) {
	address := uintptr(0x0046C769)
	randomBytesToWrite := []byte{0xC7, 0x45, 0x5C, 0x01, 0x00, 0x00, 0x00}
	renewBytesToWrite := []byte{0x89, 0x45, 0x5C, 0x8B, 0xC6, 0x90, 0x90}

	if selectState {
		for {
			select {
			case <-stopChan:
				// 收到停止信号，退出循环
				return
			default:
				randomValue := rand.Intn(29) + 1
				hexString := fmt.Sprintf("%02X", randomValue)
				byteValue, _ := hex.DecodeString(hexString)
				randomBytesToWrite[3] = byteValue[0]
				//fmt.Printf("%02X ", randomBytesToWrite)
				if err := writeMemory(processHandle(), address, randomBytesToWrite[:]); err != nil {
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	} else {
		//fmt.Println("关闭")
		if err := writeMemory(processHandle(), address, renewBytesToWrite[:]); err != nil {
			return
		}
	}
}

// 玉米投手锁定黄油
func lockButter(selectState bool) {
	address := uintptr(0x0045F1EC)
	BytesToWrite := []byte{0x90, 0x90}
	restBytesToWrite := []byte{0x75, 0x3F}
	if selectState {
		if err := writeMemory(processHandle(), address, BytesToWrite[:]); err != nil {
			return
		}
	} else {
		if err := writeMemory(processHandle(), address, restBytesToWrite[:]); err != nil {
			return
		}
	}

}

func lockPotato(selectState bool) {
	address := uintptr(0x0052FCF0)
	BytesToWrite := []byte{0xC7, 0x40, 0x54, 0x00, 0x00, 0x00, 0x00}
	restBytesToWrite := []byte{0xC7, 0x40, 0x54, 0xDC, 0x05, 0x00, 0x00}

	if selectState {
		if err := writeMemory(processHandle(), address, BytesToWrite[:]); err != nil {
			return
		}

	} else {
		if err := writeMemory(processHandle(), address, restBytesToWrite[:]); err != nil {
			return
		}

	}
}

func nutHpMax(selectState bool) {
	address1 := uintptr(0x0086EC6A)
	address2 := uintptr(0x0086EC76)
	BytesToWrite2 := []byte{0x83, 0x69, 0x40, 0x00, 0x90, 0x90, 0x90}
	restBytesToWrite2 := []byte{0x81, 0x69, 0x40, 0xE9, 0x03, 0x00, 0x00}

	if selectState {
		//if err := writeMemory(processHandle(), address1, BytesToWrite[:]); err != nil {
		//	return
		//}
		if err := writeMemory(processHandle(), address1, BytesToWrite2[:]); err != nil {
			return
		}
		if err := writeMemory(processHandle(), address2, BytesToWrite2[:]); err != nil {
			return
		}
	} else {
		//if err := writeMemory(processHandle(), address, restBytesToWrite[:]); err != nil {
		//	return
		//}
		if err := writeMemory(processHandle(), address1, restBytesToWrite2[:]); err != nil {
			return
		}
		if err := writeMemory(processHandle(), address2, restBytesToWrite2[:]); err != nil {
			return
		}
	}
}

func plantOverlap(selectState bool) {
	address := uintptr(0x0040FE2F)
	BytesToWrite := []byte{0xE9, 0x20, 0x09, 0x00, 0x00, 0x90}
	restBytesToWrite := []byte{0x0F, 0x84, 0x1F, 0x09, 0x00, 0x00}

	if selectState {
		if err := writeMemory(processHandle(), address, BytesToWrite[:]); err != nil {
			return
		}
	} else {
		if err := writeMemory(processHandle(), address, restBytesToWrite[:]); err != nil {
			return
		}
	}
}

func charmMushroom(selectState bool) {
	address := uintptr(0x004633FB)
	BytesToWrite := []byte{0x8B, 0x07, 0x90}
	restBytesToWrite := []byte{0x8B, 0x47, 0x50}

	if selectState {
		if err := writeMemory(processHandle(), address, BytesToWrite[:]); err != nil {
			return
		}
	} else {
		if err := writeMemory(processHandle(), address, restBytesToWrite[:]); err != nil {
			return
		}
	}
}
