package unmarshalxml

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/Svyatoslav2324/task-3/internal/data"
	"golang.org/x/text/encoding/charmap"
)

func UnMarshalXML(inputFile *os.File, valutes *data.DataStruct) error {

	decoder := xml.NewDecoder(inputFile)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unsupported charset: %s", charset)
		}
	}

	err := decoder.Decode(valutes)
	if err != nil {
		return err
	}

	return nil
}
