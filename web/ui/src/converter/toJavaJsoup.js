const toJavaJsoup = (url, method, header, body) => {
		let template = `package org.example;

import org.jsoup.Connection;
import org.jsoup.Jsoup;

import java.io.IOException;

public class JsoupExample {
	public static void main(String[] args) throws IOException {
`

		if (typeof url === 'string' && url !== '') {
				url = encodeURI(url)
				template+=`
		Connection.Response response = Jsoup.connect("${url}")`
		} else {
				return ''
		}

		if (typeof method === 'string' && method !== '') {
				method = method.toUpperCase()
				template += `
			.method(Connection.Method.${method})`
		}

		if (header != null && typeof header === 'object') {
				Object.entries(header).forEach((v) => {
					const value = v[1].replaceAll(/"/g, '\\"')
						template += `
			.header("${v[0]}", "${value}")`
				})
		}

		if (typeof body === 'string' && body !== '') {
            body = body.replaceAll(/"/g, '\\"')
				template += `
			.requestBody("${body}")`
		}

		template += `
			.ignoreContentType(true)
			.execute();

		System.out.println(response.parse());
	}
}
`
		return template
}


export {toJavaJsoup}