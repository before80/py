package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"py/bs"
	"py/lg"
	"py/pg"
	"strconv"
	"time"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("%v", err))
		}
	}()

	_ = err
	var browser *rod.Browser
	var page *rod.Page
	var browserHwnd win.HWND

	// 打开浏览器
	browser, err = bs.GetBrowser(strconv.Itoa(0))
	defer browser.MustClose()
	// 创建新页面
	page = browser.MustPage()
	browserHwnd = robotgo.GetHWND()
	var barMenuInfos []pg.BarMenuInfo
	barMenuInfos, err = pg.GetBarMenus(page, "https://docs.python.org/zh-cn/3.13/index.html")

	var secondMenuInfos []pg.MenuInfo
	for i, barMenuInfo := range barMenuInfos {
		lg.InfoToFileAndStdOut(fmt.Sprintf("bar正要处理 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))
		err = pg.InitIndexMdFile(i, barMenuInfo)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("bar初始化完成 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))

		secondMenuInfos, err = pg.InsertBarMenuPageData(browserHwnd, barMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("bar插入数据完成 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))

		if len(secondMenuInfos) <= 0 {
			continue
		}

		lg.InfoToFileAndStdOut(fmt.Sprintf("second处理二级菜单中 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))

		for j, menuInfo := range secondMenuInfos {
			err = pg.InitDetailPageMdFile(j, barMenuInfo, menuInfo)
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
				return
			}
			lg.InfoToFileAndStdOut(fmt.Sprintf("second初始化完成 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
				barMenuInfo.Filename, barMenuInfo.MenuName,
				menuInfo.Filename, menuInfo.MenuName))

			err = pg.InsertDetailPageData(browserHwnd, barMenuInfo, menuInfo, page)
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
				return
			}
			lg.InfoToFileAndStdOut(fmt.Sprintf("second插入数据完成 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
				barMenuInfo.Filename, barMenuInfo.MenuName,
				menuInfo.Filename, menuInfo.MenuName))

		}
	}
	time.Sleep(2000 * time.Second)
}
