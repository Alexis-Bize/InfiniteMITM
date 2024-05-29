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

import "time"

type FavoriteFilmItem struct {
	Links struct {
	} `json:"Links"`
	Name           string `json:"Name"`
	Description    string `json:"Description"`
	AssetID        string `json:"AssetId"`
	AssetVersionID any    `json:"AssetVersionId"`
	CustomData     struct {
	} `json:"CustomData"`
	VersionRatings any `json:"VersionRatings"`
	AssetKind      int `json:"AssetKind"`
	SearchableData struct {
		PublicName       string  `json:"PublicName"`
		Description      string  `json:"Description"`
		Tags             []any   `json:"Tags"`
		HasNodeGraph     bool    `json:"HasNodeGraph"`
		CloneBehavior    int     `json:"CloneBehavior"`
		PublishedVersion string  `json:"PublishedVersion"`
		IsBanned         bool    `json:"IsBanned"`
		FavoritesCount   int     `json:"FavoritesCount"`
		RatingScore      float64 `json:"RatingScore"`
		RatingCount      int     `json:"RatingCount"`
		AssetState       int     `json:"AssetState"`
		AssetCreatedDate struct {
			ISO8601Date time.Time `json:"ISO8601Date"`
		} `json:"AssetCreatedDate"`
		AssetLastModified struct {
			ISO8601Date time.Time `json:"ISO8601Date"`
		} `json:"AssetLastModified"`
	} `json:"SearchableData"`
}

type FavoriteFilmsResult struct {
	EstimatedTotal int `json:"EstimatedTotal"`
	Start          int `json:"Start"`
	Count          int `json:"Count"`
	ResultCount    int `json:"ResultCount"`
	Results        []FavoriteFilmItem `json:"Results"`
	Links struct {
	} `json:"Links"`
}
