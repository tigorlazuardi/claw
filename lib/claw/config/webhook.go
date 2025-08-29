package config

import (
	"errors"
	"net/http"
	"strings"
)

type Webhooks struct {
	JobStart             []Webhook `koanf:"job_start"`
	JobFinished          []Webhook `koanf:"job_finished"`
	ImageDownloadSuccess []Webhook `koanf:"image_downloaded"`
	ImageDownloadFailed  []Webhook `koanf:"image_failed"`
	ImageAssigned        []Webhook `koanf:"image_assigned"`
}

func (hooks Webhooks) ValidateAndNormalize() error {
	var e error
	allHooks := [][]Webhook{
		hooks.JobStart,
		hooks.JobFinished,
		hooks.ImageDownloadSuccess,
		hooks.ImageDownloadFailed,
		hooks.ImageAssigned,
	}
	for _, hookList := range allHooks {
		for i := range hookList {
			if err := hookList[i].ValidateAndNormalize(); err != nil {
				e = errors.Join(e, err)
			}
		}
	}
	return e
}

type Webhook struct {
	URL     string  `koanf:"url"`
	Method  string  `koanf:"method"`
	Headers Headers `koanf:"headers"`
}

func (web *Webhook) ValidateAndNormalize() error {
	if len(web.URL) == 0 {
		return errors.New("webhook URL cannot be empty")
	}
	if len(strings.TrimSpace(web.Method)) == 0 {
		web.Method = http.MethodPost
	}
	if web.Headers == nil {
		web.Headers = make(Headers)
	}
	return nil
}

type Headers map[string]string

func (headers Headers) ToHTTPHeader() http.Header {
	httpHeaders := make(http.Header)
	headers.Apply(httpHeaders)
	return httpHeaders
}

func (headers Headers) Apply(header http.Header) {
	for key, value := range headers {
		header.Set(key, value)
	}
}
