package pg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"io/fs"
	"os"
	"path/filepath"
	"py/cfg"
	"py/contants"
	"py/js"
	"py/lg"
	"py/myf"
	"py/wind"
	"slices"
	"strings"
	"time"
)

type HInfo struct {
	Header     string `json:"header"`
	FullHeader string `json:"fullHeader"`
	Url        string `json:"url"`
	Desc       string `json:"desc"`
}

type Hn2HInfo map[string]HInfo

// GetAllHeaderInfo 获取所有头信息
func GetAllHeaderInfo(page *rod.Page) (hn2HInfo Hn2HInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在 GetAllHeaderInfo 遇到panic：%v", r)
		}
	}()
	page.MustNavigate("https://zh.cppreference.com/w/c/header")
	page.MustWaitLoad()
	var result *proto.RuntimeRemoteObject

	result, err = page.Eval(js.InHeadersPageGetAllHeaderInfoJs)

	if err != nil {
		return nil, err
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", err)
	}

	// 定义一个结构体切片来解析 JSON 数据
	var data []HInfo

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %v", err)
	}
	hn2HInfo = make(Hn2HInfo)
	for _, st := range data {
		hn2HInfo[st.Header] = st
	}

	return hn2HInfo, nil
}

// findSectionStart 查找指定类型的起始行
func findSectionStart(lines []string, sectionType string) int {
	sectionHeader := ""
	switch sectionType {
	case "f":
		sectionHeader = "## 函数"
	case "m":
		sectionHeader = "## 宏"
	case "t":
		sectionHeader = "## 类型"
	case "e":
		sectionHeader = "## 枚举"
	default:
		return -1
	}
	for i, line := range lines {
		if strings.TrimSpace(line) == sectionHeader {
			return i
		}
	}
	return -1
}

type IdInfo struct {
	Id     string `json:"id"`     // 标识符
	Typ    string `json:"typ"`    // 类型：f=函数，m=宏，e=枚举，t=类型
	Url    string `json:"url"`    // 所在网址
	Remark string `json:"remark"` // 备注，特殊
	Desc   string `json:"desc"`   // 描述，通常在表格中
}

// GetSomeoneHeaderAllIdentifierInfo 获取某一头的所欲标识符信息
func GetSomeoneHeaderAllIdentifierInfo(hInfo HInfo, page *rod.Page) (idInfos []IdInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("GetSomeoneHeaderAllIdentifierInfo 在处理%s时遇到panic：%v", hInfo.Url, r)
		}
	}()
	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(js.InSomeoneHeaderIntroPageGetAllIdentifierInfoJs)

	if err != nil {
		return nil, err
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", err)
	}

	// 定义一个结构体切片来解析 JSON 数据
	var data []IdInfo

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %v", err)
	}
	//fmt.Println("run 1")

	if len(data) > 0 {
		//fmt.Println("run 2")
		// 根据typ将不同的标识符进行归类
		var id2IdInfo = make(map[string]IdInfo)
		var tIds, eIds, mIds, fIds []string
		for _, idInfo := range data {
			id2IdInfo[idInfo.Id] = idInfo
			if idInfo.Typ == "t" {
				tIds = append(tIds, idInfo.Id)
			}
			if idInfo.Typ == "e" {
				eIds = append(eIds, idInfo.Id)
			}
			if idInfo.Typ == "m" {
				mIds = append(mIds, idInfo.Id)
			}
			if idInfo.Typ == "f" {
				fIds = append(fIds, idInfo.Id)
			}
		}
		//fmt.Println("run 3")
		//fmt.Println("id2IdInfo=", id2IdInfo)

		slices.Sort(tIds)
		slices.Sort(eIds)
		slices.Sort(mIds)
		slices.Sort(fIds)
		lg.InfoToFileAndStdOut(fmt.Sprintf("hInfo.Header=%v\n", hInfo.Header))
		lg.InfoToFileAndStdOut(fmt.Sprintf("tIds=%v\n", tIds))
		lg.InfoToFileAndStdOut(fmt.Sprintf("eIds=%v\n", eIds))
		lg.InfoToFileAndStdOut(fmt.Sprintf("mIds=%v\n", mIds))
		lg.InfoToFileAndStdOut(fmt.Sprintf("fIds=%v\n", fIds))
		fp := filepath.Join(contants.OutputFolderName, contants.CStdFolderName, hInfo.Header+".md")

		if len(tIds) > 0 {
			err = insertSubmenus("t", tIds, id2IdInfo, fp)
			if err != nil {
				return data, fmt.Errorf("头%s在插入类型为类型的标识符对应的数据时出现错误：%v", hInfo.Header, err)
			}
		}

		if len(eIds) > 0 {
			err = insertSubmenus("e", eIds, id2IdInfo, fp)
			if err != nil {
				return data, fmt.Errorf("头%s在插入类型为枚举的标识符对应的数据时出现错误：%v", hInfo.Header, err)
			}
		}

		if len(mIds) > 0 {
			err = insertSubmenus("m", mIds, id2IdInfo, fp)
			if err != nil {
				return data, fmt.Errorf("头%s在插入类型为宏的标识符对应的数据时出现错误：%v", hInfo.Header, err)
			}
		}

		if len(fIds) > 0 {
			err = insertSubmenus("f", fIds, id2IdInfo, fp)
			if err != nil {
				return data, fmt.Errorf("头%s在插入类型为函数的标识符对应的数据时出现错误：%v", hInfo.Header, err)
			}
		}
	}

	return data, nil
}

// insertSubmenus 插入菜单，其中ids必须已经按照升序进行排序
func insertSubmenus(typ string, ids []string, id2IdInfo map[string]IdInfo, fp string) (err error) {
	var file *os.File
	file, err = os.OpenFile(fp, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("打开文件 %s 时出错: %v\n", fp, err)
	}
	defer file.Close()

	var totalLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalLines = append(totalLines, scanner.Text())
	}

	start := findSectionStart(totalLines, typ)
	if start == -1 {
		return fmt.Errorf("未找到类型为 %s 的起始行", typ)
	}

	var url string
	var newTotalLines []string
	insertIndex := start + 1
	hadExisted := false
	oldTotalL := len(totalLines)
	for _, id := range ids {
		idInfo, ok := id2IdInfo[id]
		if !ok {
			continue
		}

		hadExisted = false
		oldTotalL = len(totalLines)
		for i := insertIndex; i < oldTotalL; i++ {
			line := totalLines[i]
			if (i + 1) == oldTotalL {
				insertIndex = i
				break
			}

			if !strings.HasPrefix(line, "### ") && !strings.HasPrefix(line, "## ") {
				continue
			}

			if strings.HasPrefix(line, "### ") {
				line = strings.TrimPrefix(line, "### ")
				sps := strings.Split(line, "<-")
				if len(sps) > 1 {
					line = sps[0]
				}
				line = strings.TrimSpace(line)
				if strings.HasPrefix(id, "__") {
					line = strings.Trim(line, "`")
				}

				if id == line {
					hadExisted = true
					break
				}

				if id < line {
					insertIndex = i - 1
					break
				}
			} else if strings.HasPrefix(line, "## ") {
				insertIndex = i - 1
				break
			}
		}

		if !hadExisted {
			newTotalLines = []string{}
			newLines := make([]string, 0)
			if strings.HasPrefix(id, "__") {
				newLines = append(newLines, fmt.Sprintf("### `%s`\n", id))
			} else {
				newLines = append(newLines, fmt.Sprintf("### %s\n", id))
			}

			url = strings.TrimSpace(idInfo.Url)
			if url != "" {
				url = fmt.Sprintf(`[%s](%s)`, url, url)
			}

			newLines = append(newLines, fmt.Sprintf("原址：%s\n", url))
			newLines = append(newLines, fmt.Sprintf("作用：%s\n", idInfo.Desc))
			newLines = append(newLines, fmt.Sprintf("备注：%s\n", idInfo.Remark))
			newLines = append(newLines, fmt.Sprintf("\n"))
			newLines = append(newLines, fmt.Sprintf("\n"))

			if insertIndex <= oldTotalL {
				newTotalLines = append(newTotalLines, totalLines[:insertIndex]...)
			} else {
				newTotalLines = append(newTotalLines, totalLines[:insertIndex-1]...)
			}
			newTotalLines = append(newTotalLines, newLines...)

			if insertIndex < oldTotalL {
				newTotalLines = append(newTotalLines, totalLines[insertIndex:]...)
			}
			totalLines = newTotalLines
		}
	}

	_, _ = file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range totalLines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	return
}

func getFileLines(fp string) (lines []string, err error) {
	var file *os.File
	file, err = os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("打开文件 %s 时出错: %v\n", fp, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return
}

// InitSpecialHeaderMdFile 初始化指定头md文件
func InitSpecialHeaderMdFile(index int, hInfo HInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("InitSpecialHeaderMdFile 在处理%s时遇到panic：%v", hInfo.Url, r)
		}
	}()
	// 打开各个头文件的介绍页，获取该头文件包含哪些标识符（类型、函数、宏、枚举）
	page.MustNavigate(hInfo.Url)
	//page.MustNavigate("https://en.cppreference.com/w/c/header/stdckdint")
	page.MustWaitLoad()

	// 获取h1标签的内容
	h1Content := page.MustElement("#firstHeading").MustText()
	h1Content = strings.TrimSpace(strings.Replace(h1Content, "标准库标头", "", -1))

	err = os.MkdirAll(filepath.Join(contants.OutputFolderName, contants.CStdFolderName), 0777)
	if err != nil {
		return fmt.Errorf("无法创建%s目录：%v\n", filepath.Join(contants.OutputFolderName, contants.CStdFolderName), err)
	}

	newMdFp := filepath.Join(contants.OutputFolderName, contants.CStdFolderName, hInfo.Header+".md")
	var newMd *os.File
	_, err1 := os.Stat(newMdFp)
	// 当文件不存在的情况下，新建文件并初始化该文件
	if err1 != nil && errors.Is(err1, fs.ErrNotExist) {
		//fmt.Println("err=", err1)
		newMd, err = os.OpenFile(newMdFp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("创建文件 %s 时出错: %w", newMdFp, err)
		}
		defer newMd.Close()
		_, err = newMd.WriteString(fmt.Sprintf(`
+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++

## 类型




## 枚举




## 宏




## 函数




`, h1Content, time.Now().Format(time.RFC3339), index*10, hInfo.Desc))
		if err != nil {
			return fmt.Errorf("写入%s文件时出错: %v", newMdFp, err)
		}
	}
	return nil
}

func GetIdentifierData(browserHwnd win.HWND, hn string, idInfo IdInfo, page *rod.Page) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("在 GetIdentifierData 遇到panic：%v", r)
		}
	}()
	lg.InfoToFileAndStdOut(fmt.Sprintf("idInfo.Url=%s\n", idInfo.Url))
	//page.MustNavigate(strings.Replace(idInfo.Url, "https://en.cppreference.com", "https://zh.cppreference.com", -1))
	page.MustNavigate(idInfo.Url)
	page.MustWaitLoad()

	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(fmt.Sprintf(`() => { %s }`, js.InIdentifierPageRemoveAndReplaceJs))
	_ = result
	if err != nil {
		return fmt.Errorf("在网页%s中执行InIdentifierPageRemoveAndReplaceJs遇到错误：%v", idInfo.Url, err)
	}

	uniqueMdFilepath := cfg.Default.UniqueMdFilepath
	// 获取文件名
	spSlice := strings.Split(uniqueMdFilepath, "\\")
	mdFilename := spSlice[len(spSlice)-1]

	// 清空唯一共用的markdown文件的文件内容
	err = myf.TruncFileContent(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理头%s中的%s的数据时，清空%q文件内容出现错误：%v", hn, idInfo.Id, uniqueMdFilepath, err)
	}

	// 打开 唯一共用的markdown文件
	err = wind.OpenTypora(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理头%s中的%s的数据时，打开窗口名为%q时出现错误：%v", hn, idInfo.Id, uniqueMdFilepath, err)
	}
	// 适当延时保证能打开 typora
	time.Sleep(3 * time.Second)

	var typoraHwnd win.HWND
	typoraWindowName := fmt.Sprintf("%s - Typora", mdFilename)
	typoraHwnd, err = wind.FindWindowHwndByWindowTitle(typoraWindowName)
	if err != nil {
		return fmt.Errorf("在处理头%s中的%s的数据时，找不到%q窗口：%v", hn, idInfo.Id, typoraWindowName, err)
	}

	wind.SelectAllAndCtrlC(browserHwnd)
	wind.SelectAllAndDelete(typoraHwnd)
	wind.CtrlV(typoraHwnd)

	wind.CtrlS(typoraHwnd)
	robotgo.CloseWindow()
	time.Sleep(1 * time.Second)
	_, err = myf.ReplaceMarkdownFileContent(uniqueMdFilepath)
	if err != nil {
		return fmt.Errorf("在处理头%s中的%s的数据时，替换出现错误：%v", hn, idInfo.Id, err)
	}

	fpDst := filepath.Join(contants.OutputFolderName, contants.CStdFolderName, hn+".md")

	//将md的内容放到ids对应的三级菜单下
	return insertIdentifierData(idInfo, uniqueMdFilepath, fpDst)
}

func insertIdentifierData(idInfo IdInfo, fpSrc, fpDst string) (err error) {
	var dstFile, srcFile *os.File
	dstFile, err = os.OpenFile(fpDst, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("打开文件 %s 时出错: %v\n", fpDst, err)
	}
	defer dstFile.Close()

	var dstFileTotalLines []string
	scanner := bufio.NewScanner(dstFile)
	for scanner.Scan() {
		dstFileTotalLines = append(dstFileTotalLines, scanner.Text())
	}

	start := findSectionStart(dstFileTotalLines, idInfo.Typ)
	if start == -1 {
		return fmt.Errorf("未找到类型为 %s 的起始行", idInfo.Typ)
	}

	insertIndex := start + 1
	hadExisted := false
	dstFileLineLen := len(dstFileTotalLines)
	for i := start + 1; i < dstFileLineLen; i++ {
		line := dstFileTotalLines[i]
		if !strings.HasPrefix(line, "### ") && !strings.HasPrefix(line, "## ") {
			continue
		}
		if strings.HasPrefix(line, "### ") {
			line = strings.TrimPrefix(line, "### ")
			sps := strings.Split(line, "<-")
			if len(sps) > 1 {
				line = sps[0]
			}
			line = strings.TrimSpace(line)

			if idInfo.Id == line {
				hadExisted = true
				insertIndex = i + 7 // +7 表示找到备注所在行之后
				break
			} else if idInfo.Id < line {
				insertIndex = i
				break
			}
			continue
		} else if strings.HasPrefix(line, "## ") {
			insertIndex = i - 1
			break
		}
	}

	if !hadExisted {
		return fmt.Errorf("在文件%s中未找到%s", fpDst, idInfo.Id)
	}

	// 读取uniqueMd文件中的内容
	srcFile, err = os.OpenFile(fpSrc, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("打开文件 %s 时出错: %v\n", fpSrc, err)
	}
	defer srcFile.Close()

	var srcFileTotalLines []string
	scanner = bufio.NewScanner(srcFile)
	for scanner.Scan() {
		srcFileTotalLines = append(srcFileTotalLines, scanner.Text())
	}

	var newTotalLines []string

	if insertIndex <= dstFileLineLen {
		newTotalLines = append(newTotalLines, dstFileTotalLines[:insertIndex]...)
	} else {
		newTotalLines = append(newTotalLines, dstFileTotalLines[:insertIndex-1]...)
	}
	newTotalLines = append(newTotalLines, srcFileTotalLines...)

	if insertIndex < dstFileLineLen {
		newTotalLines = append(newTotalLines, dstFileTotalLines[insertIndex:]...)
	}

	_, _ = dstFile.Seek(0, 0)
	writer := bufio.NewWriter(dstFile)
	for _, line := range newTotalLines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	return nil
}
