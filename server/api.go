package server

import (
	"482.solutions-node-storage/storage"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func NewApi(storageRepository storage.Repository) Api {
	return &HttpApi{
		StorageRepository: storageRepository,
	}
}

type Api interface {
	Get(string) (*bytes.Buffer, error)
	Set(string) (*bytes.Buffer, error)
	GetAll() (*bytes.Buffer, error)
	BuildRecoverResponse() (*bytes.Buffer, error)
	Recover(string)
}

type HttpApi struct {
	StorageRepository storage.Repository
}

func (a *HttpApi) Get(message string) (*bytes.Buffer, error) {
	response := Response{}

	var request Request
	if err := json.NewDecoder(strings.NewReader(message)).Decode(&request); err != nil {
		log.Printf("json decoding error: %s", err)
		return nil, errors.New("wrong message format")
	}

	value, err := a.StorageRepository.Get(request.Key)
	if err != nil {
		log.Printf("storage reading error: %s", err)
		return nil, errors.New("error getting data from the storage")

	}

	response.Value = value
	respBytes, err := json.Marshal(response)

	if err != nil {
		log.Printf("json encoding error: %s", err)
		return nil, errors.New("error result json encoding")
	}

	return bytes.NewBuffer(respBytes), nil
}

func (a *HttpApi) Set(message string) (*bytes.Buffer, error) {
	response := &Response{}

	var request Request
	if err := json.NewDecoder(strings.NewReader(message)).Decode(&request); err != nil {
		log.Printf("json decoding error: %s", err)
		return nil, errors.New("wrong message format")
	}

	if err := a.StorageRepository.Set(request.Key, request.Value); err != nil {
		log.Printf("storage inserting error: %s", err)
		return nil, errors.New("error inserting data into the storage")
	}

	response.Key = request.Key
	response.Value = request.Value

	respBytes, err := json.Marshal(response)

	if err != nil {
		log.Printf("json encoding error: %s", err)
		return nil, errors.New("error result json encoding")

	}

	return bytes.NewBuffer(respBytes), nil
}

func (a *HttpApi) GetAll() (*bytes.Buffer, error) {
	data := a.StorageRepository.GetAll()
	response := make([]*Response, len(data))

	i := 0
	for key, value := range data {
		response[i] = &Response{
			Key:   key,
			Value: value,
		}

		i++
	}

	respBytes, err := json.Marshal(response)

	if err != nil {
		log.Printf("json encoding error: %s", err)
		return nil, errors.New("error result json encoding")
	}

	return bytes.NewBuffer(respBytes), nil
}

func (a *HttpApi) BuildRecoverResponse() (*bytes.Buffer, error) {
	resp, err := a.GetAll()
	if err != nil {
		return nil, errors.New("error getting data from storage")
	}

	respStr := fmt.Sprintf("%s||%s", RECOVERPrefix, resp.String())

	return bytes.NewBufferString(respStr), nil
}

func (a *HttpApi) Recover(cmd string) {
	var data []*Response
	buf := bytes.NewBufferString(cmd)
	if err := json.NewDecoder(buf).Decode(&data); err != nil {
		panic(err)
	}

	for _, item := range data {
		if err := a.StorageRepository.Set(item.Key, item.Value); err != nil {
			log.Println("error restoring form parent node")
			os.Exit(1)
		}
	}
}

type Request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
