const toCurl = (url, method, header, body) => {
    let template = `curl`

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += ` -X ${method}`
    }

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += ` '${url}' \\`
    } else {
        return ''
    }

    if (header != null && typeof header === 'object') {
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            template += `
-H '${v[0]}: ${value}' \\`
        })
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/'/g, '\\\'')
        template += `
--data '${body}' \\`
    }

    template += `
--compressed`
    return template
}


export {toCurl}