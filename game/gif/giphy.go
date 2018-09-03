package gif

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	giphyBaseURI   = "http://api.giphy.com/v1/gifs/%s?api_key=dc6zaTOxFJmzC&%s"
	giphyRandomApi = "random"
)

func NewGiphy() Gif {
	return &giphy{}
}

type giphy struct {
}

type giphyResponse struct {
	Data *giphyData `json:"data"`
	Meta giphyMeta  `json:"meta"`
}

type giphyData struct {
	FixedHeightDownsampledUrl string `json:"fixed_height_downsampled_url"`
}

type giphyMeta struct {
	Status int64 `json:"status"`
}

func (g *giphy) Random(searchQuery string) (string, error) {
	return g.doRequest(giphyRandomApi, fmt.Sprintf("&tag=%s", searchQuery))
}

func (g *giphy) doRequest(api string, searchQuery string) (string, error) {

	callURI := fmt.Sprintf(giphyBaseURI, api, searchQuery)

	res, err := http.Get(callURI)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data giphyResponse

	err = json.Unmarshal(rawBody, &data)
	if err != nil {
		return "", err
	}

	if data.Meta.Status == 200 && data.Data != nil {
		return data.Data.FixedHeightDownsampledUrl, nil
	}

	return "", fmt.Errorf("unexpected response from giphy api. %s", string(rawBody))
}
