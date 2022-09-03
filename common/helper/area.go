package helper

import (
	"strings"
)

const (
	Province    = "provinsi"
	City        = "kabupaten"
	District    = "kecamatan"
	SubDistrict = "kelurahan"
)

const (
	ServicesName = "SERVICE_AREA"
	AreaRedisKey = "INDONESIA"
)

func BuildKey(areaType ...string) string {
	return strings.Join(areaType, "_")
}

func GetKey(areaType string, key string) string {
	if key == "" && areaType == Province {
		return Province
	}
	return areaType + "_" + key
}
