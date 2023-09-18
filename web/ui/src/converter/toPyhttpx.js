const toPyhttpx = (url, method, header, body) => {
    let template = `import httpx
client = httpx.Client(http2=True, verify=False)
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
        switch (method) {
            case 'GET':
                template += `
	response = client.get('${url}'`
                break
            case 'POST':
                template += `
	response = client.post('${url}'`
                break
            case 'HEAD':
                template += `
	response = client.head('${url}'`
                break
            case 'PUT':
                template += `
	response = client.put('${url}'`
                break
            case 'DELETE':
                template += `
	response = client.delete('${url}'`
                break
            case 'PATCH':
                template += `
	response = client.patch('${url}'`
                break
            case 'OPTIONS':
                template += `
	response = client.options('${url}'`
                break
            default:
                template += `
	response = client.request('${method}', '${url}'`
        }
    } else {
        template += `
	response = requests.get('${url}'`
    }

    if (header != null && typeof header === 'object') {
        template += `,
		headers={`
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            template += `
			'${v[0]}': '${value}',`
        })
        template += `
		}`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/'/g, '\\\'')
        template += `,
		data='${body}'`
    }

    template += `)
	print(response.text)
except Exception as e:
	print(e)`
    return template
}


export {toPyhttpx}