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

# definitions
qosservers='[{"region":"SouthAfricaNorth","serverUrl":"pfmsqosprod2-0.southafricanorth.cloudapp.azure.com"},{"region":"WestEurope","serverUrl":"pfmsqosprod2-0.westeurope.cloudapp.azure.com"},{"region":"AustraliaEast","serverUrl":"pfmsqosprod2-0.australiaeast.cloudapp.azure.com"},{"region":"EastAsia","serverUrl":"pfmsqosprod2-0.eastasia.cloudapp.azure.com"},{"region":"SoutheastAsia","serverUrl":"pfmsqosprod2-0.southeastasia.cloudapp.azure.com"},{"region":"BrazilSouth","serverUrl":"pfmsqosprod2-0.brazilsouth.cloudapp.azure.com"},{"region":"EastUs","serverUrl":"pfmsqosprod2-0.eastus.cloudapp.azure.com"},{"region":"EastUs2","serverUrl":"pfmsqosprod2-0.eastus2.cloudapp.azure.com"},{"region":"CentralUs","serverUrl":"pfmsqosprod2-0.centralus.cloudapp.azure.com"},{"region":"NorthCentralUs","serverUrl":"pfmsqosprod2-0.northcentralus.cloudapp.azure.com"},{"region":"SouthCentralUs","serverUrl":"pfmsqosprod2-0.southcentralus.cloudapp.azure.com"},{"region":"WestUs","serverUrl":"pfmsqosprod2-0.westus.cloudapp.azure.com"},{"region":"WestUs2","serverUrl":"pfmsqosprod2-0.westus2.cloudapp.azure.com"},{"region":"NorthEurope","serverUrl":"pfmsqosprod2-0.northeurope.cloudapp.azure.com"},{"region":"JapanEast","serverUrl":"pfmsqosprod2-0.japaneast.cloudapp.azure.com"}]'

regions=($(echo "$qosservers" | grep -o '"region":"[^"]*"' | awk -F'"' '{print $4}'))
server_urls=($(echo "$qosservers" | grep -o '"serverUrl":"[^"]*"' | awk -F'"' '{print $4}'))

# styles
RESET=$(tput sgr0)
STYLE_BOLD=$(tput bold)
COLOR_RED=$(tput setaf 1)
COLOR_GREEN=$(tput setaf 2)
COLOR_YELLOW=$(tput setaf 3)

get_ping_time() {
	local server_url=$1
	ping_time=$(ping -c 1 "$server_url" | grep 'time=' | awk -F'time=' '{print $2}' | awk '{print $1 " ms"}')
	if [[ -z "$ping_time" ]]; then echo "FAIL"
	else echo "$ping_time"
	fi
}

copy_to_clipboard() {
	local text="$1"

	case "$OSTYPE" in
	darwin*)  # macOS
		echo "$text" | pbcopy
		echo "${STYLE_BOLD}✔ The configuration has been copied to your clipboard.${RESET}"
		;;
	msys*|cygwin*|win32*)  # Windows
		echo "$text" | powershell.exe -Command "Set-Clipboard"
		echo "${STYLE_BOLD}✔ The configuration has been copied to your clipboard.${RESET}"
		;;
	*)
		echo "Unsupported OS: $OSTYPE"
		exit 1
		;;
	esac
}

declare -A response_times

echo "${STYLE_BOLD}Servers:${RESET}"
for i in "${!server_urls[@]}"; do
	url=${server_urls[$i]}
	region=${regions[$i]}
	ping=`get_ping_time "$url" | cut -d \. -f 1`

	color_ms="${ping}"
	if [ "$ping" != "FAIL" ]; then
		if [ "$ping" -le 50 ]; then
			color_ms=$(echo ${COLOR_GREEN}${ping})
		elif [ "$ping" -le 100 ] && [ "$ping" -ge 51 ]; then
			color_ms=$(echo ${COLOR_YELLOW}${ping})
		elif [ "$ping" -ge 101 ]; then
			color_ms=$(echo ${COLOR_RED}${ping})
		fi

		color_ms+="ms${RESET}"
	fi

	if [ "$ping" != "FAIL" ]; then
		response_times["$url"]=$ping
	fi

	prefix="├──"
	if [[ "${#server_urls[@]}" -eq i+1 ]]; then
		prefix="└──"
	fi

	echo "${prefix} ${STYLE_BOLD}${region}${RESET}: ${STYLE_BOLD}${color_ms}${RESET}"
done

min_server_url=""
min_ping=999999

max_server_url=""
max_ping=0

for url in "${!response_times[@]}"; do
	ping=${response_times[$url]}

	if [ "$ping" -lt "$min_ping" ]; then
		min_server_url="${url}"
		min_ping=$ping
	fi

	if [ "$ping" -gt "$max_ping" ]; then
		max_server_url="${url}"
		max_ping=$ping
	fi
done

echo ""
echo "> Enter the ${STYLE_BOLD}regions${RESET} you want to ${STYLE_BOLD}prioritize${RESET} (separated by space):"
read -a selected_regions

json="["
for i in "${!regions[@]}"; do
	region="${regions[$i]}"
	server_url="${server_urls[$i]}"
	assigned_server_url="${server_url}"

	if [[ " ${selected_regions[@]} " =~ " ${region} " ]]; then
		assigned_server_url="${min_server_url}"
	else
		assigned_server_url="${max_server_url}"
	fi

	json+="{\"region\":\"$region\",\"serverUrl\": \"$assigned_server_url\"}"
	if [[ "${#regions[@]}" -ne i+1 ]]; then
		json+=","
	fi
done
json+="]"

echo ""
copy_to_clipboard "${json}"
echo ""
