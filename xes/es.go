package xes

import (
	"github.com/olivere/elastic/v7"
)

var (
	_configs map[string]Config
	_clients map[string]*elastic.Client
)

func Init(configs ...Config) error {
	_configs = make(map[string]Config, 8)
	for _, config := range configs {
		conf := config.Default()
		_configs[conf.Id] = conf
	}

	_clients = make(map[string]*elastic.Client, 8)
	for _, config := range _configs {
		client, err := NewClient(config)
		if err != nil {
			Finally()
			return err
		}
		_clients[config.Id] = client
	}

	return nil
}

func NewClient(config Config) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(config.Hosts...),
		elastic.SetHttpClient(config.Client),
		elastic.SetBasicAuth(config.Username, config.Password),
		elastic.SetInfoLog(config.Logger),
		elastic.SetErrorLog(config.Logger),
		elastic.SetTraceLog(config.Logger),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
}

func Finally() {
	for _, client := range _clients {
		client.Stop()
	}
}

func Get(id string) *elastic.Client {
	return _clients[id]
}

func GetDefault() *elastic.Client {
	return Get(DefaultId)
}
