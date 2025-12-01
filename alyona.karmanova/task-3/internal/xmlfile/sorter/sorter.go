package sorter

import (
	"sort"

	model "github.com/HuaChenju/task-3/internal/xmlfile/model"
)

func SortValCursByValue(doc *model.ValCurs) {
	sort.Slice(doc.Valutes, func(i, j int) bool {
		return doc.Valutes[i].Value > doc.Valutes[j].Value
	})
}
