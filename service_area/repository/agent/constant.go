package agent

const (
	IKSGroup    = "IKSGROUP"
	IKSGroupAPI = "https://api.iksgroup.co.id/apilokasi/provinsi"
)

type ProvinceRawData struct {
	Output    string `json:"output"`
	Generator string `json:"generator"`
	Data      []struct {
		ID   string `json:"provinsi_id"`
		Name string `json:"provinsi_nama"`
	} `json:"data"`
}

type CityRawData struct {
	Output    string `json:"output"`
	Generator string `json:"generator"`
	Data      []struct {
		ID   string `json:"kabupatenkota_id"`
		Name string `json:"kabupatenkota_nama"`
	} `json:"data"`
}

type DistrictRaw struct {
	Output    string            `json:"output"`
	Generator string            `json:"generator"`
	Data      []DistrictRawData `json:"data"`
}

type DistrictRawData struct {
	Key  string `json:"kecamatan_id"`
	Name string `json:"kecamatan_nama"`
}

type SubDistrictRawData struct {
	Output    string `json:"output"`
	Generator string `json:"generator"`
	Data      []struct {
		Key  string `json:"kelurahan_id"`
		Name string `json:"kelurahan_nama"`
	} `json:"data"`
}
