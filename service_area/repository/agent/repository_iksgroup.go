package agent

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/service_area/constant"
	"github.com/syukri21/mercari/service_area/model"
	"github.com/syukri21/mercari/service_area/repository"
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	IKSURL      = "https://api.iksgroup.co.id/apilokasi"
	TotalWorker = 500
)

type IKSRepository struct {
	l *log.Logger
}

func (I *IKSRepository) GetALLAreaData(ctx context.Context, saveFunc model.SaveAreaDataRedis) error {
	wg := sync.WaitGroup{}

	// Get All provinces
	provinces := make(map[string]model.Province)
	areaData := make([]model.AreaRedis, 0)

	data, err := I.getProvinces()
	if err != nil {
		return err
	}
	prov, _ := json.Marshal(data)
	_ = saveFunc(ctx, model.AreaRedis{
		Key:   constant.Province,
		Value: string(prov),
	})

	// dispatcher
	type JobFunc func(*map[string]model.Province, *[]model.AreaRedis)
	jobs := make(chan JobFunc)

	for i := 0; i < TotalWorker; i++ {
		go func(jobs <-chan JobFunc, wg *sync.WaitGroup, prov *map[string]model.Province, sub *[]model.AreaRedis) {
			for job := range jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("Recovered in f", r)
						}
						wg.Done()
					}()
					job(prov, sub)
				}()
			}
		}(jobs, &wg, &provinces, &areaData)
	}

	I.l.Printf("Fetch City length: %v", len(data))
	// fetch city
	wg.Add(len(data))
	for _, p := range data {
		p := p
		jobs <- func(prov *map[string]model.Province, subDistrict *[]model.AreaRedis) {
			data, _ = I.getCities(p.Key)
			*subDistrict = append(areaData, data...)
			cities, _ := json.Marshal(data)
			_ = saveFunc(ctx, model.AreaRedis{
				Key:   helper.BuildKey(constant.Province, p.Key),
				Value: string(cities),
			})
		}
	}

	wg.Wait()
	I.l.Printf("Fetch Success length: %v")

	// fetch areaData
	I.l.Printf("Fetch Districts length: %v", len(areaData))
	tempAreaData := areaData
	areaData = make([]model.AreaRedis, 0)
	wg.Add(len(tempAreaData))
	for _, p := range tempAreaData {
		p := p
		jobs <- func(prov *map[string]model.Province, subDistrict *[]model.AreaRedis) {
			rawData, _ := I.getDistricts(p.Key)
			*subDistrict = append(areaData, rawData...)
			data, _ := json.Marshal(rawData)
			_ = saveFunc(ctx, model.AreaRedis{
				Key:   helper.BuildKey(constant.City, p.Key),
				Value: string(data),
			})
		}
	}

	wg.Wait()
	I.l.Printf("Fetch Districts Success")

	I.l.Printf("Fetch SubDistrict length: %v", len(areaData))
	wg.Add(len(areaData))
	for _, d := range areaData {
		jobs <- func(prov *map[string]model.Province, subDistrict *[]model.AreaRedis) {
			rawData, _ := I.getSubDistricts(d.Key)
			data, _ := json.Marshal(rawData)
			_ = saveFunc(ctx, model.AreaRedis{
				Key:   helper.BuildKey(constant.District, d.Key),
				Value: string(data),
			})
		}
	}

	wg.Wait()
	I.l.Printf("Fetch Districts Success")
	return err

}

func (I *IKSRepository) GetAreaData(ctx context.Context, key string, areaType string) (result []model.AreaRedis, err error) {
	switch areaType {
	case constant.Province:
		return I.getCities(key)
	case constant.City:
		return I.getDistricts(key)
	case constant.District:
		return I.getSubDistricts(key)
	}
	return result, err
}

func (I *IKSRepository) getProvinces() (data []model.AreaRedis, err error) {
	url := fmt.Sprintf("%s/%s", IKSURL, constant.Province)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return data, err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return data, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := ProvinceRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return data, err
	}

	for _, d := range results.Data {
		data = append(data, model.AreaRedis{
			Key:   d.ID,
			Value: d.Name,
		})
	}
	return data, nil
}

func (I *IKSRepository) getCities(key string) (data []model.AreaRedis, err error) {
	url := fmt.Sprintf("%s/%s?%s=%s", IKSURL, constant.City, constant.Province, key)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return data, err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return data, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := CityRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return data, err
	}

	for _, datum := range results.Data {
		data = append(data, model.AreaRedis{
			Value: datum.Name,
			Key:   datum.ID,
		})
	}

	return
}

func (I *IKSRepository) getDistricts(cityKey string) (data []model.AreaRedis, err error) {
	url := fmt.Sprintf("%s/%s?%s=%s", IKSURL, constant.District, constant.City, cityKey)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return data, err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return data, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := DistrictRaw{}
	err = json.NewDecoder(resp.Body).Decode(&results)

	for _, datum := range results.Data {
		data = append(data, model.AreaRedis{
			Value: datum.Name,
			Key:   datum.Key,
		})
	}

	return
}

func (I *IKSRepository) getSubDistricts(districtKey string) (data []model.AreaRedis, err error) {
	url := fmt.Sprintf("%s/%s?%s=%s", IKSURL, constant.SubDistrict, constant.District, districtKey)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return data, err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return data, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := SubDistrictRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return data, err
	}

	for _, datum := range results.Data {
		data = append(data, model.AreaRedis{
			Value: datum.Name,
			Key:   datum.Key,
		})
	}
	return data, nil
}

func NewIKSRepository(l *log.Logger) repository.Agent {
	return &IKSRepository{l: l}
}
