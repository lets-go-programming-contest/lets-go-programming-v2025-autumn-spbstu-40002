package cbr

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var (
	ErrOpenInputXML       = errors.New("open input xml")
	ErrDecodeInputXML     = errors.New("decode input xml")
	ErrUnsupportedCharset = errors.New("unsupported charset")
)

type Document struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ReadFile(path string) (Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("%s: %w", ErrOpenInputXML, err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		cs := strings.ToLower(strings.TrimSpace(charset))
		switch cs {
		case "utf-8", "utf8", "":
			return input, nil
		case "windows-1251", "windows1251", "cp1251":
			return transform.NewReader(input, charmap.Windows1251.NewDecoder()), nil
		default:
			return nil, fmt.Errorf("%w: %q", ErrUnsupportedCharset, charset)
		}
	}

	var doc Document

	if err := decoder.Decode(&doc); err != nil {
		return Document{}, fmt.Errorf("%s: %w", ErrDecodeInputXML, err)
	}

	return doc, nil
}
