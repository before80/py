() => {
    // // toctree-wrapper compound
    const tocTreeUl = document.querySelector("div.toctree-wrapper.compound > ul")
    let menuInfos = []
    if (tocTreeUl) {
        tocTreeUl.querySelectorAll(":scope > li").forEach(li => {
            const a = li.querySelector(':scope > a')
            const menu_name = a.textContent.trim()
            const url = a.href.trim()
            const names = url.split('/')
            let filename = names[names.length - 1].replace(/\.html$/, '')
                .replace(/[\.\/]/g, '_')
            menuInfos.push({
                menu_name: menu_name,
                filename: filename,
                url: url,
            })
        })
    }
    console.log(menuInfos)
    return menuInfos
}

