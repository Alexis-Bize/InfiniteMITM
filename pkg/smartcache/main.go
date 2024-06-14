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

package smartcache

import (
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/resources"
	"infinite-mitm/pkg/utilities"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type StrategyType string

type SmartCache struct {
	strategy StrategyType
	items    map[string]*SmartCacheItem
}

type SmartCacheItem struct {
	Body       []byte
	Header     http.Header
	persisted  bool
	expires    time.Time
}

const (
	Memory     StrategyType = "memory"
	Persistent StrategyType = "persistent"
	expirationDelta = 48 * time.Hour
)

var readMutex, writeMutex, flushMutex sync.Mutex

func init() {
	gob.Register(http.Header{})
	gob.Register(time.Time{})
}

func New(strategy StrategyType) *SmartCache {
	if strategy != Memory && strategy != Persistent {
		strategy = Memory
	}

	sc := &SmartCache{strategy: strategy, items: make(map[string]*SmartCacheItem)}
	return sc
}

func Flush() {
	flushMutex.Lock()
	defer flushMutex.Unlock()

	os.RemoveAll(resources.GetSmartCacheDirPath())
	os.MkdirAll(resources.GetSmartCacheDirPath(), 0755)
}

func CreateHash(input string) string {
	input = strings.ToLower(input)

	hash := sha1.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func IsURLSmartCachable(target string, method string) bool {
	if method != http.MethodGet {
		return false
	}

	parse, _ := url.Parse(target)
	hostname := parse.Hostname()
	isSupportedDomain := utilities.Contains(domains.SmartCachableHostnames, hostname)

	if isSupportedDomain {
		path := strings.ToLower(parse.Path)
		if hostname == domains.DomainToHostname(domains.Skill) {
			return strings.HasSuffix(path, "/skill")
		} else if hostname == domains.DomainToHostname(domains.HaloStats) {
			return strings.HasSuffix(path, "/stats")
		} else if hostname == domains.DomainToHostname(domains.Authoring) {
			return !strings.Contains(path, "/favorites") && !strings.HasSuffix(path, "/myfiles") && !strings.Contains(path, "/ratings") && !strings.HasSuffix(path, "/latest") && !strings.HasSuffix(path, "/assets")
		} else if hostname == domains.DomainToHostname(domains.Discovery) {
			return strings.Contains(path, "/versions/") && !strings.HasSuffix(path, "/latest")
		}

		return true
	}

	return false
}

func (s *SmartCache) Read(key string) *SmartCacheItem {
	readMutex.Lock()
	defer readMutex.Unlock()

	if item, exists := s.items[key]; exists && !s.isExpired(item) {
		if item.Header.Get(request.DateHeaderKey) != "" {
			item.Header.Set(request.DateHeaderKey, time.Now().Format(time.RFC1123))
		}

		return item
	}

	if s.strategy == Persistent {
		file, err := os.Open(path.Join(resources.GetSmartCacheDirPath(), key))
		if err == nil {
			defer file.Close()

			var item SmartCacheItem
			if err = gob.NewDecoder(file).Decode(&item); err == nil {
				item.persisted = true
				s.Write(key, &item)
				return &item
			}
		}
	}

	return nil
}

func (s *SmartCache) Write(key string, item *SmartCacheItem) {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	item.expires = time.Now().Add(expirationDelta)
	if item.Header.Get(request.ExpiresHeaderKey) != "" {
		item.Header.Set(request.ExpiresHeaderKey, item.expires.Format(time.RFC1123))
	}

	if item.Header.Get(request.DateHeaderKey) != "" {
		item.Header.Set(request.DateHeaderKey, time.Now().Format(time.RFC1123))
	}

	if item.Header.Get(request.AgeHeaderKey) != "" {
		item.Header.Set(request.AgeHeaderKey, "0")
	}

	s.items[key] = item

	if !item.persisted && s.strategy == Persistent {
		file, err := os.Create(path.Join(resources.GetSmartCacheDirPath(), key)); if err != nil {
			return
		}

		defer file.Close()
		item.persisted = true
		gob.NewEncoder(file).Encode(item)
	}
}

func (s *SmartCache) isExpired(item *SmartCacheItem) bool {
	if item.expires == (time.Time{}) {
		return true
	}

	return time.Now().After(item.expires)
}

func (s *SmartCache) CreateKey(input string) string {
	parse, err := url.Parse(input)
	if err == nil {
		hostname := parse.Hostname()

		if hostname == domains.DomainToHostname(domains.GameCMS) {
			queryParams := parse.Query()
			queryParams.Del("flight")
			parse.RawQuery = queryParams.Encode()
		}

		normalizedPath := strings.ReplaceAll(parse.Path, "//", "/")
		parse.Path = normalizedPath
		input = request.StripPort(input)
	}

	return CreateHash(input)
}