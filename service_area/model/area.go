package model

import "context"

type Area struct {
	Provinces []Province `json:"-"`
}

type Province struct {
	Name   string
	Key    string
	Cities map[string]City `json:"-"`
}

type City struct {
	Name      string
	Key       string
	Districts map[string]District `json:"-"`
}

type District struct {
	Name         string
	Key          string
	SubDistricts map[string]SubDistrict `json:"-"`
}

type SubDistrict struct {
	Name string
	Key  string
}

type AreaRedis struct {
	Key   string `json:"Key"`
	Value string `json:"Name"`
}

type AreaData struct {
	Key   string               `json:"Key"`
	Value string               `json:"Name"`
	Data  map[string]AreaRedis `json:"data"`
}

type SaveAreaDataRedis = func(ctx context.Context, area AreaRedis) error

type GetAreaInfoRequest struct {
	AreaType string `json:"AreaType" from:"area_type" validate:"required"`
	Key      string `json:"Key" from:"key"`
}
type GetAreaInfoResponse struct {
}
