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

package MITMApplicationMITMServiceHelpers

import (
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

func MatchRequestURL(re *regexp.Regexp) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return MatchURL(req, re)
	}
}

func MatchResponseURL(re *regexp.Regexp) goproxy.RespConditionFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		return MatchURL(resp.Request, re)
	}
}

func MatchURL(req *http.Request, re *regexp.Regexp) bool {
	url := req.URL.Hostname() + req.URL.Path
	query := req.URL.RawQuery
	if query != "" {
		url += "?" + query
	}

	return re.MatchString(url)
}
