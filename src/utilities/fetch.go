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
	QueryParams         map[string]string
	PostData            interface{}
	Headers             map[string][]string
	URL                 string
	Method              string
	DisableHumanMessage bool
}

// Fetch - convenient method to make request with querystring and post data
func Fetch(option FetchOption) ([]byte, string, error) {
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
	postDataBuffer := new(bytes.Buffer)
	if option.PostData != nil {
		postDataMap := option.PostData
		json.NewEncoder(postDataBuffer).Encode(postDataMap)
	} else {
		postDataMap := map[string]string{}
		json.NewEncoder(postDataBuffer).Encode(postDataMap)
	}

	// prepare headers
	headers := map[string][]string{}
	if option.Headers != nil {
		for k, v := range option.Headers {
			headers[k] = v
		}
	}

	// append request config and make request
	req, fetchErr := http.NewRequest(option.Method, requestURL.String(), postDataBuffer)
	req.Header = headers
	client := &http.Client{}
	res, fetchErr := client.Do(req)

	bytesContent, _ := ioutil.ReadAll(res.Body)

	// log response
	var responseMessage strings.Builder
	if !option.DisableHumanMessage {
		responseMessage.WriteString("Response:\n```\n")
		responseMessage.WriteString(string(bytesContent))
		responseMessage.WriteString("\n```\nAny error:\n```\n")
		if fetchErr != nil {
			responseMessage.WriteString("🔴 ")
			responseMessage.WriteString(fetchErr.Error())
		} else {
			responseMessage.WriteString("🟢 No error")
		}
		responseMessage.WriteString("\n```\n")
	}

	return bytesContent, responseMessage.String(), fetchErr
}
