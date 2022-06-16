package util

import (
	"encoding/json"
	"go-klikdokter/helper/_struct"
	"io/ioutil"
	"net/http"
)

func CallGetDetailMedicalFacility(uid string) (*ResponseHttp, error) {
	url := _struct.MedicalFacilityDomain + _struct.MedicalFacilityPath + uid

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := ResponseHttp{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
