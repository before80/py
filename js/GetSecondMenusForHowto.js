() => {
    // // toctree-wrapper compound
    const uls = document.querySelectorAll("ul.simple")
    let menuInfos = []
    let exists = {}
    if (uls.length > 0) {
        uls.forEach(ul => {
            ul.querySelectorAll(":scope > li").forEach(li => {
                const a = li.querySelector('a')
                console.log("a=", a)
                const menu_name = a.textContent.trim()
                const url = a.href.trim()
                let names = url.split('/')
                names = names[names.length - 1].split("#")
                let filename = names[0].replace(/\.html$/, '')
                    .replace(/[\.\/]/g, '_')
                const noJhaoUrl = url.split("#")[0]
                console.log("noJhaoUrl=", noJhaoUrl)
                if (!exists[noJhaoUrl]) {
                    menuInfos.push({
                        menu_name: menu_name,
                        filename: filename,
                        url: noJhaoUrl,
                    })
                    exists[noJhaoUrl] = true
                }
            })
        })
    }
    console.log(menuInfos)
    return menuInfos
}

