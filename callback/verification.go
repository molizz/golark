package callback

import (
	"encoding/json"

	"github.com/molizz/golark/utils"
	"github.com/pkg/errors"
)

// {"challenge":"7379fb3e-c512-4118-926b-8f2540b40a5e","token":"UyKEV6uP5h4CruMmxigfkgke1vmL1yWh","type":"url_verification"}

type verification struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type Verification struct {
}

func NewVerification() *Verification {
	return &Verification{}
}

func (v *Verification) TypeLabel() string {
	return TypeVerification
}

func (v *Verification) Process(params string) (interface{}, error) {
	ver := new(verification)

	err := json.Unmarshal([]byte(params), ver)
	if err != nil {
		utils.DefaultLog.Printf("lark event url verify was err: %+v params: %s\n", err, params)
		return nil, errors.WithStack(err)
	}

	result := map[string]string{
		"challenge": ver.Challenge,
	}
	return result, nil
}
