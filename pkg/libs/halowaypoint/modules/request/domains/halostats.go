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

package HaloWaypointLibRequestModuleDomains

type MatchSpectateResponse struct {
	FilmStatusBond int `json:"FilmStatusBond"`
	CustomData     struct {
		FilmLength int `json:"FilmLength"`
		Chunks     []struct {
			Index                            int `json:"Index"`
			ChunkStartTimeOffsetMilliseconds int `json:"ChunkStartTimeOffsetMilliseconds"`
			DurationMilliseconds             int `json:"DurationMilliseconds"`
			ChunkSize                        int `json:"ChunkSize"`
			FileRelativePath                 string `json:"FileRelativePath"`
			ChunkType                        int `json:"ChunkType"`
		} `json:"Chunks"`
		HasGameEnded           bool `json:"HasGameEnded"`
		ManifestRefreshSeconds int `json:"ManifestRefreshSeconds"`
		MatchID                string `json:"MatchId"`
		FilmMajorVersion       int `json:"FilmMajorVersion"`
	} `json:"CustomData"`
	BlobStoragePathPrefix      string `json:"BlobStoragePathPrefix"`
	AssetID                    string `json:"AssetId"`
}
