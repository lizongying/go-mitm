const toJsAxios = (url, method, header, body) => {
    let template = `// ES6 module
import axios from 'axios';

// Node.js
// const axios = require('axios');
`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `,`
        switch (method) {
            case 'POST':
                template += `
axios.post('${url}'`
                break
            case 'GET':
                template += `
axios.get('${url}'`
                break
            default:
                break
        }
    } else {
        template += `
axios.get('${url}'`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/`/g, '\\`')
        template += `,`
        template += `
	\`${body}\``
    }

    if (header != null && typeof header === 'object') {
        template += `,`
        template += `
	{
		headers: {`
        let headers = []
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            headers.push(`'${v[0]}': '${value}'`)
        })
        template += `
		    ${headers.join(',\n		    ')}`
        template += `
		}
	}`
    }

    template += `)
	.then(response => {
		console.log('Response:', response.data);
	})
	.catch(error => {
		console.error('Error:', error);
	});`
    return template
}


export {toJsAxios}