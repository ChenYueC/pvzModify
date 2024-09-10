package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/windows"
	"image/color"
	"time"
)

var (
	_               fyne.Resource
	statusLabelGame = widget.NewLabel("状态：未检测到游戏")
	statusLabelSun  = widget.NewLabel("阳光：读取中")
	ProcessHandle   windows.Handle
	getSunValue     int
	getSilverValue  int
	//selectState           bool
	//randomSelectState     bool
	//lockButterSelectState bool
	//_ []byte
)

type myTheme struct {
	fyne.Theme
	font fyne.Resource
}

func (m *myTheme) Font(fyne.TextStyle) fyne.Resource {
	return m.font
}

func (m *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func main() {
	customFont := resourceAWanminTtf
	myApp := app.New()
	icon, _ := fyne.LoadResourceFromPath("Pvz.ico")
	myApp.SetIcon(icon)
	myWindow := myApp.NewWindow("Pvz 飞升之路")
	myApp.Settings().SetTheme(&myTheme{font: customFont})

	sunUpdate := widget.NewButton("阳光修改9990", func() {
		if ProcessHandle == 0 {
		} else {
			increaseSunValue()
			statusLabelSun.SetText("Tips:阳光已修改")
		}
	})

	clearTimeStop := make(chan struct{})
	clearTime := widget.NewCheck("植物卡槽种植无CD", func(checked bool) {
		if checked {
			go coolingTimeClear(clearTimeStop)
		} else {
			close(clearTimeStop)
			clearTimeStop = make(chan struct{})
		}
	})

	cannonSelectCheckCd := widget.NewCheck("毁灭加农炮发射无CD", func(checked bool) {
		if checked {
			cannonCd(true)
		} else {
			cannonCd(false)
		}
	})

	var lockButterCheck *widget.Check
	randomSelectCheckStop := make(chan struct{})
	randomSelectCheck := widget.NewCheck("全体植物-随机子弹", func(checked bool) {
		if checked {
			go randomBullet(true, randomSelectCheckStop)
			//randomSelectState = true
			lockButterCheck.Disable()
		} else {
			close(randomSelectCheckStop)
			//randomSelectState = false
			randomBullet(false, make(chan struct{}))
			randomSelectCheckStop = make(chan struct{})
			lockButterCheck.Enable()
		}
	})

	lockButterCheck = widget.NewCheck("玉米投手-锁定黄油", func(checked bool) {
		if checked {
			lockButter(true)
			//lockButterSelectState = true
			randomSelectCheck.Disable()
		} else {
			lockButter(false)
			//lockButterSelectState = false
			randomSelectCheck.Enable()
		}
	})

	potatoCd := widget.NewCheck("阳光土豆雷种植无CD", func(checked bool) {
		if checked {
			lockPotato(true)
		} else {
			lockPotato(false)
		}
	})

	hpMax := widget.NewCheck("冰帝Plus(免疫巨人、冰车不减血量)", func(checked bool) {
		if checked {
			nutHpMax(true)
		} else {
			nutHpMax(false)
		}
	})

	plantOverlapCheck := widget.NewCheck("植物叠加-无限种植", func(checked bool) {
		if checked {
			plantOverlap(true)
		} else {
			plantOverlap(false)
		}
	})

	charmMushroomTime := widget.NewCheck("魅惑菇射手-种植不消失", func(checked bool) {
		if checked {
			charmMushroom(true)
		} else {
			charmMushroom(false)
		}
	})

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if ProcessHandle == 0 {
				getProcessHandle := processHandle()
				ProcessHandle = getProcessHandle
				statusLabelGame.SetText("状态：未检测到游戏")
				statusLabelSun.SetText("阳光：" + "读取中")
			} else {
				statusLabelGame.SetText("状态：已检测到游戏")
				readSunValue()
				readSilverValue()
				statusLabelSun.SetText("阳光：" + fmt.Sprintf("%d", getSunValue))
				if getSunValue == 0 && getSilverValue == 0 {
					getProcessHandle := processHandle()
					ProcessHandle = getProcessHandle
					continue
				}
				potatoCd.Enable()
				hpMax.Enable()
				plantOverlapCheck.Enable()
			}
		}
	}()

	defer func() {
		err := windows.CloseHandle(ProcessHandle)
		if err != nil {
			fmt.Println("Failed to close handle:", err)
		} else {
			fmt.Println("Handle closed successfully")
		}
	}()

	content := container.NewHBox(statusLabelGame, statusLabelSun)
	contentStatusLabe := container.NewCenter(content)

	layout1 := container.NewHBox(
		randomSelectCheck,
		cannonSelectCheckCd,
	)

	layout2 := container.NewHBox(
		lockButterCheck,
		potatoCd,
	)

	layout3 := container.NewHBox(
		plantOverlapCheck,
		clearTime,
	)
	//layoutContainer = container.NewCenter(layoutContainer)

	myWindow.SetContent(container.NewVBox(
		contentStatusLabe,
		layout1, layout2, layout3, hpMax, charmMushroomTime,
		sunUpdate,
	))

	myWindow.Resize(fyne.NewSize(330, 239))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
