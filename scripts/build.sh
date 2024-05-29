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
assets_path="${root_path}/assets"

build_dir="${root_path}/build"
build_bin_dir="${build_dir}/bin"
build_archive_dir="${build_dir}/archive"

config_file="${root_path}/configs/application.yaml"

name=$(grep "^name:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
description=$(grep "^description:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
version=$(grep "^version:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
author=$(grep "^author:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
repository=$(grep "^repository:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')

name_lower="${name,,}"
package_name="${name}-${version}"

target_platforms=(
	"windows/amd64"
	"darwin/amd64"
)

IFS='.' read -r major minor patch <<< "$version"
build=0

for platform in "${target_platforms[@]}"
do
	platform_split=(${platform//\// })

	GOOS="${platform_split[0]}"
	GOARCH="${platform_split[1]}"

	base_output_name="${package_name}-${GOOS}-${GOARCH}"
	output_name="${base_output_name}"

	#region windows

	if [ $GOOS = "windows" ]; then
		output_name+='.exe'

		versioninfo='{
			"FixedFileInfo": {
				"FileVersion": {
					"Major": '"$major"',
					"Minor": '"$minor"',
					"Patch": '"$patch"',
					"Build": '"$build"'
				},
				"ProductVersion": {
					"Major": '"$major"',
					"Minor": '"$minor"',
					"Patch": '"$patch"',
					"Build": '"$build"'
				},
				"FileFlagsMask": "3f",
				"FileFlags ": "00",
				"FileOS": "040004",
				"FileType": "01",
				"FileSubType": "00"
			},
			"StringFileInfo": {
				"Comments": "",
				"CompanyName": "'"$author"'",
				"FileDescription": "'"$description"'",
				"FileVersion": "v'"$version.$build"'",
				"InternalName": "'"$name".exe'",
				"LegalCopyright": "Copyright (c) '$(date +"%Y")' '"$author"'",
				"LegalTrademarks": "",
				"OriginalFilename": "main.go",
				"PrivateBuild": "",
				"ProductName": "'"$name"'",
				"ProductVersion": "v'"$version.$build"'",
				"SpecialBuild": ""
			},
			"VarFileInfo": {
				"Translation": {
					"LangID": "0409",
					"CharsetID": "04B0"
				}
			},
			"IconPath": "",
			"ManifestPath": ""
		}'

		echo $versioninfo >| ${root_path}/versioninfo.json

		go generate
	fi

	#endregion

	env GOOS=$GOOS GOARCH=$GOARCH go build -o "${build_bin_dir}/${output_name}"
	(cd "${build_bin_dir}" && zip -rj "${build_archive_dir}/${base_output_name}.zip" "${output_name}")
done
