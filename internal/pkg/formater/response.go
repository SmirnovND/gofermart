package formater

import (
	"encoding/json"
	"fmt"
)

func JSONResponse(data interface{}) ([]byte, error) {
	JSONResponse, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metric to JSON")
	}

	return JSONResponse, nil
}
