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
	"fmt"
	"infinite-mitm/pkg/domains"
	"infinite-mitm/pkg/request"
	"infinite-mitm/pkg/resources"
	"infinite-mitm/pkg/utilities"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type StrategyType string

type SmartCache struct {
	strategy StrategyType
	duration time.Duration
	items    map[string]*SmartCacheItem
}

type SmartCacheYAMLOptions struct {
	Enabled  bool
	Strategy StrategyType
	TTL      string
}

type SmartCacheItem struct {
	Body       []byte
	Header     http.Header
	Created    time.Time
	Expires    time.Time
}

const (
	Memory     StrategyType = "memory"
	Persistent StrategyType = "persistent"
	defaultDuration = 7 * 24 * time.Hour
)

var RWMutex = &sync.RWMutex{}
var fileOSMutex = &sync.Mutex{}
var flushSmartCacheMutex = &sync.Mutex{}

func init() {
	gob.Register(http.Header{})
	gob.Register(time.Time{})
}

func New(strategy StrategyType, ttl string) *SmartCache {
	if strategy != Memory && strategy != Persistent {
		strategy = Memory
	}

	sc := &SmartCache{
		strategy: strategy,
		duration: parseDuration(ttl),
		items:    make(map[string]*SmartCacheItem),
	}

	return sc
}

func Flush() {
	flushSmartCacheMutex.Lock()
	defer flushSmartCacheMutex.Unlock()

	os.RemoveAll(resources.GetSmartCacheDirPath())
}

func (s *SmartCache) Get(key string) *SmartCacheItem {
	if s.strategy == Persistent {
		fileOSMutex.Lock()
		defer fileOSMutex.Unlock()

		target := filepath.Join(resources.GetSmartCacheDirPath(), key)
		file, err := os.Open(target)
		if err != nil {
			return nil
		}
		defer file.Close()

		var item *SmartCacheItem
		err = gob.NewDecoder(file).Decode(&item);
		if err != nil {
			return nil
		}

		if s.isExpired(item) {
			os.Remove(target)
			return nil
		}

		return item
	}

	RWMutex.RLock()
	defer RWMutex.RUnlock()

	if item, exists := s.items[key]; exists && !s.isExpired(item) {
		item.Header.Set(request.DateHeaderKey, time.Now().Format(time.RFC1123))
		if item.Created != (time.Time{}) {
			since := time.Since(item.Created)
			seconds := int(since.Seconds())
			item.Header.Set(request.AgeHeaderKey, fmt.Sprintf("%d", seconds))
		}

		return item
	}

	return nil
}

func (s *SmartCache) Write(key string, item *SmartCacheItem) {
	item.Created = time.Now()
	item.Expires = time.Now().Add(s.duration)

	item.Header.Set(request.ExpiresHeaderKey, item.Expires.Format(time.RFC1123))
	item.Header.Set(request.DateHeaderKey, item.Created.Format(time.RFC1123))
	item.Header.Set(request.AgeHeaderKey, "0")
	item.Header.Del(request.CacheControlHeaderKey)

	if s.strategy == Memory {
		RWMutex.Lock()
		s.items[key] = item
		RWMutex.Unlock()
		return
	}

	go func() {
		fileOSMutex.Lock()
		defer fileOSMutex.Unlock()

		target := filepath.Join(resources.GetSmartCacheDirPath(), key)
		file, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			defer file.Close()
			gob.NewEncoder(file).Encode(item)
		}
	}()
}

func (s *SmartCache) isExpired(item *SmartCacheItem) bool {
	if item.Expires == (time.Time{}) {
		return true
	}

	return time.Now().After(item.Expires)
}

func (s *SmartCache) CreateKey(input string, extra...string) string {
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

	return CreateHash(strings.Join(append([]string{input}, extra...), ":"))
}

func parseDuration(ttl string) time.Duration {
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return defaultDuration
	}

	return duration
}

func CreateHash(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
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
