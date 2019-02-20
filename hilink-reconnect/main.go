package main

import (
	"fmt"
	"log"
	"flag"
	"time"
	"os"
	"strings"
	"strconv"
	"io/ioutil"
)

type DeviceLoggerWriter struct {
	IP string
}

func (writer DeviceLoggerWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s %s %s", time.Now().Format("2006/01/02 15:04:05"), writer.IP, bytes)
}

type Instance struct {
	Gateway string
	Interface string
	Mark int64
	Table int
	UID int
	Port int
	User string
	Password string
	ReconnectMethod string
	ReconnectInterval int
}

func main() {
	var defaultConfigFile string
	var configFile string
	var instances []Instance
	var instancesBytes []byte
	var err error
	docker := os.Getenv("DOCKER") == "true"
	finished := make(chan bool)

	if docker {
		defaultConfigFile = "/etc/instances"
	} else {
		defaultConfigFile = "/opt/multiproxy/instances"
	}

	flag.StringVar(&configFile, "config", "", "set configuration file (default: " + defaultConfigFile + ")")
	flag.Parse()

	if len(configFile) == 0 {
		configFile = defaultConfigFile
	}

	instancesBytes, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	instancesStrings := strings.Split(string(instancesBytes), "\n")

	for _, element := range instancesStrings {
		var err error

		if len(element) == 0 || element[0] == '#' {
			continue
		}

		instanceStrings := strings.Split(element, "\t")
		instance := *new(Instance)
		instance.Gateway = instanceStrings[0]
		instance.Interface = instanceStrings[1]
		instance.Mark, err = strconv.ParseInt(instanceStrings[2], 0, 64)
		if err != nil {
			log.Print(err)
		}
		instance.Table, err = strconv.Atoi(instanceStrings[3])
		if err != nil {
			log.Print(err)
		}
		instance.UID, err = strconv.Atoi(instanceStrings[4])
		if err != nil {
			log.Print(err)
		}
		instance.Port, err = strconv.Atoi(instanceStrings[5])
		if err != nil {
			log.Print(err)
		}
		instance.User = instanceStrings[6]
		instance.Password = instanceStrings[7]
		instance.ReconnectMethod = instanceStrings[8]
		instance.ReconnectInterval, err = strconv.Atoi(instanceStrings[9])
		if err != nil {
			log.Print(err)
		}
		instances = append(instances, instance)
	}

	for index, element := range instances {
		if len(element.ReconnectMethod) == 0 {
			continue
		}

		if len(element.Gateway) == 0 {
			log.Printf("Instance %d error: Instance does not have gateway IP address", index)
			continue
		}

		go func(element Instance) {
			deviceLoggerWriter := new(DeviceLoggerWriter)
			deviceLoggerWriter.IP = element.Gateway

			deviceLogger := log.New(deviceLoggerWriter, "", 0)

			switch element.ReconnectMethod {
			case "reconnect", "reboot":
			default:
				deviceLogger.Print("Error: Unknown reconnect method")
				return
			}

			deviceLogger.Print("Huawei HiLink Device Reconnector started successfully")
			deviceLogger.Printf("Method: %s, Interval: %d", element.ReconnectMethod, element.ReconnectInterval)

			for {
				time.Sleep(time.Duration(element.ReconnectInterval) * time.Second)

				sesInfo, tokInfo, err := getTokens(element.Gateway)
				if err != nil {
					deviceLogger.Print(err)
					continue
				}

				switch element.ReconnectMethod {
				case "reconnect":
					connectionStatus, err := getConnectionStatus(element.Gateway, sesInfo)
					if err != nil {
						deviceLogger.Print(err)
						continue
					}

					switch connectionStatus {
					case 0:
						deviceLogger.Print("Error: Unknown connecion status")
						continue
					case 900:
						for connectionStatus != 901 {
							connectionStatus, err = getConnectionStatus(element.Gateway, sesInfo)
							if err != nil {
								deviceLogger.Print(err)
								break
							}
						}

						tokInfo, err = doDisconnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}


						connectionStatus, err = getConnectionStatus(element.Gateway, sesInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}

						for connectionStatus != 902 {
							connectionStatus, err = getConnectionStatus(element.Gateway, sesInfo)
							if err != nil {
								deviceLogger.Print(err)
								break
							}
						}

						tokInfo, err = doConnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}
					case 901:
						tokInfo, err = doDisconnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}


						connectionStatus, err = getConnectionStatus(element.Gateway, sesInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}

						for connectionStatus != 902 {
							connectionStatus, err = getConnectionStatus(element.Gateway, sesInfo)
							if err != nil {
								deviceLogger.Print(err)
								break
							}
						}

						tokInfo, err = doConnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}
					case 902:
						tokInfo, err = doConnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}
					case 903:
						for connectionStatus != 902 {
							connectionStatus, err = getConnectionStatus(element.Gateway, sesInfo)
							if err != nil {
								deviceLogger.Print(err)
								break
							}
						}

						tokInfo, err = doConnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}
					default:
						tokInfo, err = doConnect(element.Gateway, sesInfo, tokInfo)
						if err != nil {
							deviceLogger.Print(err)
							continue
						}
					}

					deviceLogger.Print("Device successfully reconnected")
				case "reboot":
					tokInfo, err = doReboot(element.Gateway, sesInfo, tokInfo)
					if err != nil {
						deviceLogger.Print(err)
						continue
					}

					deviceLogger.Print("Device successfully rebooted")
				}
			}
		}(element)
	}

	<- finished
}
