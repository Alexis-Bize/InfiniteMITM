// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package UtilitiesRequestModule

import (
	"bytes"
	"fmt"
	errors "infinite-mitm/pkg/modules/errors"
	"io"
	"net/http"
	"strings"
)

func Send(method, url string, payload []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Create(errors.ErrHTTPRequestException, err.Error())
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Create(errors.ErrHTTPRequestException, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Create(errors.ErrHTTPError, fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Create(errors.ErrHTTPBodyReadException, err.Error())
	}

	return body, nil
}

func ComputeUrl(baseUrl string, path string) string {
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return baseUrl + path
}

func HeadersToMap(header http.Header) map[string]string {
	headersMap := make(map[string]string)
	for key, values := range header {
		headersMap[key] = strings.Join(values, ", ")
	}

	return headersMap
}

func ExtractHeaderValue(req *http.Request, key string) string {
	value := req.Header.Get(key)
	return value
}
