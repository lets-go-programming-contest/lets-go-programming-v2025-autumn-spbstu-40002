package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/xanaoxx/task-3/internal/xml"
)

type CurrencyRecord struct {
	NumCode int `json:"num_code"`
	CharCode string `json:"char_code"`
	Value float64 `json:"value"`
}

func SaveCurrencies(data *xml.Currencies, path string) error {
	dir:=filepath.Dir(path)
	if err:=os.MkdirAll(dir,0750);err!=nil{
		return fmt.Errorf("create dir: %w",err)
	}
	file,err:=os.Create(path)
	if err!=nil{
		return fmt.Errorf("create file: %w",err)
	}
	defer func(){ 
		if err:=file.Close();err!=nil{
			fmt.Printf("close error: %v\n",err)
		}
	}()
	records:=make([]CurrencyRecord,len(data.Currencies))
	for i,c:=range data.Currencies{
		v,err:=c.ToFloat()
		if err!=nil{
			return fmt.Errorf("convert: %w",err)
		}
		records[i]=CurrencyRecord{NumCode:c.NumCode,CharCode:c.CharCode,Value:v}
	}
	enc:=json.NewEncoder(file)
	enc.SetIndent("","  ")
	if err:=enc.Encode(records);err!=nil{
		return fmt.Errorf("encode: %w",err)
	}
	return nil
}
