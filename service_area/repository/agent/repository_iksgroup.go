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
	TotalWorker = 20
)

type IKSRepository struct {
	l *log.Logger
}

func (I *IKSRepository) GetALLAreaData(ctx context.Context, saveFunc model.SaveAreaDataRedis) error {
	wg := sync.WaitGroup{}

	// Get All provinces
	provinces := make(map[string]model.Province)

	err := I.getProvinces(&provinces)
	if err != nil {
		return err
	}

	prov, _ := json.Marshal(provinces)
	p, _ := json.Marshal(prov)
	_ = saveFunc(ctx, model.AreaRedis{
		Key:   helper.BuildKey(constant.Province),
		Value: string(p),
	})

	// dispatcher
	jobs := make(chan func(*map[string]model.Province))

	for i := 0; i < TotalWorker; i++ {
		i := i
		go func(cityJobs <-chan func(*map[string]model.Province), wg *sync.WaitGroup, prov *map[string]model.Province) {
			counter := 1
			for job := range cityJobs {
				func() {
					defer wg.Done()
					job(prov)
					I.l.Printf("Worker no := %v and counter := %v \n", i, counter)
				}()
				counter++
			}
		}(jobs, &wg, &provinces)
	}

	// fetch city
	for _, p := range provinces {
		p := p
		wg.Add(1)
		jobs <- func(prov *map[string]model.Province) {
			err = I.getCities(prov, p.Key)
			cities, _ := json.Marshal(p.Cities)
			_ = saveFunc(ctx, model.AreaRedis{
				Key:   helper.BuildKey(constant.City, p.Key),
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
			jobs <- func(prov *map[string]model.Province) {
				err = I.getDistricts(&provinces, p.Key, c.Key)
				data, _ := json.Marshal(provinces[p.Key].Cities[c.Key])
				_ = saveFunc(ctx, model.AreaRedis{
					Key:   helper.BuildKey(constant.District, p.Key),
					Value: string(data),
				})
			}
		}
	}
	wg.Wait()
	for _, p := range provinces {
		for _, c := range p.Cities {
			for _, d := range c.Districts {
				d := d
				p := p
				c := c
				wg.Add(1)
				jobs <- func(prov *map[string]model.Province) {
					err = I.getSubDistricts(&provinces, p.Key, c.Key, d.Key)
					data, _ := json.Marshal(provinces[p.Key].Cities[c.Key].Districts[d.Key])
					_ = saveFunc(ctx, model.AreaRedis{
						Key:   helper.BuildKey(constant.SubDistrict, p.Key),
						Value: string(data),
					})
				}
			}
		}
	}
	return err

}

func (I *IKSRepository) getProvinces(provinces *map[string]model.Province) (err error) {
	url := fmt.Sprintf("%s/%s", IKSURL, constant.Province)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := ProvinceRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return err
	}

	temp := *provinces
	for _, datum := range results.Data {
		temp[datum.ID] = model.Province{
			Name:   datum.Name,
			Key:    datum.ID,
			Cities: make(map[string]model.City),
		}
	}

	*provinces = temp
	return
}

func (I *IKSRepository) getCities(provinces *map[string]model.Province, key string) (err error) {
	url := fmt.Sprintf("%s/%s?%s=%s", IKSURL, constant.City, constant.Province, key)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := CityRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return err
	}

	temp := *provinces
	for _, datum := range results.Data {
		temp[key].Cities[datum.ID] = model.City{
			Name:      datum.Name,
			Key:       datum.ID,
			Districts: make(map[string]model.District),
		}
	}

	*provinces = temp
	return
}

func (I *IKSRepository) getDistricts(provinces *map[string]model.Province, provinceKey string, cityKey string) (err error) {
	url := fmt.Sprintf("%s/%s?%s=%s", IKSURL, constant.District, constant.City, cityKey)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := DistrictRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return err
	}

	temp := *provinces
	for _, datum := range results.Data {
		temp[provinceKey].Cities[cityKey].Districts[datum.ID] = model.District{
			Name:         datum.Name,
			Key:          datum.ID,
			SubDistricts: make(map[string]model.SubDistrict),
		}
	}

	*provinces = temp
	return
}

func (I *IKSRepository) getSubDistricts(provinces *map[string]model.Province, provinceKey string, cityKey string, districtKey string) (err error) {
	url := fmt.Sprintf("%s/%s?%s=%s", IKSURL, constant.SubDistrict, constant.District, districtKey)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	results := SubDistrictRawData{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return err
	}

	temp := *provinces
	for _, datum := range results.Data {
		temp[provinceKey].Cities[cityKey].Districts[districtKey].SubDistricts[datum.ID] = model.SubDistrict{
			Name: datum.Name,
			Key:  datum.ID,
		}
	}

	*provinces = temp
	return
}

func NewIKSRepository(l *log.Logger) repository.Agent {
	return &IKSRepository{l: l}
}
