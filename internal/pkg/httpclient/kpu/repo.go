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
)

const (
	host             = "https://sirekap-obj-data.kpu.go.id"
	webHost          = "https://pemilu2024.kpu.go.id/pilpres/hitung-suara"
	hhcwInfoPath     = "/pemilu/hhcw/ppwp"
	ppwtListPath     = "/wilayah/pemilu/ppwp"
	stdJsonExtension = ".json"
)

type repo struct {
	client *http.Client
}

func (r repo) GetPPWTList(req model.PPWTEntity) ([]model.PPWTEntity, error) {
	endpoint, err := url.JoinPath(host, ppwtListPath, req.GetCanonicalCode()+stdJsonExtension)
	if err != nil {
		return nil, err
	}

	body, err := r.jsonRequest(endpoint)
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

func (r repo) GetPageLink(req []string) (string, error) {
	return url.JoinPath(webHost, req...)
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

func (r repo) GetHHWCInfo(req *model.HHCWEntity) error {
	endpoint, err := url.JoinPath(host+"/"+hhcwInfoPath, req.GetCanonicalCode()...)
	if err != nil {
		return err
	}

	endpoint += ".json"

	body, err := r.jsonRequest(endpoint)
	if err != nil {
		return nil
	}

	var respData dto2.HHCWEntity
	if err := json.Unmarshal(body, &respData); err != nil {
		return err
	}

	obj, err := respData.ToModel()
	if err != nil {
		return err
	}

	obj.Code = req.Code
	*req = obj
	return nil
}

func New() dao.KPU {
	return &repo{
		client: http.DefaultClient,
	}
}
