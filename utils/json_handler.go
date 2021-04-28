package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadToMap(resp *http.Response, maps *map[string]interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(string(body)), &maps)
	return nil
}
