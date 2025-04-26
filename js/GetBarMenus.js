() => {
    const url = "%s"
    let baseUrl = url.replace(/\/index\.html$/, '');
    baseUrl = baseUrl.replace(/\/$/, '');
    baseUrl = baseUrl + '/'
    let menuInfos = []

    document.querySelectorAll("table.contentstable").forEach((t, i) => {
        const ps = t.querySelectorAll("p")
        if (ps.length > 0) {
            ps.forEach(p => {
                const a = p.querySelector('a')
                const menu_name = a.textContent.trim()
                const url = a.href.trim()
                let filename = url.replace(baseUrl, '')
                    .replace(/\/index\.html$/, '')
                    .replace(/\.html$/, '')
                    .replace(/[\.\/]/g, '_')
                if (i === 0 || (i > 0 && ["术语对照表", "Python 的历史与许可证"].includes(menu_name))) {
                    menuInfos.push({
                        menu_name: menu_name,
                        filename: filename,
                        url: url,
                    })
                }
            })
        }
    })
    console.log(menuInfos)
    return menuInfos
}
