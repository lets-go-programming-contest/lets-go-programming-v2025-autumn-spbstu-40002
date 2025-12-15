package unmarshalxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Svyatoslav2324/task-3/internal/structures"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnsuppotedCharset error = errors.New("unsupported charset")

func UnMarshalXML(inputFile *os.File, valutes *structures.ValuteStruct) error {
	decoder := xml.NewDecoder(inputFile)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, ErrUnsuppotedCharset
		}
	}

	err := decoder.Decode(valutes)
	if err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	return nil
}
