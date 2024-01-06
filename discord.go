package discord

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Discord interface {
	Info(keyVal ...interface{}) error
	Debug(keyVal ...interface{}) error
	Warn(keyVal ...interface{}) error
	Error(keyVal ...interface{}) error
	Fatal(keyVal ...interface{}) error
	Trace(keyVal ...interface{}) error
	Panic(keyVal ...interface{}) error
}

const (
	InfoColor  = 3447003
	DebugColor = 15105570
	WarnColor  = 16776960
	ErrorColor = 15158332
	FatalColor = 10181046
	TraceColor = 9807270
	PanicColor = 10038562
)

type Config struct {
	Webhook string
	Title   string
}

func NewDiscordLogger(loggerConfig *Config) Discord {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: config,
	}

	netClient := &http.Client{
		Transport: transport,
	}

	return &discordLogger{
		webhook:   loggerConfig.Webhook,
		netClient: netClient,
		title:     loggerConfig.Title,
	}
}

type discordLogger struct {
	webhook   string
	netClient *http.Client
	title     string
}

type Fields struct {
	Name  interface{} `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type Embeds struct {
	Description string    `json:"description,omitempty"`
	Color       int       `json:"color,omitempty"`
	Fields      []*Fields `json:"fields,omitempty"`
}

type Params struct {
	Content interface{} `json:"content,omitempty"`
	Embeds  []*Embeds   `json:"embeds"`
}

func (s *discordLogger) Info(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, InfoColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) Debug(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, DebugColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) Warn(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, WarnColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) Error(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, ErrorColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) Fatal(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, FatalColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) Trace(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, TraceColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) Panic(keyVal ...interface{}) error {
	err := s.sendMessage(keyVal, PanicColor)

	if err != nil {
		return err
	}

	return nil
}

func (s *discordLogger) sendMessage(data []interface{}, color int) error {
	if s.webhook == "" {
		return errors.New("webhook is empty")
	}

	_, err := s.netClient.Post(s.webhook, "application/json", prepareData(data, s.title, color))
	if err != nil {
		return err
	}

	return nil
}

func prepareData(keyVal []interface{}, title string, color int) *bytes.Buffer {

	var fields []*Fields
	var value interface{}
	description := ""
	keyValLen := len(keyVal)
	key := 0

	for i := 0; i < keyValLen; i += 2 {
		if i%2 == 0 {
			key = i
		}

		if (i + 1) == keyValLen {
			value = "-"
		} else if fmt.Sprintf("%s", keyVal[key+1]) == "" {
			value = "-"
		} else {
			value = fmt.Sprintf("%s", keyVal[key+1])
		}

		if keyValLen == i+1 {
			fields = append(fields, &Fields{
				Name:  fmt.Sprintf("%s", keyVal[key]),
				Value: "-",
			})

			break
		}

		switch strings.ToUpper(fmt.Sprintf("%s", keyVal[key])) {
		case "DESCRIPTION":
			description = fmt.Sprintf("%s", keyVal[key+1])
			continue
		case "COLOR":
			color, _ = strconv.Atoi(keyVal[key+1].(string))
			continue
		}

		fields = append(fields, &Fields{
			Name:  fmt.Sprintf("%s", keyVal[key]),
			Value: value,
		})

	}

	var embeds []*Embeds

	embeds = append(embeds, &Embeds{
		Description: fmt.Sprintf("%v", description),
		Color:       color,
		Fields:      fields,
	})

	postBody, _ := json.Marshal(Params{
		Content: title,
		Embeds:  embeds,
	})

	responseBody := bytes.NewBuffer(postBody)

	return responseBody
}
