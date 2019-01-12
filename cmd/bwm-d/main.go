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
import "github.com/rickb777/date/period"

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

type Publisher struct {
	apiToken   string
	deviceName string
}

func EnvIf(varName string, defaultValue string) string {
	result := defaultValue

	if newResult := os.Getenv(varName); "" != newResult {
		result = newResult
	}

	return result
}

func NewPublisher() (*Publisher, error) {
	deviceName := EnvIf("BALENA_DEVICE_NAME_AT_INIT", "DEV")
	apiToken := os.Getenv("UBIDOTS_API_TOKEN")

	return &Publisher{
		deviceName: deviceName,
		apiToken:   apiToken,
	}, nil
}

func (p *Publisher) Publish(b *BandwidthMeasurement) error {
	resp, err := resty.R().
		SetHeader("X-Auth-Token", p.apiToken).
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
		}).Post("https://industrial.api.ubidots.com/api/v1.6/devices/" + p.deviceName)

	if nil == err {
		log.Infof("resp: %+v", resp)
	} else {
		log.Warn("Oops: %s", err)
	}

	return err
}

func (p *Publisher) ParseAndPublishLine(s string) error {
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

	return p.Publish(rec)
}

func BackgroundTask() {
	serverId := EnvIf("SERVER_ID", "13536")

	publisher, err := NewPublisher()

	if nil != err {
		log.Warn("Oops: %s", err)

		return
	}

	cliPath, err := exec.LookPath("speedtest-cli")

	if nil != err {
		log.Warn("Oops: %s", err)

		return
	}

	reportingPeriod, _ := period.MustParse(EnvIf("REPORTING_PERIOD", "P60M")).Duration()

	log.Infof("Using duration: %s", reportingPeriod.String())

	for {
		log.Infof("Using server id: %s", serverId)

		log.Infof("Waiting a minute")
		time.Sleep(1 * time.Minute)

		log.Infof("Running")

		cmdToRun := exec.Command(cliPath, "--csv", "--server", serverId)

		outputBytes, err := cmdToRun.Output()

		if nil != err {
			log.Warnf("Oops: %s", err)
			continue
		}

		commandOutput := string(outputBytes)

		err = publisher.ParseAndPublishLine(commandOutput)

		if nil != err {
			log.Warn("Oops: %s", err)
			continue
		}

		log.Infof("Output: %s", commandOutput)

		time.Sleep(reportingPeriod)
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{})

	log.Infof("Starting")

	BackgroundTask()
}
