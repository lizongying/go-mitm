const toJsFetch = (url, method, header, body) => {
    let template = `fetch(`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += `'${url}'`
    } else {
        return ''
    }

    let methodOk = false
    if (typeof method === 'string' && method !== '') {
        methodOk = true
        method = method.toUpperCase()
        template += `, {`
        template += `
	method: '${method}'`
    }

    let headerOk = false
    if (header != null && typeof header === 'object') {
        headerOk = true
        if (methodOk) {
            template += `,`
        } else {
            template += `, {`
        }

        template += `
	headers: {`
        let headers = []
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            headers.push(`'${v[0]}': '${value}'`)
        })
        template += `
		${headers.join(',\n		')}`
        template += `
	}`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/`/g, '\\`')
        if (methodOk || headerOk) {
            template += `,`
        } else {
            template += `, {`
        }
        template += `
	body: \`${body}\``
    }

    if (methodOk || headerOk || body) {
        template += `}`
    }
    template += `);`
    return template
}


export {toJsFetch}