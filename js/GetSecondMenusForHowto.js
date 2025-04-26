() => {
    // // toctree-wrapper compound
    const uls = document.querySelectorAll("ul.simple")
    let menuInfos = []
    if (uls.length > 0) {
        uls.forEach(ul => {
            ul.querySelectorAll(":scope > li").forEach(li => {
                const a = li.querySelector(':scope > a')
                const menu_name = a.textContent.trim()
                const url = a.href.trim()
                let names = url.split('/')
                names = names[names.length - 1].split("#")
                let filename = names[0].replace(/\.html$/, '')
                    .replace(/[\.\/]/g, '_')
                menuInfos.push({
                    menu_name: menu_name,
                    filename: filename,
                    url: url,
                })
            })
        })
    }
    console.log(menuInfos)
    return menuInfos
}

