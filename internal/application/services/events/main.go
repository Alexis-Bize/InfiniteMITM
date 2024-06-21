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

package MITMApplicationEventsService

import "encoding/json"

const (
	RestartServer = "server.restart"
	ProxyRequestSent = "request.sent"
	ProxyResponseReceived = "response.received"
	ProxyStatusMessage = "proxy.status_message"
)

type ProxyRequestEventData struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	Headers     map[string]string `json:"headers"`
	Body        []byte `json:"body"`
	Proxified   bool `json:"proxified"`
	SmartCached bool `json:"smart_cached"`
}

type ProxyResponseEventData struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	Status      int `json:"status"`
	Headers     map[string]string `json:"headers"`
	Body        []byte `json:"body"`
	Proxified   bool `json:"proxified"`
	SmartCached bool `json:"smart_cached"`
}

func StringifyRequestEventData(data ProxyRequestEventData) string {
	marshal, _ := json.Marshal(data)
	return string(marshal)
}

func StringifyResponseEventData(data ProxyResponseEventData) string {
	marshal, _ := json.Marshal(data)
	return string(marshal)
}

func ParseRequestEventData(data string) ProxyRequestEventData {
	var unmarshal ProxyRequestEventData
	json.Unmarshal([]byte(data), &unmarshal)
	return unmarshal
}

func ParseResponseEventData(data string) ProxyResponseEventData {
	var unmarshal ProxyResponseEventData
	json.Unmarshal([]byte(data), &unmarshal)
	return unmarshal
}
