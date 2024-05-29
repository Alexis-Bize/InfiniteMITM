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

package Utilities

func GetValueAndIndex(m map[string]string, keys []string, key string) (string, int) {
	value, exists := m[key]
	if !exists {
		return "", -1
	}

	for index, k := range keys {
		if k == key {
			return value, index
		}
	}

	return value, -1
}
