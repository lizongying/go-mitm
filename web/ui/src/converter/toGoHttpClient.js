const toGoHttpClient = (url, method, header, body) => {
    let template = `package main

import (
    "bytes"
    "fmt"
    "net/http"
)

func main() {`

    if (typeof url === 'string' && url !== '') {
        url = encodeURI(url)
    } else {
        return ''
    }

    if (typeof method === 'string' && method !== '') {
        method = method.toUpperCase()
        template += `
	req, err := http.NewRequest("${method}", "${url}"`
    } else {
        template += `
	req, err := http.NewRequest("GET", "${url}"`
    }

    if (typeof body === 'string' && body !== '') {
        body = body.replaceAll(/`/g, '\\`')
        template += `, bytes.NewBuffer(\`${body}\`))`
    } else {
        template += `, nil)`
    }
    template+=`
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}`

    if (header != null && typeof header === 'object') {
        template += `
        
	// Set request headers.`
        Object.entries(header).forEach((v) => {
            const value = v[1].replaceAll(/"/g, '\\"')
            template += `
	req.Header.Set("${v[0]}", "${value}")`
        })
    }


    template += `
    
	// Create an HTTP client and send a request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response.
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request failed with status code:", resp.StatusCode)
		return
	}

	fmt.Println("Request was successful")

	// Read the response content.
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response body:", string(responseBody))
}`

    return template
}


export {toGoHttpClient}