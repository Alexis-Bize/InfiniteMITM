#!/usr/bin/env bash

# Copyright 2024 Alexis Bize
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

script_path=$(dirname "$0")
script_path=$(cd "$script_path" && pwd)
root_path="${script_path}/.."
build_dir="${root_path}/build"

config_file="${root_path}/configs/application.yaml"
name=$(grep "^name:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
port=$(grep "^port:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')

enableProxy() {
	echo "Enabling ${name} proxy..."
	networksetup -setwebproxy Wi-Fi localhost "${port}"
	networksetup -setsecurewebproxy Wi-Fi localhost "${port}"
	networksetup -setwebproxystate Wi-Fi on
	networksetup -setsecurewebproxystate Wi-Fi on
}

disableProxy() {
	echo "Disabling ${name} proxy..."
	networksetup -setwebproxystate Wi-Fi off
	networksetup -setsecurewebproxystate Wi-Fi off
}

if [ "$1" == "on" ]; then
  enableProxy
elif [ "$1" == "off" ]; then
  disableProxy
else
  echo "Usage: $0 {on|off}"
  exit 1
fi
