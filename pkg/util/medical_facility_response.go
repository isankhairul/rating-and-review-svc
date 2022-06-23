package util

type ResponseHttp struct {
	Meta MetaResponse `json:"meta"`
	Data Data         `json:"data"`
}

type MetaResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type Data struct {
	Record Record `json:"record"`
}

type Record struct {
	Name string `json:"name"`
}
