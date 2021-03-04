package fs

import (
	"encoding/json"
	"os"
)

func ReadJSON(path string, e interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	return nil
}
