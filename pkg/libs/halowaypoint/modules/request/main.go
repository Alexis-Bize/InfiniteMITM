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

package HaloWaypointLibRequestModule

import (
	ErrorsModule "infinite-mitm/pkg/modules/errors"
	"net/http"
)

type RequestAttributes struct {
	SpartanToken string
	UserAgent    string
	ExtraHeaders map[string]string
}

func ValidateResponseStatusCode(resp *http.Response) error {
	var err error

	if resp.StatusCode == 401 {
		err = ErrorsModule.Log(ErrorsModule.ErrHTTPUnauthorized, "current STv4 is invalid or has expired")
	} else if resp.StatusCode == 403 {
		err = ErrorsModule.Log(ErrorsModule.ErrHTTPForbidden, "your are not allowed to perform this action")
	} else if resp.StatusCode == 404 {
		err = ErrorsModule.Log(ErrorsModule.ErrHTTPNotFound, "desired entity does not exist")
	} else if resp.StatusCode < 200 || (resp.StatusCode >= 300 && resp.StatusCode < 500) {
		err = ErrorsModule.Log(ErrorsModule.ErrHTTPRequestException, "something went wrong")
	} else if resp.StatusCode >= 500 {
		err = ErrorsModule.Log(ErrorsModule.ErrHTTPInternal, "something went wrong")
	}

	return err
}
