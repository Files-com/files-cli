package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go"
	u2f "github.com/marshallbrekka/go-u2fhost"
)

func u2fResponse(paramsSessionCreate files_sdk.SessionCreateParams, responseError files_sdk.ResponseError) (files_sdk.SessionCreateParams, error) {
	u2fSIgnRequests := responseError.Data.U2fSIgnRequests
	request := &u2f.AuthenticateRequest{
		Challenge: u2fSIgnRequests[0].Challenge,
		AppId:     u2fSIgnRequests[0].AppId,
		Facet:     u2fSIgnRequests[0].AppId,
		KeyHandle: u2fSIgnRequests[0].SignRequest.KeyHandle,
	}
	response, err := u2fDeviceInput(request, u2f.Devices())
	responseJson, _ := json.Marshal(response)
	paramsSessionCreate.Otp = string(responseJson)
	paramsSessionCreate.PartialSessionId = responseError.Data.PartialSessionId
	paramsSessionCreate.Password = ""

	return paramsSessionCreate, err
}

func u2fDeviceInput(req *u2f.AuthenticateRequest, devices []*u2f.HidDevice) (*u2f.AuthenticateResponse, error) {
	var openDevices []u2f.Device
	for i, device := range devices {
		err := device.Open()
		if err == nil {
			openDevices = append(openDevices, u2f.Device(devices[i]))
			defer func(i int) {
				devices[i].Close()
			}(i)
			version, err := device.Version()
			if err != nil {
				fmt.Printf("Device version error: %s", err.Error())
			} else {
				fmt.Printf("Device version: %s", version)
			}
		}
	}
	if len(openDevices) == 0 {
		return nil, errors.New("failed to find any devices")
	}
	prompted := false
	timeout := time.After(time.Second * 25)
	interval := time.NewTicker(time.Millisecond * 250)
	defer interval.Stop()
	for {
		select {
		case <-timeout:
			return nil, errors.New("failed to get authentication response after 25 seconds")
		case <-interval.C:
			for _, device := range openDevices {
				response, err := device.Authenticate(req)
				if err == nil {
					return response, nil
				} else if err.Error() == "Device is requesting test of use presence to fulfill the request." && !prompted {
					fmt.Println("\nTouch the flashing U2F device to authenticate")
					prompted = true
				} else {
					fmt.Printf(".")
				}
			}
		}
	}
	return nil, nil
}
