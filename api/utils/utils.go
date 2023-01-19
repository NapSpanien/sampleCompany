package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//Contains check if slice of strings contains string
func Contains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}

	return false, 0
}

func NotifyAdmin(level, employee, message string) error {
	var err error
	type output struct {
		Level                string `json:"level"`
		EmployeeAbbreviation string `json:"employeeAbbreviation"`
		Message              string `json:"message"`
	}
	o := output{
		Level:                level,
		EmployeeAbbreviation: employee,
		Message:              message,
	}

	payloadBytes, err := json.Marshal(o)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(payloadBytes)
	url := "http://admin-notif:8080/api/notify"
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return err
}
