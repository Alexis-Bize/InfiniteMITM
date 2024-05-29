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

type FilmResponse struct {
	CustomData struct {
		FilmLength int `json:"FilmLength"`
		Chunks     []struct {
			Index                            int    `json:"Index"`
			ChunkStartTimeOffsetMilliseconds int    `json:"ChunkStartTimeOffsetMilliseconds"`
			DurationMilliseconds             int    `json:"DurationMilliseconds"`
			ChunkSize                        int    `json:"ChunkSize"`
			FileRelativePath                 string `json:"FileRelativePath"`
			ChunkType                        int    `json:"ChunkType"`
		} `json:"Chunks"`
		HasGameEnded           bool   `json:"HasGameEnded"`
		ManifestRefreshSeconds int    `json:"ManifestRefreshSeconds"`
		MatchID                string `json:"MatchId"`
		FilmMajorVersion       int    `json:"FilmMajorVersion"`
	} `json:"CustomData"`
	Tags               []any  `json:"Tags"`
	MapLink            any    `json:"MapLink"`
	UgcGameVariantLink any    `json:"UgcGameVariantLink"`
	AssetID            string `json:"AssetId"`
	VersionID          string `json:"VersionId"`
	PublicName         string `json:"PublicName"`
	Description        string `json:"Description"`
	Files              struct {
		Prefix            string   `json:"Prefix"`
		FileRelativePaths []string `json:"FileRelativePaths"`
		PrefixEndpoint    struct {
			AuthorityID                              string `json:"AuthorityId"`
			Path                                     string `json:"Path"`
			QueryString                              any    `json:"QueryString"`
			RetryPolicyID                            string `json:"RetryPolicyId"`
			TopicName                                string `json:"TopicName"`
			AcknowledgementTypeID                    int    `json:"AcknowledgementTypeId"`
			AuthenticationLifetimeExtensionSupported bool   `json:"AuthenticationLifetimeExtensionSupported"`
			ClearanceAware                           bool   `json:"ClearanceAware"`
		} `json:"PrefixEndpoint"`
	} `json:"Files"`
	Contributors []any `json:"Contributors"`
	AssetHome    int   `json:"AssetHome"`
	AssetStats   struct {
		PlaysRecent      int     `json:"PlaysRecent"`
		PlaysAllTime     int     `json:"PlaysAllTime"`
		Favorites        int     `json:"Favorites"`
		Likes            int     `json:"Likes"`
		Bookmarks        int     `json:"Bookmarks"`
		ParentAssetCount int     `json:"ParentAssetCount"`
		AverageRating    float64 `json:"AverageRating"`
		NumberOfRatings  int     `json:"NumberOfRatings"`
	} `json:"AssetStats"`
	InspectionResult int `json:"InspectionResult"`
	CloneBehavior    int `json:"CloneBehavior"`
	Order            int `json:"Order"`
	PublishedDate    struct {
		ISO8601Date time.Time `json:"ISO8601Date"`
	} `json:"PublishedDate"`
	VersionNumber int    `json:"VersionNumber"`
	Admin         string `json:"Admin"`
}

func GetFilmByAssetID(attr HaloWaypointLibRequestModule.RequestAttributes, assetID string) (FilmResponse, error) {
	url := UtilitiesRequestModule.ComputeUrl(HaloWaypointLibRequestModule.GetConfig().Urls.Discovery, fmt.Sprintf("/hi/films/%s", assetID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FilmResponse{}, ErrorsModule.Log(ErrorsModule.ErrHTTPRequestException, err.Error())
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
		return FilmResponse{}, ErrorsModule.Log(ErrorsModule.ErrHTTPRequestException, err.Error())
	}
	defer resp.Body.Close()

	err = HaloWaypointLibRequestModule.ValidateResponseStatusCode(resp)
	if err != nil {
		return FilmResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FilmResponse{}, ErrorsModule.Log(ErrorsModule.ErrIOReadException, err.Error())
	}

	var unmarshal FilmResponse
	if err := json.Unmarshal(body, &unmarshal); err != nil {
		return FilmResponse{}, ErrorsModule.Log(ErrorsModule.ErrJSONUnmarshalException, err.Error())
	}

	return unmarshal, nil
}

