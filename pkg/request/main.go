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

package request

import (
	"bytes"
	"fmt"
	"infinite-mitm/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)
const (
	MITMCacheHeaderKey    = "X-Infinite-MITM-Smart-Cache"
	MITMVersionHeaderKey  = "X-Infinite-MITM-Version"
	MITMProxyHeaderKey    = "X-Infinite-MITM-Proxy"

	CacheControlHeaderKey = "Cache-Control"
	ContentTypeHeaderKey  = "Content-Type"
	ExpiresHeaderKey      = "Expires"
	PragmaHeaderKey       = "Pragma"
	DateHeaderKey         = "Date"
	AgeHeaderKey          = "Age"
)

const (
	MITMCacheHeaderHitValue  = "HIT"
	MITMCacheHeaderMissValue = "MISS"
	MITMProxyEnabledValue    = "enabled"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

func Send(method, url string, payload []byte, header http.Header) ([]byte, http.Header, *errors.MITMError) {
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, http.Header{}, errors.Create(errors.ErrHTTPRequestException, err.Error())
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, http.Header{}, errors.Create(errors.ErrHTTPRequestException, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, http.Header{}, errors.Create(errors.ErrHTTPCodeError, fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.Header{}, errors.Create(errors.ErrHTTPBodyReadException, err.Error())
	}

	return body, resp.Header, nil
}

func StripPort(u string) string {
	parse, err := url.Parse(u)
	if err != nil {
		return u
	}

	return strings.Replace(parse.String(), fmt.Sprintf(":%s", parse.Port()), "", -1)
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
