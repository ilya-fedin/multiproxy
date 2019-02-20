package main

import (
	"fmt"
	"strings"
	"encoding/xml"
)

type XmlError struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code"`
}

type XmlSesTokInfo struct {
	XMLName xml.Name `xml:"response"`
	SesInfo string   `xml:"SesInfo"`
	TokInfo string   `xml:"TokInfo"`
}

type XmlMonitoringStatus struct {
	XMLName              xml.Name `xml:"response"`
	ConnectionStatus     int      `xml:"ConnectionStatus"`
}

func getTokens(ip string) (string, string, error) {
	body, err := makeGETRequest(ip, "api/webserver/SesTokInfo", "")
	if err != nil {
		return "", "", err
	}

	var parsed XmlSesTokInfo
	if err := xml.Unmarshal(body, &parsed); err != nil {
		var parsedError XmlError
		if err := xml.Unmarshal(body, &parsedError); err != nil {
			return "", "", fmt.Errorf("GET request getTokens() failed:\n%s", body)
		}
		return "", "", fmt.Errorf("GET request getTokens() failed: Error %s", parsedError.Code)
	}

	if strings.Contains(parsed.SesInfo, "=") {
		parsed.SesInfo = strings.Split(parsed.SesInfo, "=")[1]
	}

	return parsed.SesInfo, parsed.TokInfo, nil
}

func getConnectionStatus(ip string, sesInfo string) (int, error) {
	body, err := makeGETRequest(ip, "api/monitoring/status", sesInfo)
	if err != nil {
		return 0, err
	}

	var parsed XmlMonitoringStatus
	if err := xml.Unmarshal(body, &parsed); err != nil {
		var parsedError XmlError
		if err := xml.Unmarshal(body, &parsedError); err != nil {
			return 0, fmt.Errorf("GET request getConnectionStatus() failed:\n%s", body)
		}
		return 0, fmt.Errorf("GET request getConnectionStatus() failed: Error %s", parsedError.Code)
	}

	return parsed.ConnectionStatus, nil
}
