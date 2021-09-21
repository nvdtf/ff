package follower

import (
	"regexp"
)

var knownImports = map[string]string{
	"c1e4f4f4c4257510": "TopShotMarket",
	"329feb3ab062d289": "RaceDay",
	"64f83c60989ce555": "ChainmonstersMarket",
	"3c5959b568896393": "FUSD",
}

type CadenceImport struct {
	Contract string
	Address  string
}

func GetImports(code string) []CadenceImport {
	r, _ := regexp.Compile("import (?P<Contract>\\w*) from 0x(?P<Address>[0-9a-f]*)")
	matches := r.FindAllStringSubmatch(code, -1)

	var returnList []CadenceImport
	for i := range matches {
		returnList = append(returnList, CadenceImport{
			Contract: matches[i][1],
			Address:  matches[i][2],
		})
	}

	return returnList
}
