package auth_test

import (
	"github.com/statistico/statistico-strategy/internal/trader/auth"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestAwsTokenParser_Parse(t *testing.T) {
	t.Run("parses token from JWT string", func(t *testing.T) {
		t.Helper()

		region := os.Getenv("AWS_REGION")
		userPoolID := os.Getenv("AWS_USER_POOL_ID")

		if region == "" || userPoolID == "" {
			t.Skip("AWS Region and User Pool ID required to run this test suite")
		}

		clock := MockClock{t: time.Unix(1617126949, 0)}

		parser := auth.NewAwsTokenParser(region, userPoolID, clock)

		token := "eyJraWQiOiJMWVpcL1RjK0V3S2xtQ2hcL1czcnRVRHA3bkR0Z3Rwbm04UHI3MUpYSHEzT0E9IiwiYWxnIjoiUlMyNTYifQ." +
			"eyJzdWIiOiI5MDRlMDNmYi1iYTAxLTQ1YjAtYWI2Ni0wMWE0ZDIyYTM2YjgiLCJldmVudF9pZCI6ImE4MmNkNjVkLTVmMzktNGVlMi1h" +
			"YmMyLTU3NTBlODM2MWE0NCIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iL" +
			"CJhdXRoX3RpbWUiOjE2MTY3MDgzMzUsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTIuYW1hem9uYXdzLmNvbVwvZ" +
			"XUtd2VzdC0yX0xvODZOVUtFUiIsImV4cCI6MTYxNzEyODY1NywiaWF0IjoxNjE3MTI1MDU3LCJqdGkiOiJjYWMyYzQ0OS05YzAxLTQxZ" +
			"jQtOTk4MS1kNWIwZWY1MjRiMmYiLCJjbGllbnRfaWQiOiIxcHQ2azNwMHRibjQ2amJiMGtlZmVibDJtdSIsInVzZXJuYW1lIjoiOTA0Z" +
			"TAzZmItYmEwMS00NWIwLWFiNjYtMDFhNGQyMmEzNmI4In0.Ev2dcbHtrigECvbqcOfn0ree0d3jrHEDQmI2Wd9uZ-FHeskDGXYEVpzvP" +
			"KQXx38-gR2iKfGtPXxLR-uEh86IV2gKmIrrN3-fSyzmLCyEUlVp-ntkTZVss2HA_PIwHUZCRXF54Y350KXBzoFvaex7VPfd7wRM6I5U2" +
			"NmBE_qlMG4IGdhN46Wh6A8WBMe2MVf9Txz-KT7YYOzq_G8Seyy-RblqKYgg31ERFGwCby7YsOTSt28iI7VMe3fY71evJuI-2XUDaE7Hu" +
			"--iUr5XxLWQ7YAv_n9rKVyBsQpGEShgrUdcOTcR6Hmdv0XTpt0N5eXtRCi0QVWu5ihi1dTr_Zk61Q"

		tn, err := parser.Parse(token)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, "904e03fb-ba01-45b0-ab66-01a4d22a36b8", tn.UserID)
	})

	t.Run("returns KeySetParseError error if unable to parse key set", func(t *testing.T) {
		t.Helper()

		clock := MockClock{t: time.Unix(1617126949, 0)}

		parser := auth.NewAwsTokenParser("eu-west-1", "invalid-id", clock)

		_, err := parser.Parse("token")

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		assert.Equal(t, "error parsing key set from JWT provider: failed to fetch remote JWK (status = 400)", err.Error())
	})

	t.Run("returns TokenValidationError if token provided is invalid", func(t *testing.T) {
		t.Helper()

		region := os.Getenv("AWS_REGION")
		userPoolID := os.Getenv("AWS_USER_POOL_ID")

		if region == "" || userPoolID == "" {
			t.Skip("AWS Region and User Pool ID required to run this test suite")
		}

		clock := MockClock{t: time.Unix(1617126949, 0)}

		parser := auth.NewAwsTokenParser(region, userPoolID, clock)

		token := "eyJraWQiOiJMWVpcL1RjK0V3S2xtQ2hcL1czcnRVRHA3bkR0Z3Rwbm04UHI3MUpYSHEzT0E9IiwiYWxnIjoiUlMyNTYifQ." +
			"eyJzdWIiOiI5MDRlMDNmYi1iYTAxLTQ1YjAtYWI2Ni0wMWE0ZDIyYTM2YjgiLCJldmVudF9pZCI6ImE4MmNkNjVkLTVmMzktNGVlMi1h" +
			"YmMyLTU3NTBlODM2MWE0NCIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iL" +
			"CJhdXRoX3RpbWUiOjE2MTY3MDgzMzUsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTIuYW1hem9uYXdzLmNvbVwvZ" +
			"XUtd2VzdC0yX0xvODZOVUtFUiIsImV4cCI6MTYxNzEyODY1NywiaWF0IjoxNjE3MTI1MDU3LCJqdGkiOiJjYWMyYzQ0OS05YzAxLTQxZ" +
			"jQtOTk4MS1kNWIwZWY1MjRiMmYiLCJjbGllbnRfaWQiOiIxcHQ2azNwMHRibjQ2amJiMGtlZmVibDJtdSIsInVzZXJuYW1lIjoiOTA0Z" +
			"TAzZmItYmEwMS00NWIwLWFiNjYtMDFhNGQyMmEzNmI4In0.Ev2dcbHtrigECvbqcOfn0ree0d3jrHEDQmI2Wd9uZ-FHeskDGXYEVpzvP" +
			"KQXx38-gR2iKfGtPXxLR-uEh86IV2gKmIrrN3-fSyzmLCyEUlVp-ntkTZVss2HA_PIwHUZCRXF54Y350KXBzoFvaex7VPfd7wRM6I5U2" +
			"NmBE_qlMG4IGdhN46Wh6A8WBMe2MVf9Txz-KT7YYOzq_G8Seyy-RblqKYgg31ERFGwCby7YsOTSt28iI7VMe3fY71evJuI-2XUDaE7Hu"

		_, err := parser.Parse(token)

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		assert.Equal(
			t,
			"error validating token: failed to find matching key for verification: failed to parse token data: failed to decode signature: failed to decode source: illegal base64 data at input byte 264",
			err.Error(),
		)
	})
}

type MockClock struct {
	t time.Time
}

func (m MockClock) Now() time.Time {
	return m.t
}
