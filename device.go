package edgecli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RegisterOption struct {
	Name         string
	HardwareUUID string
	OS           string
}

type RegisterService struct {
	HttpOption
}

func (s *RegisterService) Register(opt *RegisterOption) (*DeviceResponse, error) {
	var url string
	var body map[string]string
	url = s.BaseUrl + "/devices"

	body = map[string]string{
		"name":          opt.Name,
		"hardware_uuid": opt.HardwareUUID,
		"platform":      opt.OS,
	}
	postBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", s.Token)
	resp, _ := HandleCall(req)
	log.Tracef("Register response %+v", resp)
	switch resp.(type) {
	case *SuccessResponse:
		deviceJson, _ := json.Marshal(resp.(*SuccessResponse).Data)
		device := DeviceResponse{}
		if err := json.Unmarshal(deviceJson, &device); err != nil {
			return nil, errors.New(fmt.Sprintf("Fail to unmarshal response's data ,err is %+v", err))
		}
		log.Debugf("Registerdevice result is %+v", device)
		return &device, nil
	case *ErrorResponse:
		return nil, errors.New(fmt.Sprintf("Fail to register device, error message: %s", resp.(*ErrorResponse).Message))
	default:
		return nil, errors.New(fmt.Sprint("This client has some unpredictable problems, please contact the omniedge team."))
	}
}
