package http_mock

import (
	"encoding/json"

	"github.com/stretchr/testify/mock"
)

type HttpRequestMock struct {
	Mock mock.Mock
}

func (repository *HttpRequestMock) PerformRequest(
	method string,
	url string,
	body []byte,
	bodyQuery []byte,
	headers map[string]string,
) (int, []byte, map[string]interface{}) {
	arguments := repository.Mock.Called(method, url, body, bodyQuery, headers)

	if arguments.Get(0) == nil {
		return 0, nil, map[string]interface{}{}
	} else {
		response, _ := json.Marshal(arguments.Get(0))
		return 200, response, map[string]interface{}{}
	}
}
