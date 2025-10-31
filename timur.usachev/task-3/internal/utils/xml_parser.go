package utils

import (
	"bytes"
	"encoding/xml"
)

func parseXML(data []byte) (valCurs, error) {
	var root valCurs

	decoder := xml.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&root)
	if err != nil {
		return valCurs{}, err
	}

	return root, nil
}
