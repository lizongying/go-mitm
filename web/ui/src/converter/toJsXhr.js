const toJsXhr = (url, method, header, body) => {
    let template = `const xhr = new XMLHttpRequest();`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `
xhr.open('${method}', '${url}');`
    } else {
        template += `
xhr.open('GET', '${url}');`
    }

    if (header != null && typeof header === 'object') {
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            template += `
xhr.setRequestHeader('${v[0]}', '${value}');`
        })
    }

    template += `
xhr.onload = () => {
	console.log(xhr.response);
};`

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/`/g, '\\`')
        template += `
xhr.send(\`${body}\`);`
    } else {
        template += `
xhr.send(null);`
    }

    return template
}


export {toJsXhr}