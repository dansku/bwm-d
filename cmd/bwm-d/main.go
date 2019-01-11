package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type BandwidthMeasurement struct {
	ServerID   string
	Sponsor    string
	Servername string
	Timestamp  string
	Distance   float64
	Ping       float64
	Download   float64
	Upload     float64
	Share      string
	IPAddress  string
}

type ValueHolder struct {
	Value float64 `json:"value"`
}

func (b *BandwidthMeasurement) Publish() error {
	deviceName := "DEV"
	apiToken := os.Getenv("UBIDOTS_API_TOKEN")

	if newDeviceName := os.Getenv("BALENA_DEVICE_NAME_AT_INIT"); "" != newDeviceName {
		deviceName = newDeviceName
	}

	resp, err := resty.R().
		SetHeader("X-Auth-Token", apiToken).
		SetHeader("Content-Type", "application/json").
		SetBody(struct {
			Distance ValueHolder `json:"distance"`
			Ping     ValueHolder `json:"ping"`
			Upload   ValueHolder `json:"upload"`
			Download ValueHolder `json:"download"`
		}{
			Distance: ValueHolder{Value: b.Distance},
			Ping:     ValueHolder{Value: b.Ping},
			Upload:   ValueHolder{Value: b.Upload},
			Download: ValueHolder{Value: b.Download},
		}).Post("https://industrial.api.ubidots.com/api/v1.6/devices/" + deviceName)

	if nil == err {
		log.Infof("resp: %+v", resp)
	}

	return err
}

func parseAndPublishLine(s string) error {
	elements := strings.SplitN(strings.Trim(s, "\r\n "), ",", 10)

	atof := func(s string) float64 {
		v, _ := strconv.ParseFloat(s, 64)

		return v
	}

	rec := &BandwidthMeasurement{
		ServerID:   elements[0],
		Sponsor:    elements[1],
		Servername: elements[2],
		Timestamp:  elements[3],
		Distance:   atof(elements[4]),
		Ping:       atof(elements[5]),
		Download:   atof(elements[6]),
		Upload:     atof(elements[7]),
		Share:      elements[8],
		IPAddress:  elements[9],
	}

	log.Infof("rec: %+v", rec)

	return rec.Publish()
}

func BackgroundTask() {
	serverId := "13536"

	if newServerId := os.Getenv("SERVER_ID"); "" != newServerId {
		serverId = newServerId
	}

	for {
		log.Infof("Using server id: %s", serverId)

		log.Infof("Waiting a minute")
		time.Sleep(1 * time.Minute)

		log.Infof("Running")

		cmdToRun := exec.Command("/usr/local/bin/speedtest-cli", "--csv", "--server", serverId)

		outputBytes, err := cmdToRun.Output()

		if nil != err {
			log.Warnf("Oops: %s", err)
			continue
		}

		commandOutput := string(outputBytes)

		err = parseAndPublishLine(commandOutput)

		if nil != err {
			log.Warn("Oops: %s", err)
			continue
		}

		log.Infof("Output: %s", commandOutput)

		time.Sleep(1 * time.Hour)
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{})

	log.Infof("Starting")

	BackgroundTask()
}
