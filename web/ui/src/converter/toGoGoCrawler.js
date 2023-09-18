const toGoGoCrawler = (url, method, header, body) => {
    let template = `package main

import (
	"github.com/lizongying/go-crawler/pkg"
	"github.com/lizongying/go-crawler/pkg/request"
)

func main() {
	_ = request.NewRequest()`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += `.
		SetUrl("${url}")`
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `.
		SetMethod(pkg.${method})`
    }

    if (header != null && typeof header === 'object') {
        template += `.
		SetHeaders(map[string]string{`
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/"/g, '\\"')
            template += `
			"${v[0]}": "${value}",`
        })
        template += `
		})`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/`/g, '\\`')
        template += `.
		SetBodyStr(\`${body}\`)`
    }

    template += `
}`

    return template
}


export {toGoGoCrawler}