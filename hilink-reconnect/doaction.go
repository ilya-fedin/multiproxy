package main

import (
	"fmt"
	"strings"
	"encoding/xml"
)

type XmlResponse struct {
	XMLName xml.Name    `xml:"response"`
	Value   string      `xml:",innerxml"`
}

func doConnect(ip, sesInfo, tokInfo string) (string, error) {
	body, verificationToken, err := makePOSTRequest(ip, "api/dialup/mobile-dataswitch", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><request><dataswitch>1</dataswitch></request>`), sesInfo, tokInfo)
	if err != nil {
		return tokInfo, err
	}

	var parsed XmlResponse
	if err := xml.Unmarshal(body, &parsed); err != nil {
		var parsedError XmlError
		if err := xml.Unmarshal(body, &parsedError); err != nil {
			return verificationToken, fmt.Errorf("POST request doConnect() failed:\n%s", body)
		}
		return verificationToken, fmt.Errorf("POST request doConnect() failed: Error %s", parsedError.Code)
	}

	return verificationToken, nil
}

func doDisconnect(ip, sesInfo, tokInfo string) (string, error) {
	body, verificationToken, err := makePOSTRequest(ip, "api/dialup/mobile-dataswitch", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><request><dataswitch>0</dataswitch></request>`), sesInfo, tokInfo)
	if err != nil {
		return tokInfo, err
	}

	var parsed XmlResponse
	if err := xml.Unmarshal(body, &parsed); err != nil {
		var parsedError XmlError
		if err := xml.Unmarshal(body, &parsedError); err != nil {
			return verificationToken, fmt.Errorf("POST request doDisconnect() failed:\n%s", body)
		}
		return verificationToken, fmt.Errorf("POST request doDisconnect() failed: Error %s", parsedError.Code)
	}

	return verificationToken, nil
}

func doReboot(ip, sesInfo, tokInfo string) (string, error) {
	body, verificationToken, err := makePOSTRequest(ip, "api/device/control", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><request><Control>1</Control></request>`), sesInfo, tokInfo)
	if err != nil {
		return tokInfo, err
	}

	var parsed XmlResponse
	if err := xml.Unmarshal(body, &parsed); err != nil {
		var parsedError XmlError
		if err := xml.Unmarshal(body, &parsedError); err != nil {
			return verificationToken, fmt.Errorf("POST request doReboot() failed:\n%s", body)
		}
		return verificationToken, fmt.Errorf("POST request doReboot() failed: Error %s", parsedError.Code)
	}

	return verificationToken, nil
}
