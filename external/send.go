package external

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Micah-Shallom/departure-times/utility"
)

type SendRequestObject struct {
	Name         string
	Logger       *utility.Logger
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	Data         interface{}
	DecodeMethod string
	UrlPrefix    string
}

func GetNewSendRequestObject(logger *utility.Logger, name, path, method, urlPrefix, decodeMethod string, headers map[string]string, successCode int, data interface{}) *SendRequestObject {
	return &SendRequestObject{
		Logger:       logger,
		Name:         name,
		Path:         path,
		Method:       method,
		UrlPrefix:    urlPrefix,
		DecodeMethod: decodeMethod,
		Headers:      headers,
		SuccessCode:  successCode,
		Data:         data,
	}
}

func (r *SendRequestObject) SendRequest(response interface{}) error {
	var (
		data   = r.Data
		logger = r.Logger
		name   = r.Name
		err    error
	)

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		logger.Error("Error encoding data for request %s: %s", name, err)
	}

	logger.Info("before prefix", name, r.Path, data, buf)
	if r.UrlPrefix != "" {
		r.Path += r.UrlPrefix
	}
	logger.Info("after prefix", name, r.Path, data, buf)

	var req *http.Request
	client := &http.Client{}
	if r.Method == http.MethodGet {
		req, err = http.NewRequest(r.Method, r.Path, nil)
	} else {
		switch r.Headers["Content-Type"] {
		case "application/x-www-form-urlencoded":
			req, err = http.NewRequest(r.Method, r.Path, data.(io.Reader))
		default:
			req, err = http.NewRequest(r.Method, r.Path, buf)
		}
	}
	if err != nil {
		logger.Error("Error creating request %s: %s", name, err)
	}

	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	logger.Info("request", name, r.Path, r.Method, r.Headers)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("client do", name, err.Error())
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error reading response body", name, err.Error())
		return err
	}

	if r.DecodeMethod == "xml" {
		err = xml.Unmarshal(body, response)
		if err != nil {
			logger.Error("json decoding error", name, err.Error())
			return err
		}
	} else {
		err = json.Unmarshal(body, response)
		if err != nil {
			logger.Error("json decoding error", name, err.Error())
			return err
		}
	}

	defer res.Body.Close()
	if res.StatusCode == r.SuccessCode {
		return nil
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("external requests error for request %v, code %v", name, strconv.Itoa(res.StatusCode))
	}

	return nil
}
