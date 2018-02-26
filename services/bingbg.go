package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"git.m/svcman/common"
	"git.m/svcman/data"
)

const (
	bingImageTTL      = 86400
	bingDailyImageURL = "http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
)

//BingBGFetcher ...
type BingBGFetcher struct {
	httpClient *http.Client
	cache      *data.CacheService
}

//NewBingBGFetcher ...
func NewBingBGFetcher(db *data.DataStore) *BingBGFetcher {
	var http = &http.Client{
		Timeout: time.Second * 2,
	}
	return &BingBGFetcher{
		httpClient: http,
		cache:      db.Cache,
	}
}

//GetBGImage ...
func (api *BingBGFetcher) GetBGImage(resp http.ResponseWriter, r *http.Request) {
	if imageInfo, err := api.getImageInfo(); err == nil {
		req, _ := http.NewRequest("GET", "http://bing.com/"+imageInfo, nil)
		if response, err := api.httpClient.Do(req); err == nil {
			if image, err := ioutil.ReadAll(response.Body); err == nil {
				common.Logger.Debugln("Got")
				resp.Write(image)
			} else {
				common.CreateFailureResponse(err, "getbgimage", 500)
			}
		}
	} else {
		resp.WriteHeader(404)
		common.CreateFailureResponse(err, "getbgimage", 500)
	}
}
func (api *BingBGFetcher) getImageInfo() (string, error) {
	var imageInfo data.BingDailyImage
	var imageURL string

	if imageURL = api.cache.GetString("ui", "bgimage"); imageURL == "" {
		req, _ := http.NewRequest("GET", bingDailyImageURL, nil)
		if response, err := api.httpClient.Do(req); err == nil {
			if response.StatusCode == 200 {
				body, _ := ioutil.ReadAll(response.Body)
				if err = json.Unmarshal(body, &imageInfo); err == nil {
					api.cache.PutStringWithExpiration("ui", "bgimage", imageInfo.Images[0].URL, bingImageTTL)
					return imageInfo.Images[0].URL, nil
				} else {
					common.CreateFailureResponse(err, "getImageInfo", 500)
					return "", err
				}
			} else {
				common.CreateFailureResponse(err, "getImageInfo", 500)
				return "", errors.New(response.Status)
			}
		} else {
			common.CreateFailureResponse(err, "getImageInfo", 500)
			return "", err
		}
	} else {
		return imageURL, nil
	}
}
