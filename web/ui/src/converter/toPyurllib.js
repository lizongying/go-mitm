const toPyurllib = (url, method, header, body) => {
    let template = `import urllib.request
`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
    } else {
        return ''
    }

    template += `
try:`
    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `
	request = urllib.request.Request('${url}', 
		method='${method}'`
    } else {
        template += `
	request = urllib.request.Request('${url}', 
		method='GET'`
    }

    if (header != null && typeof header === 'object') {
        template += `,
		headers={`
        Object.entries(header).forEach((v) => {
            template += `
			'${v[0]}': '${v[1]}',`
        })
        template += `
		}`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/'/g, '\\\'')
        template += `,
		data=b'${body}'`
    }

    template += `)
	response = urllib.request.urlopen(request)
	print(response.read().decode())
except Exception as e:
	print(e)`
    return template
}


export {toPyurllib}