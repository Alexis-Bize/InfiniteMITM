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

package MITMApplicationPromptService

type PromptOption int

const (
	StartProxyServer PromptOption = iota
	InstallRootCertificate
	ForceKillProxy
	Exit
)

var optionToString = map[PromptOption]string{
	StartProxyServer:       "Start Proxy Server",
	InstallRootCertificate: "Install Root Certificate",
	ForceKillProxy:         "Force Kill Proxy",
	Exit:                   "Exit",
}

func (d PromptOption) String() string {
	return optionToString[d]
}

func (d PromptOption) Is(option string) bool {
	return d.String() == option
}
