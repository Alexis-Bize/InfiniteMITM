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
	errors "infinite-mitm/pkg/modules/errors"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

func ComputeUrl(baseUrl string, path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return baseUrl + path
}

func ReplayRequestWithJSONAccept(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()

	clonedReq := req.Clone(req.Context())
	clonedReq.Body = io.NopCloser(bytes.NewReader(body))
	clonedReq.Header.Set("Accept", "application/json")
	clonedReq.RequestURI = ""

	_, err = httputil.DumpRequest(clonedReq, true)
	if err != nil {
		return nil, errors.Create(errors.ErrHTTPInternal, err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(clonedReq)
	if err != nil {
		return nil, errors.Create(errors.ErrHTTPInternal, err.Error())
	}

	return resp, nil
}

func AssignHeaders(extraHeaders map[string]string) map[string]string {
	headers := map[string]string{}
	for k, v := range extraHeaders {
		headers[k] = v
	}

	return headers
}

func ExtractHeaderValue(req *http.Request, key string) string {
	spartanToken := req.Header.Get(key)
	return spartanToken
}
