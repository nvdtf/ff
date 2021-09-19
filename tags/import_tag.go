package tags

import (
	"regexp"
)

var knownImports = map[string]string{
	"c1e4f4f4c4257510": "TopShotMarket",
	"329feb3ab062d289": "RaceDay",
	"64f83c60989ce555": "ChainmonstersMarket",
	"3c5959b568896393": "FUSD",
}

func ProcessImportTags(code string) []string {
	r, _ := regexp.Compile("import \\w* from 0x(?P<Address>[0-9a-f]*)")
	matches := r.FindAllStringSubmatch(code, -1)

	returnSet := make(map[string]struct{})
	for i := range matches {
		knownTag := knownImports[matches[i][1]]
		returnSet[knownTag] = struct{}{}
	}

	var returnList []string
	for tag := range returnSet {
		returnList = append(returnList, tag)
	}
	return returnList
}
