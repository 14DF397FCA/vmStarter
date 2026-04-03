package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type YCVM struct {
	ID          string `json:"id"`
	FolderID    string `json:"folderId"`
	CreatedAt   string `json:"createdAt"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Labels      string `json:"labels"`
	ZoneID      string `json:"zoneId"`
	PlatformID  string `json:"platformId"`
	Status      string `json:"status"`
}

func vmStart(vmId string) (string, error) {
	url := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/instances/%s:start", vmId)
	req, err := makePostRequest(url)
	if err != nil {
		return "", err
	}
	resp, err := doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var vm YCVM
	err = json.NewDecoder(resp.Body).Decode(&vm)
	if err != nil {
		log.Errorf("Error unmarshalling response: %s", err)
		return "", err
	}
	return vm.Status, nil
}

func vmGetStatus(vmId string) (string, error) {
	url := fmt.Sprintf("https://compute.api.cloud.yandex.net/compute/v1/instances/%s", vmId)

	req, err := makeGetRequest(url)
	if err != nil {
		return "", err
	}
	resp, err := doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var vm YCVM
	err = json.NewDecoder(resp.Body).Decode(&vm)
	if err != nil {
		log.Errorf("Error unmarshalling response: %s", err)
		return "", err
	}
	return vm.Status, nil
}

func vmIsRunning(status string) bool {
	if status == "RUNNING" {
		return true
	}
	return false
}
