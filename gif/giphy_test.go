package gif

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// test where data response is a object
// test where data response is an empty slice


type GiphyTestSuite struct {
	suite.Suite

	httpClient HttpClientMock
	g giphy
}

func(s *GiphyTestSuite) SetupTest() {
	s.httpClient = HttpClientMock{}
	s.g = giphy{
		httpClient: &s.httpClient,
	}
}

func(s *GiphyTestSuite) TestResultFound() {

	gifUrl := "https://foo.bar/bux"
	responseBody := fmt.Sprintf(`{"data": {"fixed_height_downsampled_url": "%s"}, "meta": {"status": 200}}`, gifUrl)
	responseStatusCode := 200

	s.httpClient.OnDo(responseBody, &responseStatusCode, nil)

	url, found, err := s.g.Random("foo")
	s.Require().NoError(err)
	s.Require().True(found)
	s.Require().EqualValues(gifUrl, url)
}

func(s *GiphyTestSuite) TestNoResultFound() {
	responseBody := `{"data": [], "meta": {"status": 200}}`
	responseStatusCode := 200

	s.httpClient.OnDo(responseBody, &responseStatusCode, nil)

	url, found, err := s.g.Random("foo")
	s.Require().NoError(err)
	s.Require().False(found)
	s.Require().EqualValues("", url)
}

func(s *GiphyTestSuite) TestErrorWhileDoingReq() {
	responseBody := ""

	s.httpClient.OnDo(responseBody, nil, fmt.Errorf("I am a test error"))

	url, found, err := s.g.Random("foo")
	s.Require().Error(err)
	s.Require().False(found)
	s.Require().EqualValues("", url)
}

func TestGiphy(t *testing.T) {
	suite.Run(t, &GiphyTestSuite{})
}