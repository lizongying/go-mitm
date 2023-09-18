const toPyaiohttp = (url, method, header, body) => {
    let template = `import asyncio

import aiohttp


async def main():
    async with aiohttp.ClientSession() as session:`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        switch (method) {
            case 'GET':
                template += `
        async with session.get('${url}'`
                break
            case 'POST':
                template += `
        async with session.post('${url}'`
                break
            case 'HEAD':
                template += `
        async with session.head('${url}'`
                break
            case 'PUT':
                template += `
        async with session.put('${url}'`
                break
            case 'DELETE':
                template += `
        async with session.delete('${url}'`
                break
            case 'PATCH':
                template += `
        async with session.patch('${url}'`
                break
            case 'OPTIONS':
                template += `
        async with session.options('${url}'`
                break
            default:
                template += `
        async with session.request('${method}', '${url}'`
        }
    } else {
        template += `
        async with session.get('${url}'`
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

    template += `) as response:
		html = await response.text()
		print(html)


asyncio.run(main())`
    return template
}


export {toPyaiohttp}