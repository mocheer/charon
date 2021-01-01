package fs

import (
	"encoding/json"
	"io/ioutil"
)

func ReadJSON(path string, e interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	return nil
}
