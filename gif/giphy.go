package gif

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/bpoetzschke/bin.go/logger"
)

const (
	giphyBaseURI   = "http://api.giphy.com/v1/gifs/%s?api_key=%s%s"
	giphyRandomApi = "random"
)

// HTTPDoer is a 3rd party interface for the HTTP do method
type HTTPDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewGiphy() (Gif, error) {

	apiKey := os.Getenv("GIPHY_API_KEY")
	if apiKey == "" {
		err := fmt.Errorf("giphy api key is not set.")
		logger.Error("%s", err)
		return nil, err
	}

	return &giphy{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
	}, nil
}

type giphy struct {
	httpClient HTTPDoer
	apiKey     string
}

type giphyResponse struct {
	RawData interface{} `json:"data"`
	Meta    giphyMeta   `json:"meta"`
}

type giphyData struct {
	FixedHeightDownsampledUrl string `json:"fixed_height_downsampled_url"`
}

type giphyMeta struct {
	Status int64 `json:"status"`
}

func (g *giphy) Random(searchQuery string) (string, bool, error) {
	return g.doRequest(giphyRandomApi, fmt.Sprintf("&tag=%s", url.QueryEscape(searchQuery)))
}

func (g *giphy) doRequest(api string, searchQuery string) (string, bool, error) {

	callURI := fmt.Sprintf(giphyBaseURI, api, g.apiKey, searchQuery)

	req, err := http.NewRequest("GET", callURI, nil)
	if err != nil {
		return "", false, fmt.Errorf("error while creating http request: Error: %s", err)
	}

	res, err := g.httpClient.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("error while sending http request. Error: %s", err)
	}

	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)

	var data giphyResponse

	err = json.Unmarshal(rawBody, &data)
	if err != nil {
		return "", false, err
	}

	if _, ok := data.RawData.(map[string]interface{}); ok {
		rawGiphyData, err := json.Marshal(data.RawData)
		if err != nil {
			return "", false, fmt.Errorf("failed to marshal giphy response into json. Error: %s", err)
		}

		var d giphyData
		err = json.Unmarshal(rawGiphyData, &d)
		if err != nil {
			return "", false, fmt.Errorf("failed to unmarshal giphy response into value of type %T. Error: %s", d, err)
		}

		return d.FixedHeightDownsampledUrl, true, nil
	}

	return "", false, nil
}
