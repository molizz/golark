package callback

import (
	"encoding/json"

	"github.com/molizz/golark/utils"
	"github.com/pkg/errors"
)

/*

事件订阅回调

请求示例
{
    "uuid": "ae23aeea711e710b359bfd85b9e78ab4",
    "event": {
        "app_id": "cli_9f8ff7a9beb5100e",
        "chat_id": "oc_d8bb949ee92a06299d4a3089d5414f1a",
        "operator": {
            "open_id": "ou_4db0e9e9f1313c180b78bafa1cfb1e4e",
            "user_id": "b7dc42d8"
        },
        "tenant_key": "2c60c8c7f5cfd75e",
        "type": "p2p_chat_create",
        "user": {
            "name": "moli",
            "open_id": "ou_4db0e9e9f1313c180b78bafa1cfb1e4e",
            "user_id": "b7dc42d8"
        }
    },
    "token": "mvvmPUXRiLcXOwAmbWMfgd66AvFokAkB",
    "ts": "1596025030.119123",
    "type": "event_callback"
}
*/

type EventProcessor interface {
	TypeLabel() string
	Process(eventBody []byte) error
}

type ProcessorFinder interface {
	Find(typeName string) EventProcessor
}

type callBack struct {
	UUID  string          `json:"uuid"`
	Type  string          `json:"type"`
	Ts    string          `json:"ts"`
	Token string          `json:"token"`
	Event json.RawMessage `json:"event"`
}

//
type Callback struct {
	eventProcessorFinder ProcessorFinder
}

func NewCallback(callbackFinder ProcessorFinder) *Callback {
	return &Callback{
		eventProcessorFinder: callbackFinder,
	}
}

func (c *Callback) TypeLabel() string {
	return TypeCallback
}

func (c *Callback) Process(params string) (interface{}, error) {
	cb := new(callBack)

	err := json.Unmarshal([]byte(params), cb)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	typeSet := make(map[string]interface{})
	err = json.Unmarshal(cb.Event, &typeSet)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	process := c.eventProcessorFinder.Find(typeSet["type"].(string))
	if process == nil {
		utils.DefaultLog.Printf("not found callback, params: %s\n", params)
		return nil, nil
	} else {
		return nil, process.Process(cb.Event)
	}
}
