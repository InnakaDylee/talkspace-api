package databases

import (
	"crypto/tls"
	"log"
	"talkspace-api/app/configs"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

func ConnectElasticsearch() *elasticsearch.Client {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load Elasticsearch configuration: %v", err)
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ELASTICSEARCH.ELASTICSEARCH_URL},
		Username:  config.ELASTICSEARCH.ELASTICSEARCH_USER,
		Password:  config.ELASTICSEARCH.ELASTICSEARCH_PASS,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})

	if err != nil {
		log.Fatalf("failed to connect to Elasticsearch: %v", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("failed to get Elasticsearch info: %v", err)
	}
	defer res.Body.Close()

	logrus.Info("connected to Elasticsearch")

	return es
}
