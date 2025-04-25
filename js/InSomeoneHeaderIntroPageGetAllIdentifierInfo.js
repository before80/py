() => {

    function replaceSpanSourceC() {
        document.querySelectorAll('span.mw-geshi.c.source-c').forEach(span => {
            if (!span.querySelector(".kw486")) {
                // 在原有内容前后添加反引号（保留内部HTML结构）
                span.innerHTML = "\u0060" + span.innerHTML + "\u0060";
            }
        });
    }

    replaceSpanSourceC();

// 根据表格第二列中提示，返回类型标识
    function calTypeFlag(str) {
        // (函数) -> f <-函数
        // (常量) (枚举) -> e <- 枚举
        // (宏常量) (关键词宏) (宏函数) -> m <- 宏
        // (typedef) (结构体) -> t <- 类型
        // 没有的情况 -> 暂时归类为 t <- 类型
        if(str === "(函数)") {
            return "f"
        }

        if(str === "(常量)" || str === "(枚举)") {
            return "e"
        }

        if(str === "(宏常量)" || str === "(关键词宏)"|| str === "(宏函数)") {
            return "m"
        }

        if(str === "(typedef)" || str === "(结构体)") {
            return "t"
        }

        return "t"
    }


    function GetIdentifierFromGaiYao() {
        let ids = []
        let exist = {}

        const divMwGeshis = document.querySelectorAll("div.mw-geshi")
        divMwGeshis.forEach(mw => {
            let typeNames = [];
            let macroNames = [];
            let functionNames = [];
            // 判断其紧邻的元素中的上一个p标签是否有如下内容作为开头：仅当实现定义了
            let remark = ""
            const prevE = mw.previousElementSibling;
            console.log("prevE.textContent=", prevE.textContent)
            if (prevE.tagName === "P" && prevE.textContent.startsWith("仅当实现定义了")) {
                remark = prevE.textContent
            }

            const pre = mw.querySelector("div.c.source-c > pre.de1")
            const preCode = pre.textContent
            const lines = preCode.split('\n');
            lines.forEach(line => {
                // 去除注释
                line = line.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, '');
                line = line.trim();

                // 匹配类型名称
                const typeMatch = line.match(/^typedef\s+(?:struct|union|enum)?\s*(\w+);/);
                if (typeMatch) {
                    typeNames.push(typeMatch[1]);
                }

                // 匹配宏名称
                const macroMatch = line.match(/^#define\s+(\w+)\b/);
                if (macroMatch) {
                    macroNames.push(macroMatch[1]);
                }

                // 匹配函数名
                const functionMatch = line.match(/^([\w\s*_]+)\s+(\w+)\s*\(/);
                if (functionMatch) {
                    // 排除类型定义、宏定义和关键字开头的情况
                    if (!/^(typedef|#define|struct|union|enum)/.test(line)) {
                        functionNames.push(functionMatch[2]);
                    }
                }
            });
            // console.log(1,typeNames)
            // console.log(1,macroNames)
            // console.log(1,functionNames)
            // console.log("---------------------------")
            typeNames = [...new Set(typeNames)]
            macroNames = [...new Set(macroNames)]
            functionNames = [...new Set(functionNames)]
            // console.log(2,typeNames)
            // console.log(2,macroNames)
            // console.log(2,functionNames)
            if (typeNames.length > 0) {
                typeNames.forEach(id => {
                    exist[id] = true
                    ids.push({
                        id: id,
                        typ: "t",
                        url: "",
                        remark: remark,
                        desc: "",
                    })
                })
            }

            if (macroNames.length > 0) {
                macroNames.forEach(id => {
                    exist[id] = true
                    ids.push({
                        id: id,
                        typ: "m",
                        url: "",
                        remark: remark,
                        desc: "",
                    })
                })
            }

            if (functionNames.length > 0) {
                functionNames.forEach(id => {
                    exist[id] = true
                    ids.push({
                        id: id,
                        typ: "f",
                        url: "",
                        remark: remark,
                        desc: "",
                    })
                })
            }

            const aLinks = mw.querySelectorAll("div.c.source-c > pre.de1 a")
            aLinks.forEach(a => {
                for(const obj of ids) {
                    if(obj.id === a.textContent.trim()) {
                        obj.url = a.href
                    }
                }
            })
        })

        console.log(ids)

        const trs = document.querySelectorAll("table.t-dsc-begin tr.t-dsc")
        trs.forEach(tr => {
            const firstTd = tr.querySelector("td:first-child")
            let secondTd = tr.querySelector("td:nth-child(2)").cloneNode(true)
            if (secondTd.querySelector(".editsection")) {
                secondTd.removeChild(secondTd.querySelector(".editsection"))
            }
            // if (secondTd.querySelector("br")) {
            //     secondTd.removeChild(secondTd.querySelector("br"))
            // }

            const aE = firstTd.querySelector("a")
            const spans = firstTd.querySelectorAll("div.t-dsc-member-div div:first-child span.t-lines span")
            const typFlagSpan = tr.querySelector("td:nth-child(2) span.t-mark")
            // console.log("secondTd.textContent.trim()=", secondTd.textContent.trim())
            if (spans) {
                spans.forEach(span => {
                    const id = span.textContent.trim()
                    if(!exist[id]) {
                        ids.push({
                            id: id,
                            typ: calTypeFlag(typFlagSpan.textContent.trim()),
                            url: aE?.href ? aE.href : "",
                            remark: "",
                            desc: secondTd.textContent.trim(),
                        })
                        exist[id] = true
                    } else {
                        for(const obj of ids) {
                            if (obj.id === id) {
                                obj.url = aE?.href ? aE.href : ""
                                obj.desc = secondTd.textContent.trim()
                            }
                        }
                    }
                })
            } else {
                const codeB = firstTd.querySelector("code > b")
                const id = codeB.textContent.trim()
                if(!exist[id]) {
                    ids.push({
                        id: id,
                        typ: calTypeFlag(typFlagSpan.textContent.trim()),
                        url: aE?.href ? aE.href : "",
                        remark: "",
                        desc: secondTd.textContent.trim(),
                    })
                    exist[id] = true
                } else {
                    for(const obj of ids) {
                        if (obj.id === id) {
                            obj.url = aE?.href ? aE.href : ""
                            obj.desc = secondTd.textContent.trim()
                        }
                    }
                }
            }
        })
        console.log("ids=", ids)
        return ids
    }

    return GetIdentifierFromGaiYao()
}