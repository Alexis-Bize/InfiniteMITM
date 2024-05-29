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

import (
	"encoding/json"
	"fmt"
	HaloWaypointLibRequestModule "infinite-mitm/pkg/libs/halowaypoint/modules/request"
	ErrorsModule "infinite-mitm/pkg/modules/errors"
	UtilitiesRequestModule "infinite-mitm/pkg/modules/utilities/request"
	"io"
	"net/http"
	"time"
)

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

// Partial
type MatchStatsResponse struct {
	MatchID   string `json:"MatchId"`
	MatchInfo struct {
		StartTime           time.Time `json:"StartTime"`
		EndTime             time.Time `json:"EndTime"`
		Duration            string    `json:"Duration"`
		LifecycleMode       int       `json:"LifecycleMode"`
		GameVariantCategory int       `json:"GameVariantCategory"`
		LevelID             string    `json:"LevelId"`
		MapVariant          struct {
			AssetKind int    `json:"AssetKind"`
			AssetID   string `json:"AssetId"`
			VersionID string `json:"VersionId"`
		} `json:"MapVariant"`
		UgcGameVariant struct {
			AssetKind int    `json:"AssetKind"`
			AssetID   string `json:"AssetId"`
			VersionID string `json:"VersionId"`
		} `json:"UgcGameVariant"`
		ClearanceID string `json:"ClearanceId"`
		Playlist    struct {
			AssetKind int    `json:"AssetKind"`
			AssetID   string `json:"AssetId"`
			VersionID string `json:"VersionId"`
		} `json:"Playlist"`
		PlaylistExperience  int `json:"PlaylistExperience"`
		PlaylistMapModePair struct {
			AssetKind int    `json:"AssetKind"`
			AssetID   string `json:"AssetId"`
			VersionID string `json:"VersionId"`
		} `json:"PlaylistMapModePair"`
		SeasonID            string `json:"SeasonId"`
		PlayableDuration    string `json:"PlayableDuration"`
		TeamsEnabled        bool   `json:"TeamsEnabled"`
		TeamScoringEnabled  bool   `json:"TeamScoringEnabled"`
		GameplayInteraction int    `json:"GameplayInteraction"`
	} `json:"MatchInfo"`
}

func GetMatchStats(attr HaloWaypointLibRequestModule.RequestAttributes, matchID string) (MatchStatsResponse, error) {
	url := UtilitiesRequestModule.ComputeUrl(HaloWaypointLibRequestModule.GetConfig().Urls.Stats, fmt.Sprintf("/hi/matches/%s/stats", matchID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MatchStatsResponse{}, ErrorsModule.Log(ErrorsModule.ErrHTTPRequestException, err.Error())
	}

	for k, v := range UtilitiesRequestModule.AssignHeaders(map[string]string{
		"User-Agent": attr.UserAgent,
		"X-343-Authorization-Spartan": attr.SpartanToken,
		"Accept": "application/json",
	}) { req.Header.Set(k, v) }

	for k, v := range attr.ExtraHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return MatchStatsResponse{}, ErrorsModule.Log(ErrorsModule.ErrHTTPRequestException, err.Error())
	}
	defer resp.Body.Close()

	err = HaloWaypointLibRequestModule.ValidateResponseStatusCode(resp)
	if err != nil {
		return MatchStatsResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MatchStatsResponse{}, ErrorsModule.Log(ErrorsModule.ErrIOReadException, err.Error())
	}

	var unmarshal MatchStatsResponse
	if err := json.Unmarshal(body, &unmarshal); err != nil {
		return MatchStatsResponse{}, ErrorsModule.Log(ErrorsModule.ErrJSONUnmarshalException, err.Error())
	}

	return unmarshal, nil
}
