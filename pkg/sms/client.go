package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type SmsClient struct {
	c     *http.Client
	email string
	token string
}

type smsResult struct {
	Success bool `json:"success"`
	Data    struct {
		Number  []string
		Message string
	} `json:"data"`
}

func NewSmsClient(email, token string) *SmsClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	return &SmsClient{
		c:     client,
		token: token,
		email: email,
	}
}

func (s *SmsClient) SendSms(phone, msg string) error {
	var res smsResult
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:%s@gate.smsaero.ru/v2/sms/send?number=%s&text=%s&sign=SMS Aero", s.email, s.token, phone, msg), nil)
	resp, err := s.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return err
	}
	if !res.Success {
		if len(res.Data.Number) != 0 {
			return errors.New(fmt.Sprintf("phone %s", res.Data.Number[0]))
		}
		if res.Data.Message != "" {
			return errors.New(fmt.Sprintf("message %s", res.Data.Message))
		}
		return errors.New("Unknown error")
	}
	log.Debug().Interface("respBody", string(respBody)).Msg("SendSms")
	return nil
}
