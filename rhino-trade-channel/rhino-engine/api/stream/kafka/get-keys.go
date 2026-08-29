package kafka

import (
	"rhino-common/domain_error"
	kafkautil "rhino-common/utils/kafka"
	"sort"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type respKey struct {
	Key string `json:"system_guid"`
	Seq int64  `json:"req_msg_seq"`
}

func (k *KafkaClient) GetHistoricalSentKeysAndReqMsgSeqs() (keys []string, reqMsgSeqs []int64) {

	_, _, messages, err := kafkautil.GetHistoricMessages(k.brokers, k.respTopic)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to GetHistoricalSentKeys")
		return
	}

	for _, message := range messages {
		k := &respKey{}
		err := json.Unmarshal(message, k)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to Unmarshal json data")
			continue
		}
		if k.Key != "" {
			keys = append(keys, k.Key)
		}
		if k.Seq > 0 {
			reqMsgSeqs = append(reqMsgSeqs, k.Seq)
		}
	}

	if len(reqMsgSeqs) > 0 {
		sort.Slice(reqMsgSeqs, func(i, j int) bool {
			return reqMsgSeqs[i] < reqMsgSeqs[j]
		})
	}

	return
}
