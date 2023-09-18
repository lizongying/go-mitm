const toPHPcURL = (url, method, header, body) => {
    let template = `<?php
$ch = curl_init();`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
        template += `
curl_setopt($ch, CURLOPT_URL, '${url}');`
    } else {
        return ''
    }

    template += `
    
// If you want to retrieve the request headers.
// curl_setopt($ch, CURLOPT_HEADER, true);

curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);`

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, '${method}');`
    }

    if (header != null && typeof header === 'object') {
        template += `
curl_setopt($ch, CURLOPT_HTTPHEADER, [`
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/'/g, '\\\'')
            template += `
	'${v[0]}: ${value}',`
        })
        template += `
]);`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/'/g, '\\\'')
        template += `
curl_setopt($ch, CURLOPT_POSTFIELDS, '${body}');`
    }

    template += `
$response = curl_exec($ch);

//if ($error = curl_error($ch)) {
//  die($error);
//}

echo $response;

// Get the HTTP status code.
//$code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
//echo $code;

// Get request headers.
//$headerSize = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
//$header = substr($response, 0, $headerSize);
//echo $header;

curl_close($ch);`

    return template
}


export {toPHPcURL}