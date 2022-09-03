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
	districts := make([]model.AreaRedis, 0)

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
					defer wg.Done()
					job(prov, sub)
				}()
			}
		}(jobs, &wg, &provinces, &districts)
	}

	// fetch city
	for _, p := range data {
		p := p
		wg.Add(1)
		jobs <- func(prov *map[string]model.Province, subDistrict *[]model.AreaRedis) {
			data, _ = I.getCities(p.Key)
			cities, _ := json.Marshal(data)
			_ = saveFunc(ctx, model.AreaRedis{
				Key:   helper.BuildKey(constant.Province, p.Key),
				Value: string(cities),
			})
		}
	}
	wg.Wait()

	// fetch districts
	for _, p := range provinces {
		for _, citi := range p.Cities {
			p := p
			c := citi
			wg.Add(1)
			jobs <- func(prov *map[string]model.Province, subDistrict *[]model.AreaRedis) {
				rawData, _ := I.getDistricts(c.Key)
				*subDistrict = append(districts, rawData...)
				data, _ := json.Marshal(rawData)
				_ = saveFunc(ctx, model.AreaRedis{
					Key:   helper.BuildKey(constant.District, p.Key),
					Value: string(data),
				})
			}
		}
	}

	wg.Wait()

	for _, d := range districts {
		wg.Add(1)
		jobs <- func(prov *map[string]model.Province, subDistrict *[]model.AreaRedis) {
			rawData, _ := I.getSubDistricts(d.Key)
			data, _ := json.Marshal(rawData)
			_ = saveFunc(ctx, model.AreaRedis{
				Key:   helper.BuildKey(constant.SubDistrict, d.Key),
				Value: string(data),
			})
		}
	}

	wg.Wait()
	return err

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
