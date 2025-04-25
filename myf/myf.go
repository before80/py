package myf

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TruncFileContent 清空指定文件中的内容
func TruncFileContent(filePath string) (err error) {
	var f *os.File
	// 以写入模式打开文件，若文件不存在则创建，若存在则清空内容
	f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("清空%s文件内容出现错误：%v", filePath, err)
	}
	// 确保文件在函数结束时关闭
	defer func() {
		_ = f.Close()
	}()
	return nil
}

// AddVersionInfoToMenu 添加版本信息到菜单中
func AddVersionInfoToMenu(menu string, code string) (versionInfo string) {
	// 从 code 中寻找版本信息
	// (C11 前) 替换成 -> 11 F
	// (C11 起)  替换成 -> 11+
	// (C11 弃用) 替换成 -> 11 D
	// (C11 移除) 替换成 -> 11 R
	// (C23 前) 替换成 -> 23 F
	// (C23 起)	替换成 -> 23+
	// (C23 弃用) 替换成 -> 23 D
	// (C23 移除) 替换成 -> 23 R

	return
}

// preLine 将所有以 >=2 个空白符开头的行，拼接到上一行。
// 保留其他行原样，最后用 "\n" 重新组装。
func preLine(src string) string {
	//fmt.Printf("src=%v\n\n", src)
	var out []string
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			// 认为这是上一行的续行，去掉前导空白，拼接
			if len(out) > 0 {
				out[len(out)-1] += " " + strings.TrimSpace(line)
				continue
			}
		}

		// 否则，新起一行
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// FindIdentifierVersion 查找标识符的版本号
func FindIdentifierVersion(identifier, src string) (string, error) {
	src = preLine(src)
	// 1) 找到含有 identifier 且带注释的那一行
	// (?m)：开启多行模式，让 ^ … $ 分别匹配每一行的行首/行尾
	//
	// \b ... \b：确保完整的单词匹配
	//
	// .*//.*：该行必须包含 // 注释
	lineRe := regexp.MustCompile(`(?m)^.*\b` + regexp.QuoteMeta(identifier) + `\b.*//.*$`)
	line := lineRe.FindString(src)
	if line == "" {
		return "", errors.New("identifier not found or no comment")
	}
	//fmt.Printf("line=%s\n", line)
	//fmt.Println()

	// 2) 拆出注释部分（"// ..."）
	idx := strings.Index(line, "//")
	if idx < 0 {
		return "", errors.New("no comment found on matched line")
	}
	comment := line[idx+2:] // 去掉 "//"
	//fmt.Printf("comment=%s\n", comment)
	//fmt.Println()

	// 3) 提取所有 "(…)"
	parenRe1 := regexp.MustCompile(`\([^)]*\)`)
	all := parenRe1.FindAllString(comment, -1)
	if len(all) == 0 {
		return "", errors.New("no parenthesized version found")
	}
	if len(all) > 1 {
		parenRe2 := regexp.MustCompile(`\(\d+\)`)
		if len(parenRe2.FindAllString(all[0], -1)) > 0 {
			return strings.TrimSpace(strings.Join(all[1:], " ")), nil
		}
	}
	return all[len(all)-1], nil
}

// ReplaceMarkdownFileContent 替换文件中的特定内容
func ReplaceMarkdownFileContent(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	replacements := []struct {
		pattern     *regexp.Regexp
		replacement string
	}{
		{regexp.MustCompile(`title\s*?=\s*?"([a-zA-Z]+)"`), "title = \"<$1.h>\""},
		{regexp.MustCompile("\n+```\n{2,}&zeroWidthSpace;"), "\n```\n\n&zeroWidthSpace;"},
		{regexp.MustCompile("输出：\\s*?\n```\\s*?\n"), "输出：\n\n```txt\n"},
		{regexp.MustCompile("```\\s*?//"), "```c\n//"},
		{regexp.MustCompile("```\\s*?\n#include"), "```c\n#include"},
		{regexp.MustCompile("```\\s*?\ntypedef"), "```c\ntypedef"},
		{regexp.MustCompile("```\\s*?\nvoid"), "```c\nvoid"},
		{regexp.MustCompile("```\\s*?\n#define"), "```c\n#define"},
		{regexp.MustCompile("```\\s*?\nchar"), "```c\nchar"},
		{regexp.MustCompile("```\\s*?\nint"), "```c\nint"},
		{regexp.MustCompile("```\\s*?\nfloat"), "```c\nfloat"},
		{regexp.MustCompile("```\\s*?\ndouble"), "```c\ndouble"},
		{regexp.MustCompile("```\\s*?\nstruct"), "```c\nstruct"},
		{regexp.MustCompile("```\\s*?\nunion"), "```c\nunion"},
		{regexp.MustCompile("@!br /!@"), "<br />"},
		{regexp.MustCompile("!@"), ">"},
		{regexp.MustCompile("@!"), "<"},
		{regexp.MustCompile("### 返回值"), "**返回值**"},
		{regexp.MustCompile("### 注意"), "**注意**"},
		{regexp.MustCompile("### 注解"), "**注解**"},
		{regexp.MustCompile("### 示例"), "**示例**"},
		{regexp.MustCompile("### 参数"), "**参数**"},
		{regexp.MustCompile("### 引用"), "**引用**"},
		{regexp.MustCompile("### 参阅"), "**参阅**"},
		{regexp.MustCompile("### 错误处理"), "**错误处理**"},
		{regexp.MustCompile("### 可能的实现"), "**可能的实现**"},
		{regexp.MustCompile("### 缺陷报告"), "**缺陷报告**"},
		{regexp.MustCompile("### 外部链接"), "**外部链接**"},
		{regexp.MustCompile("### 展开"), "**展开**"},
		{regexp.MustCompile("### 展开值"), "**展开值**"},
		{regexp.MustCompile("### 复杂度"), "**复杂度**"},
		{regexp.MustCompile("### 文件访问标记"), "**文件访问标记**"},
		{regexp.MustCompile("- &zeroWidthSpace; "), "  - "},
		{regexp.MustCompile("&zeroWidthSpace;"), "​\t"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?\(C(\d+)\s*?起\)`), "### $1 <- $2+"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?\(C(\d+)\s*?起\)`), "### $1 <- $2+"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?(\d{2}\+)\s*?\(C(\d{2})\s*?移除\)`), "### $1 <- $2 $3 R"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?(\d{2}\+)\s*?\(C(\d{2})\s*?弃用\)`), "### $1 <- $2 $3 D"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?(\d{2}\+)\s*?\(C(\d{2})\s*?前\)`), "### $1 <- $2 $3 F"},
		{regexp.MustCompile(`原址：([a-zA-Z0-9_:/?.#=&-]+)`), "原址：[$1]($1)"},
		{regexp.MustCompile(`运行此代码`), ""},
		{regexp.MustCompile("`\\*\\*A\\*\\*`"), "`A`"},
		{regexp.MustCompile("`\\*\\*a\\*\\*`"), "`a`"},
		{regexp.MustCompile("`\\*\\*c\\*\\*`"), "`c`"},
		{regexp.MustCompile("`\\*\\*d\\*\\*`"), "`d`"},
		{regexp.MustCompile("`\\*\\*F\\*\\*`"), "`F`"},
		{regexp.MustCompile("`\\*\\*f\\*\\*`"), "`f`"},
		{regexp.MustCompile("`\\*\\*E\\*\\*`"), "`E`"},
		{regexp.MustCompile("`\\*\\*e\\*\\*`"), "`e`"},
		{regexp.MustCompile("`\\*\\*G\\*\\*`"), "`G`"},
		{regexp.MustCompile("`\\*\\*g\\*\\*`"), "`g`"},
		{regexp.MustCompile("`\\*\\*i\\*\\*`"), "`i`"},
		{regexp.MustCompile("`\\*\\*n\\*\\*`"), "`n`"},
		{regexp.MustCompile("`\\*\\*P\\*\\*`"), "`P`"},
		{regexp.MustCompile("`\\*\\*O\\*\\*`"), "`O`"},
		{regexp.MustCompile("`\\*\\*o\\*\\*`"), "`o`"},
		{regexp.MustCompile("`\\*\\*p\\*\\*`"), "`p`"},
		{regexp.MustCompile("`\\*\\*s\\*\\*`"), "`s`"},
		{regexp.MustCompile("`\\*\\*U\\*\\*`"), "`U`"},
		{regexp.MustCompile("`\\*\\*u\\*\\*`"), "`u`"},
		{regexp.MustCompile("`\\*\\*X\\*\\*`"), "`X`"},
		{regexp.MustCompile("`\\*\\*x\\*\\*`"), "`x`"},
		{regexp.MustCompile("`\\*\\*Z\\*\\*`"), "`Z`"},
		{regexp.MustCompile("`\\*\\*z\\*\\*`"), "`z`"},
		{regexp.MustCompile("`\\*\\*\\+\\*\\*`"), "`+`"},
		{regexp.MustCompile("`\\*\\*-\\*\\*`"), "`-`"},
		{regexp.MustCompile("`\\*\\*%\\*\\*`"), "`%`"},
		{regexp.MustCompile("`\\*\\*\\*\\*\\*`"), "`*`"},
		{regexp.MustCompile("`\\*\\*\\.\\*\\*`"), "`.`"},
		{regexp.MustCompile("`\\*\\*#\\*\\*`"), "`#`"},
		{regexp.MustCompile("`\\*\\*0\\*\\*`"), "`0`"},
		{regexp.MustCompile("`\\*\\*0X\\*\\*`"), "`0X`"},
		{regexp.MustCompile("`\\*\\*0x\\*\\*`"), "`0x`"},
		{regexp.MustCompile("`\\*\\*%p\\*\\*`"), "`%p`"},
		{regexp.MustCompile("`\\*\\*%d\\*\\*`"), "`%d`"},
		{regexp.MustCompile("`\\*\\*%f\\*\\*`"), "`%f`"},
		{regexp.MustCompile("`\\*\\*%%\\*\\*`"), "`%%`"},
		{regexp.MustCompile(`'\*\*\\0\*\*'`), "'\\0'"},
		{regexp.MustCompile("`'\\*\\*\\f\\*\\*'`"), "`'\\f'`"},
		{regexp.MustCompile(`"\*\*\\f\*\*"`), "`\"\\f\"`"},
		{regexp.MustCompile(`'\*\*\\f\*\*'`), "`'\\f'`"},
		{regexp.MustCompile(`"\*\*\\n\*\*"`), "`\"\\n\"`"},
		{regexp.MustCompile(`'\*\*\\n\*\*'`), "`'\\n'`"},
		{regexp.MustCompile(`"\*\*\\r\*\*"`), "`\"\\r\"`"},
		{regexp.MustCompile(`'\*\*\\r\*\*'`), "`'\\r'`"},
		{regexp.MustCompile(`"\*\*\\t\*\*"`), "`\"\\t\"`"},
		{regexp.MustCompile(`'\*\*\\t\*\*'`), "`'\\t'`"},
		{regexp.MustCompile(`"\*\*\\v\*\*"`), "`\"\\v\"`"},
		{regexp.MustCompile(`'\*\*\\v\*\*'`), "`'\\v'`"},
		{regexp.MustCompile("`\\*\\*\\f\\*\\*`"), "`\\f`"},
		{regexp.MustCompile("`\\*\\*\\n\\*\\*`"), "`\\n`"},
		{regexp.MustCompile("`\\*\\*\\r\\*\\*`"), "`\\r`"},
		{regexp.MustCompile("`\\*\\*\\t\\*\\*`"), "`\\t`"},
		{regexp.MustCompile("`\\*\\*\\v\\*\\*`"), "`\\v`"},
		{regexp.MustCompile("`\\*\\*INF\\*\\*`"), "`INF`"},
		{regexp.MustCompile("`\\*\\*INFINITY\\*\\*`"), "`INFINITY`"},
		{regexp.MustCompile("`\\*\\*NAN\\*\\*`"), "`NAN`"},
		{regexp.MustCompile("\n+\\s*?```\n{6,}"), "\n```\n\n\n\n\n"},
	}

	modified := false
	newContent := string(content)
	for _, r := range replacements {
		if r.pattern.MatchString(newContent) {
			newContent = r.pattern.ReplaceAllString(newContent, r.replacement)
			if !modified {
				modified = true
			}
		}
	}

	if modified {
		err = os.WriteFile(filePath, []byte(newContent), 0666)
		//fmt.Println("1")
		if err != nil {
			return false, err
		}
		//fmt.Println("2")

	}

	// 按行判断是否有 ### 或 ## 开头的行，若有则替换这些行
	var totalLines []string
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalLines = append(totalLines, scanner.Text())
	}

	for i, line := range totalLines {
		if strings.HasPrefix(line, "### ") {
			totalLines[i] = strings.Replace(line, "### ", "", -1)
		}
		if strings.HasPrefix(line, "## ") {
			totalLines[i] = strings.Replace(line, "## ", "", -1)
		}
	}
	_ = file.Truncate(0)
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

	return modified, nil
}
