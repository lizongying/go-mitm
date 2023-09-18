const toJsJQuery = (url, method, header, body) => {
    let template = `$.ajax({`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += `
	url: '${url}'`
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `,`
        template += `
	method: '${method}'`
    }

    if (header != null && typeof header === 'object') {
        template += `,`
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
        template += `,`
        template += `
	data: \`${body}\``
    }

    template += `
}).done(response => {
	console.log(response);
});`
    return template
}


export {toJsJQuery}