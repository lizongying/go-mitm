const toJavaOkHttp = (url, method, header, body) => {
    let template = `package org.example;

import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

import java.io.IOException;

public class Okhttp3Example {


	public static void main(String[] args) throws IOException {
		OkHttpClient client = new OkHttpClient();
		Request request = new Request.Builder()
`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += `			.url("${url}")`
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
    } else {
        method = 'GET'
    }

    if (header != null && typeof header === 'object') {
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/"/g, '\\"')
            template += `
			.addHeader("${v[0]}", "${value}")`
        })
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/"/g, '\\"')
        switch (method) {
            case 'POST':
                template += `
			.post(RequestBody.create("${body}", null))`
                break
            case 'PUT':
                template += `
			.put(RequestBody.create("${body}", null))`
                break
            case 'PATCH':
                template += `
			.patch(RequestBody.create("${body}", null))`
                break
            case 'DELETE':
                template += `
			.delete(RequestBody.create("${body}", null))`
                break
            default:
                return ''
        }
    }

    template += `								
			.build();

		try (Response response = client.newCall(request).execute()) {
			System.out.println(response.body().string());
		}
	}
}`
    return template
}


export {toJavaOkHttp}