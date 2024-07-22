package databases

import (
    "context"
    "log"
    "talkspace-api/app/configs"

    "github.com/olivere/elastic/v7"
)

func ConnectElasticsearch() *elastic.Client {
    config, err := configs.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load Elasticsearch configuration: %v", err)
    }

    client, err := elastic.NewClient(
        elastic.SetURL(config.ELASTICSEARCH.ELASTICSEARCH_URL),
        elastic.SetBasicAuth(config.ELASTICSEARCH.ELASTICSEARCH_USER, config.ELASTICSEARCH.ELASTICSEARCH_PASS),
    )
    if err != nil {
        log.Fatalf("failed to connect to Elasticsearch: %v", err)
    }

    info, code, err := client.Ping(config.ELASTICSEARCH.ELASTICSEARCH_URL).Do(context.Background())
    if err != nil {
        log.Fatalf("failed to ping Elasticsearch: %v", err)
    }

    log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

    return client
}
