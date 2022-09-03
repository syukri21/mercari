package helper

import (
	"fmt"
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

func GetKey(areaType string, key string, isAll bool) string {
	if key == "" && areaType == Province {
		return Province
	}

	if isAll {
		switch areaType {
		case Province:
			return fmt.Sprintf("%s_%s_%s", Province, key[:2], City)
		case City:
			return fmt.Sprintf("%s_%s_%s_%s_%s", Province, key[:2], City, key[:4], District)
		case District:
			return fmt.Sprintf("%s_%s_%s_%s_%s_%s_%s", Province, key[:2], City, key[:4], District, key[:6], SubDistrict)
		case SubDistrict:
			return fmt.Sprintf("%s_%s_%s_%s_%s_%s_%s_%s", Province, key[:2], City, key[:4], District, key[:6], SubDistrict, key)
		}
	} else {
		switch areaType {
		case Province:
			return areaType + key
		case City:
			return fmt.Sprintf("%s_%s_%s%s", Province, key[:2], City, key)
		case District:
			return fmt.Sprintf("%s_%s_%s_%s_%s_%s", Province, key[:2], City, key[:4], District, key)
		case SubDistrict:
			return fmt.Sprintf("%s_%s_%s_%s_%s_%s_%s_%s", Province, key[:2], City, key[:4], District, key[:6], SubDistrict, key)
		}
	}

	return ""
}
