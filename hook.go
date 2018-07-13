package logrus_forward

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
	"bytes"
	"errors"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

type ForwardHook struct {
	url    string
	tag    string
	levels []logrus.Level
}

type logEntry struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func New(url, tag string) (*ForwardHook) {
	hook := &ForwardHook{
		url:    url,
		tag:    tag,
		levels: defaultLevels,
	}
	return hook
}

func (hook *ForwardHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *ForwardHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}

func (hook *ForwardHook) Fire(entry *logrus.Entry) error {
	message, err := entry.String()
	if err != nil {
		return err
	}

	data := &logEntry{
		Tag:     hook.tag,
		Message: message,
	}

	if err := send(hook.url, data); err != nil {
		return err
	}

	return nil
}

func send(url string, entry *logEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("failed sending log to specified endpoint")
	}

	return nil
}
