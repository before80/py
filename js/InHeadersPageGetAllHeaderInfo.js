() => {
    // 获取所有头文件信息
    function getAllHeaderInfo() {
        const pattern = /<([^<>./]+)\.h>/;
        const table = document.querySelector('table.t-dsc-begin');
        const linkData = [];
        // 用于记录已经存在的 header
        const existingHeaders = {};
        if (table) {
            // 获取表格的所有行
            const rows = table.querySelectorAll('tr.t-dsc');
            rows.forEach(row => {
                // 获取第一列的所有链接
                const firstColumn = row.querySelector('td:first-child');
                const secondColumn = row.querySelector('td:nth-child(2)');
                const links = firstColumn.querySelectorAll('a');
                links.forEach(link => {
                    const fullHeader = link.textContent.trim();
                    const headerMatch  = fullHeader.match(pattern);
                    const url = link.href;
                    const header = headerMatch ? headerMatch[1] : "";
                    if (!existingHeaders[header]) {
                        linkData.push({
                            header: header,
                            fullHeader: fullHeader,
                            url: url,
                            desc: secondColumn.textContent.trim(),
                        });
                        // 标记该 header 已经存在
                        existingHeaders[header] = true;
                    }
                });
            });
            // 打印结果
            console.log(linkData);
        }
        return linkData
    }

    return getAllHeaderInfo()
}
