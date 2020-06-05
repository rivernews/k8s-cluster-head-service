package utilities

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// FetchOption - args for method `Fetch`
type FetchOption struct {
	QueryParams map[string]string
	PostData    map[string]string
	Headers     map[string][]string
	URL         string
	Method      string
}

// Fetch - convenient method to make request with querystring and post data
func Fetch(option FetchOption) string {
	requestURL, _ := url.Parse(option.URL)

	// prepare querystring
	params := url.Values{}
	if option.QueryParams != nil {
		for k, v := range option.QueryParams {
			params.Add(k, v)
		}
	}
	requestURL.RawQuery = params.Encode()

	// prepare post data
	postDataMap := map[string]string{}
	if option.PostData != nil {
		postDataMap = option.PostData
	}
	postDataBuffer := new(bytes.Buffer)
	json.NewEncoder(postDataBuffer).Encode(postDataMap)

	// prepare headers
	headers := map[string][]string{}
	if option.Headers != nil {
		for k, v := range option.Headers {
			headers[k] = v
		}
	}

	// append request config and make request
	req, err := http.NewRequest(option.Method, requestURL.String(), postDataBuffer)
	req.Header = headers
	client := &http.Client{}
	res, err := client.Do(req)

	// log response
	var responseMessage strings.Builder
	responseMessage.WriteString("Response:\n```\n")
	bytesContent, _ := ioutil.ReadAll(res.Body)
	responseMessage.WriteString(string(bytesContent))
	responseMessage.WriteString("\n```\nAny error:\n```\n")
	if err != nil {
		responseMessage.WriteString("🔴 ")
		responseMessage.WriteString(err.Error())
	} else {
		responseMessage.WriteString("🟢 No error")
	}
	responseMessage.WriteString("\n```\n")

	return responseMessage.String()
}
