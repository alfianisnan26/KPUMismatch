package kpu

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"kawalrealcount/internal/data/dao"
	"kawalrealcount/internal/data/model"
	dto2 "kawalrealcount/internal/pkg/httpclient/kpu/dto"
	"net/http"
	"net/url"
	"time"
)

const (
	host    = "https://sirekap-obj-data.kpu.go.id"
	webHost = "https://pemilu2024.kpu.go.id/pilpres/hitung-suara"

	ppwtListPath = "/wilayah/pemilu/ppwp"
	hhcwInfoPath = "/pemilu/hhcw/ppwp"

	stdJsonExtension = ".json"
)

type repo struct {
	client    *http.Client
	cacheRepo dao.Cache
}

func (r repo) GetPageLink(req model.PPWTEntity) (string, error) {
	return url.JoinPath(webHost, req.GetCanonicalCode())
}

type Param struct {
	CacheRepo dao.Cache
}

func (r repo) jsonRequest(endpoint string) ([]byte, error) {
	res, err := r.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("non-200")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r repo) GetPPWTParent(req model.PPWTEntity) (model.PPWTEntity, error) {
	node := &req
	for i := req.Tingkat - 1; i >= 0; i-- {
		var code = "0"
		if i > 0 {
			code = node.GetCanonicalCodeAll()[i-1]
		}
		parent := model.NewPPWT(code)
		plist, err := r.GetPPWTList(parent)
		if err != nil {
			return model.PPWTEntity{}, err
		}

		for _, entity := range plist {
			if entity.Kode == node.Kode {
				*node = entity

				if parent.Kode != "0" {
					node.Parent = &parent
				} else {
					node.Parent = nil
				}

				node = &parent
				break
			}
		}
	}

	return req, nil
}

func (r repo) GetPPWTList(req model.PPWTEntity) ([]model.PPWTEntity, error) {
	endpoint, err := url.JoinPath(host, ppwtListPath, req.GetCanonicalCode()+stdJsonExtension)
	if err != nil {
		return nil, err
	}

	body, err := r.cacheRepo.Get(endpoint, time.Hour*24*7, r.jsonRequest)
	if err != nil {
		return nil, err
	}

	var respData []dto2.PPWTEntity
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, err
	}

	var data = make([]model.PPWTEntity, len(respData))
	for i, datum := range respData {
		data[i] = datum.ToModel(req)
	}

	return data, nil
}

func (r repo) GetHHCWInfo(req model.PPWTEntity) (model.HHCWEntity, error) {
	endpoint, err := url.JoinPath(host, hhcwInfoPath, req.GetCanonicalCode()+stdJsonExtension)
	if err != nil {
		return model.HHCWEntity{}, err
	}

	body, err := r.cacheRepo.Get(endpoint, 3*time.Hour, r.jsonRequest)
	if err != nil {
		return model.HHCWEntity{}, err
	}

	var respData dto2.HHCWEntity
	if err := json.Unmarshal(body, &respData); err != nil {
		return model.HHCWEntity{}, err
	}

	return respData.ToModel(req)
}

type noCache struct{}

func (n noCache) Get(key string, expiry time.Duration, fallback func(string) ([]byte, error)) ([]byte, error) {
	return fallback(key)
}

func New(param Param) dao.KPU {

	cache := param.CacheRepo
	if cache == nil {
		cache = noCache{}
	}

	return &repo{
		client:    http.DefaultClient,
		cacheRepo: cache,
	}
}
