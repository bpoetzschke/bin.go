package gif

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/stretchr/testify/mock"
)

// HTTP client mock

type HttpClientMock struct {
	mock.Mock
}

func (m *HttpClientMock) AnythingOfType(obj interface{}) mock.AnythingOfTypeArgument {
	return mock.AnythingOfType(fmt.Sprintf("%T", obj))
}

func (m *HttpClientMock) OnDo(body string, responseCode *int, err error) *mock.Call {
	var response *http.Response

	if err == nil {
		response = &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(body)),
		}

		header := http.Header{}
		header.Set("Response-Header", "Response-Header-Value")

		response.Header = header

		if responseCode != nil {
			response.StatusCode = *responseCode
		}
	}

	return m.On(
		"Do",
		m.AnythingOfType(&http.Request{}),
	).Once().Return(response, err)
}

func (m *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	called := m.Called(req)
	return called.Get(0).(*http.Response), called.Error(1)
}
