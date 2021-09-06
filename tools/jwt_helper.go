package tools

import (
	"btradoc/entities"
)

func MakeDialectPermissionsIntoAbbreviations(dialectPermissions *[]entities.DialectSubdialectDocument) []string {
	var dperms []string
	for _, dp := range *dialectPermissions {
		// convert to runes for non-ascii character
		dialectRunes := []rune(dp.Dialect)
		subdialectRunes := []rune(dp.Subdialect)
		dperm := string(dialectRunes[:3]) + "_" + string(subdialectRunes[:3])
		dperms = append(dperms, dperm)
	}
	return dperms
}
