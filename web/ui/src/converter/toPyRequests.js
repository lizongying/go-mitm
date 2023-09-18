const toPyRequests = (url, method, header, body) => {
    let template = `import requests
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
	response = requests.get('${url}'`
                break
            case 'POST':
                template += `
	response = requests.post('${url}'`
                break
            case 'HEAD':
                template += `
	response = requests.head('${url}'`
                break
            case 'PUT':
                template += `
	response = requests.put('${url}'`
                break
            case 'DELETE':
                template += `
	response = requests.delete('${url}'`
                break
            case 'PATCH':
                template += `
	response = requests.patch('${url}'`
                break
            case 'OPTIONS':
                template += `
	response = requests.options('${url}'`
                break
            default:
                template += `
	response = requests.request('${method}', '${url}'`
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


export {toPyRequests}