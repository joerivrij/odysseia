package elastic

import (
	"bytes"
	"encoding/json"
)

func toBuffer(query map[string]interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return buf, err
	}

	return buf, nil
}
