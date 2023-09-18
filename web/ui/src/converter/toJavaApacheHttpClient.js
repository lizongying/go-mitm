const toJavaApacheHttpClient = (url, method, header, body) => {
		let template = `package org.example;

import org.apache.hc.client5.http.classic.methods.HttpDelete;
import org.apache.hc.client5.http.classic.methods.HttpGet;
import org.apache.hc.client5.http.classic.methods.HttpHead;
import org.apache.hc.client5.http.classic.methods.HttpOptions;
import org.apache.hc.client5.http.classic.methods.HttpPatch;
import org.apache.hc.client5.http.classic.methods.HttpPost;
import org.apache.hc.client5.http.classic.methods.HttpPut;
import org.apache.hc.client5.http.classic.methods.HttpTrace;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.HttpEntity;
import org.apache.hc.core5.http.HttpStatus;
import org.apache.hc.core5.http.io.HttpClientResponseHandler;
import org.apache.hc.core5.http.io.entity.EntityUtils;
import org.apache.hc.core5.http.io.entity.StringEntity;

import java.io.IOException;

public class ApacheHttpClient {
	public static void main(String[] args) throws IOException {
		CloseableHttpClient httpClient = HttpClients.createDefault();
`

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
		HttpGet request = new HttpGet("${url}");`
					break
				case 'POST':
					template += `
		HttpPost request = new HttpPost("${url}");`
					break
				case 'HEAD':
					template += `
		HttpHead request = new HttpHead("${url}");`
					break
				case 'PUT':
					template += `
		HttpPut request = new HttpPut("${url}");`
					break
				case 'DELETE':
					template += `
		HttpDelete request = new HttpDelete("${url}");`
					break
				case 'PATCH':
					template += `
		HttpPatch request = new HttpPatch("${url}");`
					break
				case 'OPTIONS':
					template += `
		HttpOptions request = new HttpOptions("${url}");`
					break
				default:
					return ''
			}
		} else {
            template += `
		HttpGet request = new HttpGet("${url}");`
        }

		if (header != null && typeof header === 'object') {
				Object.entries(header).forEach((v) => {
					const value = v[1].replaceAll(/"/g, '\\"')
						template += `
		request.setHeader("${v[0]}", "${value}");`
				})
		}

		if (typeof body === 'string' && body !== '') {
            body = body.replaceAll(/"/g, '\\"')
				template += `
		request.setEntity(new StringEntity("${body}"));`
		}

		template += `
		httpClient.execute(request, (HttpClientResponseHandler<CloseableHttpResponse>) response -> {
			if (response.getCode() == HttpStatus.SC_OK) {
				HttpEntity entity = response.getEntity();
				System.out.println(EntityUtils.toString(entity));
			} else {
				System.err.println("statusCode:" + response.getCode());
			}
			return null;
		}
	}
}`
		return template
}


export {toJavaApacheHttpClient}