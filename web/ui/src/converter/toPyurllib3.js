const toPyurllib3 = (url, method, header, body) => {
    let template = `import urllib3
client = urllib3.PoolManager()
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
	response = client.request('${method}', '${url}'`
    } else {
        template += `
	response = client.request('GET', '${url}'`
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
		body='${body}'`
    }

    template += `)
	print(response.status)
	print(response.data.decode())
except Exception as e:
	print(e)`
    return template
}


export {toPyurllib3}