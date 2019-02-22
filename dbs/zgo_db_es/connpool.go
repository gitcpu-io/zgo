package zgo_db_es

import "net/http"

type EsResource struct {
	EsClient *http.Client
}

func NewEsResource() *EsResource {
	return &EsResource{EsClient: &http.Client{}}
}

func (this *EsResource) GetEsClient() *http.Client {
	return this.EsClient
}