package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

type Store struct {
	Email  string `json:"email" `
	Name   string `json:"name"`
	Locale string `json:"locale"`
}

/*
	# curl --location --request GET 'http://igw-traefik.internal.svc.cluster.local/service/koala/api/internal/store'  --header 'Store-Id: 8735'
	{"address":"","city":"","company_name":"","country_code":"","created_at":"2021-12-29 18:02:50","currency":"USD","email":"hujuan+2@shoplazza.com",
	"financial_email":"","hour":8,"icon":{"src":"","alt":"","path":""},"id":"69e7fc33-e5fe-4ccb-b31e-bbdeab5ad926",
	"locale":"en-US","money_format":"amount","name":"xx-15","phone":"","province_code":"ALL","service_email":"",
	"store_id":"8735","symbol":"$","symbol_left":"$","symbol_right":"","time_zone":"+0800","updated_at":"2021-12-29 18:02:50","url":"","zip":""}
*/
func GetStore(StoreId string) (*Store, error) {

	if StoreId == "" || StoreId == "0" {
		return &Store{}, nil
	}
	Store := &Store{}
	url := fmt.Sprintf("%s/api/internal/store", config.Cfg.StoreServiceKoala)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Store-Id", StoreId)

	resp, err := utils.Httpclient.Do(req)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get store")
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Read response body err: %s", err)
		return nil, err
	}

	err = json.Unmarshal(content, Store)
	if err != nil {
		logrus.WithError(err).Errorln("failed to unmarshal body to struct")
		return nil, err
	}

	return Store, nil
}

type StoreInfo struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	UserLocale string    `json:"user_locale"`
	CreatedAt  time.Time `json:"created_at"`
}

func GetStoreFromTotoro(StoreId string) (*StoreInfo, error) {
	store := &StoreInfo{}
	url := fmt.Sprintf("%s/api/store/%s", config.Cfg.StoreService, StoreId)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := utils.Httpclient.Do(req)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get store")
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Read response body err: %s", err)
		return nil, err
	}

	err = json.Unmarshal(content, store)
	if err != nil {
		logrus.WithError(err).Errorln("failed to unmarshal body to struct")
		return nil, err
	}

	return store, nil
}
