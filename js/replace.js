// 移除不需要的元素 以及 替换复制不了的元素
// 移除div.related
function removeRelated() {
    document.querySelectorAll('div.related').forEach(div => {
        div.remove()
    });
}

// 移除页尾
function removeFooter() {
    document.querySelectorAll('div.footer').forEach(div => {
        div.remove()
    });
}

// 移除 H1
function removeH1() {
    document.querySelectorAll('h1').forEach(div => {
        div.remove()
    });
}

function removeSphinxsidebar() {
    document.querySelectorAll('div.sphinxsidebar').forEach(div => {
        div.remove()
    });
}

// 在em标签内容的左右两边加上反引号
function replaceEm() {
    document.querySelectorAll("em").forEach(em => {
        if (!['DT', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6'].includes(em.parentElement.tagName)) {
            const span = document.createElement("span")
            span.innerHTML = "\u0060" + em.innerHTML + "\u0060";
            em.insertAdjacentElement("afterend", span)
            em.remove()
        }
    })
}

// 在code标签内容的左右两边加上反引号
function replaceCode() {
    document.querySelectorAll("code").forEach(c => {
        if (!['A'].includes(c.parentElement.tagName)) {
            c.innerHTML = "\u0060" + c.innerHTML + "\u0060";
        }
    })
}

// 将div.versionadded 或 div.versionchanged 或div.admonition.note
// 或 div.admonition.seealso 或 div.admonition.caution
// 中的内容放在blockquote标签中
function replaceDivIntoBlockQuote() {
    document.querySelectorAll("div.versionadded,div.versionchanged,div.admonition").forEach(div => {
        const block = document.createElement("blockquote")
        while (div.firstChild) {
            block.appendChild(div.firstChild);
        }
        div.insertAdjacentElement("afterend", block)
        div.remove()
    })
}

// 将备注添加到其所指向的位置
function replaceRemarkToOriginPlace() {
    const footnoteListE = document.querySelector("aside.footnote-list")
    if (!footnoteListE) {
        return
    }
    const footnoteEs = footnoteListE.querySelectorAll("aside.footnote")
    footnoteEs.forEach(fn => {
        const content = fn.querySelector("p").textContent

        const aEs = fn.querySelectorAll(`span.label a[role="doc-backlink"]`)
        aEs.forEach(a => {
            const originId = a.href.split("#")[1].trim()
            console.log("originId=", originId)
            const xuHao = document.getElementById(originId).textContent
            const newContent = `（备注${xuHao}：${content}）`
            // 插入内容
            document.getElementById(originId).insertAdjacentText("afterend", newContent)
            document.getElementById(originId).remove()
        })
    })

    console.log(footnoteListE.previousElementSibling)
    if (footnoteListE.previousElementSibling.textContent.trim() === "备注") {
        footnoteListE.previousElementSibling.remove()
    }
    footnoteListE.remove()
}

// 去除代码块中 右上角的 >>>
function removeCopyButton() {
    document.querySelectorAll("span.copybutton").forEach(sp => {
        sp.remove()
    })
}

// 去除标题中固定链接（即一个悬浮图标）
function removeFixedHeaderLink() {
    document.querySelectorAll(".headerlink").forEach(a => {
        a.remove()
    })
}

// 增加添加标题的锚
function addHeaderAnchorAndRemoveHeaderLink() {
    const hs = document.querySelector("div.body").querySelectorAll("h2,h3,h4,h5,h6")

    hs.forEach(h => {
        const headerLinkE = h.querySelector("a.headerlink")
        const link = headerLinkE.href
        //去除后面的 ¶
        // const linkToTitle = h.textContent.replace(/¶+$/, '');
        const anchor = link.split("#")[1].trim()
        // h.setAttribute("data-href", link)
        headerLinkE.remove()
        h.insertAdjacentText("beforeend", `{#${anchor}}`)
    })
}

// 修改并替换dl.py.method 或 dl.py.attribute中的内容
function replaceDlPy() {
    const dlPyEs = document.querySelectorAll("dl.py")
    dlPyEs.forEach(dl => {
        //找到最近的标题是h3还是h4还是h5
        let prevE = dl.previousElementSibling;
        let foundH = false
        let curMustSetToHLevel = 4
        while (prevE && !foundH) {
            if (['H2', 'H3', 'H4', 'H5'].includes(prevE.tagName)) {
                foundH = true;
                curMustSetToHLevel = parseInt(prevE.tagName.replace("H", "")) + 1;
            } else {
                prevE = prevE.previousElementSibling;
            }
        }

        // 处理dl下的dd
        // 处理 dl 下的 dd
        const dd = dl.querySelector("dd");
        if (dd) {
            const divForDd = document.createElement("div");
            // 将 dd 的子节点移动到新的 div 容器中
            while (dd.firstChild) {
                divForDd.appendChild(dd.firstChild);
            }
            // 将新的 div 容器插入到 dl 后面
            dl.insertAdjacentElement("afterend", divForDd);
        }

        // 处理dl下的dt
        // let divForDtArr = []
        const fragment = document.createDocumentFragment();
        const dts = dl.querySelectorAll("dt")
        dts.forEach(dt => {
            const divForDt = document.createElement("div")
            const newH = document.createElement(`h${curMustSetToHLevel}`);
            let anchor = ""
            // 检查是否存在锚点
            if (dt.querySelector(".headerlink")) {
                const aE = dt.querySelector(".headerlink");
                anchor = aE.href.trim().split("#")[1];
                aE.remove(); // 移除原来的锚点链接
            }

            // 设置新标题的内容
            const title = dt.textContent.trim();
            newH.textContent = `\`${title}\`` + (anchor ? `{#${anchor}}` : "");
            divForDt.appendChild(newH)
            fragment.appendChild(divForDt);
            // divForDtArr.push(divForDt)
        })

        while (fragment.lastChild) {
            dl.insertAdjacentElement('afterend', fragment.lastChild);
        }
        // dl.insertAdjacentElement("afterend", fragment)
        // for(let i = divForDtArr.length - 1; i > 0; i--) {
        //     dl.insertAdjacentElement("afterend", divForDtArr[i])
        // }
        dl.remove()
    })
}

// 在p标签前面添加一个span，其内容为&zeroWidthSpace;用于后期在markdown中替换成Tab符号
function replaceP() {
    // 在p标签前面插入&zeroWidthSpace;
    document.querySelectorAll('p').forEach(function (p) {
        if (!['LI', "BLOCKQUOTE", "TH", "TD"].includes(p.parentElement.tagName)) {
            let newSpan = document.createElement('span');
            newSpan.textContent = '&zeroWidthSpace;';
            if (p.firstChild) {
                p.insertBefore(newSpan, p.firstChild);
            } else {
                // 如果 p 元素没有子节点，直接将新 span 元素添加到 p 元素中
                p.appendChild(newSpan);
            }
        }
    });
}


// span sig-name descname
removeRelated()
removeFooter();
removeH1();
removeSphinxsidebar();
replaceEm();
// replaceCode();
replaceDivIntoBlockQuote();
replaceRemarkToOriginPlace();
removeCopyButton();
addHeaderAnchorAndRemoveHeaderLink();
replaceDlPy();
replaceP();

