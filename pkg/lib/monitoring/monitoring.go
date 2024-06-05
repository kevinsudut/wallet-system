package monitoring

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

type monitoring struct {
	Name    string      `json:"name"`
	Latency string      `json:"latency"`
	Content interface{} `json:"content"`
}

func RecordMonitoring(name string, start time.Time, content interface{}) {
	monitoring, err := jsoniter.MarshalToString(monitoring{
		Name:    name,
		Latency: fmt.Sprintf("%dµs", time.Since(start)),
		Content: content,
	})
	if err != nil {
		log.Warnln("failed to record monitoring")
		return
	}

	log.Infoln(monitoring)
}
