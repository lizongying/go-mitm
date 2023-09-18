const toJavaHttpClient = (url, method, header, body) => {
    let template = `package org.example;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

public class HttpClientExample {
    public static void main(String[] args) throws IOException, InterruptedException {
        HttpClient client = HttpClient.newHttpClient();

        HttpRequest request = HttpRequest.newBuilder()`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += `
                .uri(URI.create("${url}"))`
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        switch (method) {
            case 'POST':
                template += `
                .POST`
                break
            case 'GET':
                break
            case 'PUT':
                template += `
                .PUT`
                break
            case 'DELETE':
                template += `
                .DELETE()`
                break
            default:
                return ''
        }
    }

    if (typeof body === 'string' && body !== '') {
        if (!['PUT', 'POST'].includes(method)) {
            return ''
        }
        body = body.replaceAll(/"/g, '\\"')
        template += `(HttpRequest.BodyPublishers.ofString("${body}"))`
    } else {
        if (['PUT', 'POST'].includes(method)) {
            template += `(null)`
        }
    }


    if (header != null && typeof header === 'object') {
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/"/g, '\\"')
            template += `
                .header("${v[0]}", "${value}")`
        })
    }

    template += `
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        System.out.println(response.statusCode());
        System.out.println(response.body());
    }
}
`
    return template
}


export {toJavaHttpClient}