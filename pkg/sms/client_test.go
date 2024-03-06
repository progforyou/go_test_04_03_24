package sms

import (
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
	"testing/dating/api/pkg/app"
)

func TestClient(t *testing.T) {
	config := app.CreateConfig()
	client := NewSmsClient(config.Email, config.SmsToken)

	for _, test := range []struct {
		name        string
		client      *SmsClient
		phone       string
		message     string
		expectedErr error
	}{
		{
			name:        "valid",
			client:      client,
			phone:       "79281724695",
			message:     "test",
			expectedErr: nil,
		},
		{
			name:        "phone_error",
			client:      client,
			phone:       "777777",
			message:     "test",
			expectedErr: errors.New("phone incorrect"),
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.client.SendSms(test.phone, test.message)
			fmt.Println(err)
			assert.Equal(t, err, test.expectedErr)
		})
	}
}
