package cbr

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var (
	ErrOpenInputXML       = errors.New("open input xml")
	ErrDecodeInputXML     = errors.New("decode input xml")
	ErrUnsupportedCharset = errors.New("unsupported charset")
)

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"-"`
	ValueRaw string  `json:"-"         xml:"Value"`
}

type Document struct {
	Valutes []Currency `xml:"Valute"`
}

func LoadXML(path string) (Document, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return Document{}, fmt.Errorf("%w: %s", ErrOpenInputXML, err.Error())
	}

	defer func() { _ = file.Close() }()

	var doc Document

	dec := xml.NewDecoder(file)
	dec.CharsetReader = func(charset string, rdr io.Reader) (io.Reader, error) {
		switch strings.ToLower(strings.ReplaceAll(charset, "-", "")) {
		case "utf8", "utf-8":
			return rdr, nil
		case "windows1251", "cp1251":
			return charmap.Windows1251.NewDecoder().Reader(rdr), nil
		default:
			return nil, fmt.Errorf("%w: %s", ErrDecodeInputXML, ErrUnsupportedCharset.Error())
		}
	}

	if err := dec.Decode(&doc); err != nil {
		return Document{}, fmt.Errorf("%w: %s", ErrDecodeInputXML, err.Error())
	}

	return doc, nil
}
