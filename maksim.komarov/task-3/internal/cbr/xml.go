package cbr

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var (
	ErrOpenInputXML   = errors.New("open input xml")
	ErrDecodeInputXML = errors.New("decode input xml")
)

type Document struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"-" json:"value"`
	ValueRaw string  `xml:"Value" json:"-"`
}

func LoadXML(path string) (Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("%s: %w", ErrOpenInputXML, err)
	}

	defer func() { _ = f.Close() }()

	dec := xml.NewDecoder(f)

	dec.CharsetReader = func(charset string, in io.Reader) (io.Reader, error) {
		cs := strings.ToLower(strings.TrimSpace(charset))
		switch cs {
		case "", "utf-8", "utf8":
			return in, nil
		case "windows-1251", "windows1251", "cp1251":
			return charmap.Windows1251.NewDecoder().Reader(in), nil
		default:
			return nil, fmt.Errorf("%s: %s", ErrDecodeInputXML, "unsupported charset")
		}
	}

	var d Document

	if err := dec.Decode(&d); err != nil {
		return Document{}, fmt.Errorf("%s: %w", ErrDecodeInputXML, err)
	}

	return d, nil
}
