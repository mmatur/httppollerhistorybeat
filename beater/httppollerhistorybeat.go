package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"encoding/json"
	"io/ioutil"

	"github.com/mmatur/httppollerhistorybeat/config"
)

type Httppollerhistorybeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Httppollerhistorybeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Httppollerhistorybeat) Run(b *beat.Beat) error {
	logp.Info("httppollerhistorybeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	datas := bt.getDatas()

	ticker := time.NewTicker(1 * time.Millisecond)
	endTicker := time.NewTicker(60 * time.Second)
	i := 0
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			if i < len(datas.Hit.Hits) {
				hit := datas.Hit.Hits[i]
				event := beat.Event{
					Timestamp: hit.Source.Timestamp,
					Fields: common.MapStr{
						"type": b.Info.Name,
						"url":  "https://hub.docker.com/v2/repositories/library/traefik",
						"dockerhub.repository.fullname":   "library/traefik",
						"dockerhub.repository.owner":      "library",
						"dockerhub.repository.name":       "traefik",
						"dockerhub.repository.pull_count": hit.Source.PullCount,
						"dockerhub.repository.star_count": hit.Source.StarCount,
					},
				}
				bt.client.Publish(event)
				logp.Info("Event sent %d", i)
				i++
			}
		case <-endTicker.C:
			bt.Stop()
		}

	}
	return nil
}

func (bt *Httppollerhistorybeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Httppollerhistorybeat) getDatas() Data {
	raw, err := ioutil.ReadFile(bt.config.Path)
	if err != nil {
		logp.Err("Error reading file %s error: %s", bt.config.Path, err.Error())
		fmt.Println(err.Error())
		bt.Stop()
	}

	var datas Data
	err = json.Unmarshal(raw, &datas)
	if err != nil {
		logp.Err("Error Unmarshaling data: %s", err.Error())
		bt.Stop()
	}
	return datas
}

type Data struct {
	Hit Hit `json:"hits"`
}

type Hit struct {
	Hits []Hits `json:"hits"`
}

type Hits struct {
	Source Source `json:"_source"`
}

type Source struct {
	Timestamp time.Time `json:"@timestamp"`
	StarCount int64     `json:"star_count"`
	PullCount int64     `json:"pull_count"`
}

func (p Data) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(bytes)
}
