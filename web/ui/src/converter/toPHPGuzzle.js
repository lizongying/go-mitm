const toPHPGuzzle = (url, method, header, body) => {
    let template = `<?php
require 'vendor/autoload.php';

use GuzzleHttp\\Client;

$client = new Client();`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
    } else {
        return ''
    }

    template += `
try {`
    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `
	$response = $client->request('${method}', '${url}'`
    } else {
        template += `
	$response = $client->request('GET', '${url}'`
    }

    let config = false
    if (header != null && typeof header === 'object') {
        if (config) {
            template += `,`
        } else {
            template += `, [`
        }
        config = true

        template += `
		'headers' => [`
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            template += `
			'${v[0]}' => '${value}',`
        })
        template += `
		]`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/'/g, '\\\'')
        if (config) {
            template += `,`
        } else {
            template += `, [`
        }
        template += `
		'body' => '${body}',`
        config = true
    }

    if (config) {
        template += `
	]`
    }
    template += `);
    
	echo $response->getBody();
	echo $response->getStatusCode();
	var_dump($response->getHeaders());
} catch (Exception $e) {
  // Handle the exception
  echo 'Caught exception: ' . $e->getMessage();
}`
    return template
}


export {toPHPGuzzle}