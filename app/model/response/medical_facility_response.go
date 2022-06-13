package response

type ResponseHttp struct {
	Meta metaResponse `json:"meta"`
	Data Data         `json:"data"`
}

type metaResponse struct {
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
