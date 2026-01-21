package geoblock

import (
	"os"
	"strings"
)


func GeoBlock(blockedGeos string, currentGeo string) {
	if blockedGeos == "" {
		return
	}

	blockedList := strings.Split(blockedGeos, ",")
	currentGeo = strings.ToUpper(strings.TrimSpace(currentGeo))
	
	for _, blockedGeo := range blockedList {
		blockedGeo = strings.ToUpper(strings.TrimSpace(blockedGeo))
		
		if blockedGeo == currentGeo {
			os.Exit(1)
		}
	}
}