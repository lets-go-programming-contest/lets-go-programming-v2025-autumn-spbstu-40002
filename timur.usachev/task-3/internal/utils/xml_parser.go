package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"golang.org/x/net/html/charset"
)

func parseXML(data []byte) (valCurs, error) {
	var root valCurs

	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.CharsetReader = charset.NewReaderLabel

	if err := dec.Decode(&root); err != nil {
		return valCurs{}, fmt.Errorf("%s: %w", ErrXMLParse.Error(), err)
	}

	return root, nil
}
