package aws_test

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	saws "github.com/statistico/statistico-trader/internal/trader/queue/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestQueue_ReceiveMarkets(t *testing.T) {
	t.Run("calls client and pushes messages into channel provided", func(t *testing.T) {
		t.Helper()

		client := new(saws.MockSqsClient)
		logger, _ := test.NewNullLogger()

		r := saws.NewMarketQueue(client, logger, "messages", 3600)

		input := mock.MatchedBy(func(i *sqs.ReceiveMessageInput) bool {
			assert.Equal(t, "messages", *i.QueueUrl)
			assert.Equal(t, int64(3600), *i.WaitTimeSeconds)
			return true
		})

		messages := []*sqs.Message{
			{
				ReceiptHandle: aws.String("1234"),
				Body:          &messageBody,
			},
		}

		deleteInput := mock.MatchedBy(func(i *sqs.DeleteMessageInput) bool {
			assert.Equal(t, "1234", *i.ReceiptHandle)
			return true
		})

		client.On("ReceiveMessage", input).Return(&sqs.ReceiveMessageOutput{Messages: messages}, nil)
		client.On("DeleteMessage", deleteInput).Return(&sqs.DeleteMessageOutput{}, nil)

		date, _ := time.Parse(time.RFC3339, "2019-01-14T11:00:00Z")

		mk := &queue.EventMarket{
			ID:       "1.2818721",
			EventID:  148192,
			Name:     "OVER_UNDER_25",
			EventDate: date,
			Exchange: "betfair",
			Runners: []*queue.Runner{
				{
					ID:   472671,
					Name: "Over 2.5 Goals",
					BackPrices: []queue.PriceSize{
						{
							Price: 1.95,
							Size:  1461,
						},
					},
					LayPrices: []queue.PriceSize{
						{
							Price: 1.95,
							Size:  1461,
						},
					},
				},
			},
			Timestamp: 1583971200,
		}

		ch := r.ReceiveMarkets()

		mt := <-ch

		assert.Equal(t, mk, mt)
		client.AssertExpectations(t)
	})

	t.Run("logs and returns error if error returned by SQS client", func(t *testing.T) {
		t.Helper()

		client := new(saws.MockSqsClient)
		logger, hook := test.NewNullLogger()

		r := saws.NewMarketQueue(client, logger, "messages", 3600)

		input := mock.MatchedBy(func(i *sqs.ReceiveMessageInput) bool {
			assert.Equal(t, "messages", *i.QueueUrl)
			assert.Equal(t, int64(3600), *i.WaitTimeSeconds)
			return true
		})

		e := errors.New("error happened")

		client.On("ReceiveMessage", input).Return(&sqs.ReceiveMessageOutput{}, e)

		ch := r.ReceiveMarkets()

		<-ch

		assert.Equal(t, 0, len(ch))
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Equal(t, "Unable to receive messages from queue \"messages\", error happened.", hook.LastEntry().Message)

		client.AssertExpectations(t)
	})

	t.Run("logs error if unable to parse message body in market struct", func(t *testing.T) {
		t.Helper()

		client := new(saws.MockSqsClient)
		logger, hook := test.NewNullLogger()

		r := saws.NewMarketQueue(client, logger, "messages", 3600)

		input := mock.MatchedBy(func(i *sqs.ReceiveMessageInput) bool {
			assert.Equal(t, "messages", *i.QueueUrl)
			assert.Equal(t, int64(3600), *i.WaitTimeSeconds)
			return true
		})

		body := "invalid body"

		messages := []*sqs.Message{
			{
				Body: &body,
			},
		}

		client.On("ReceiveMessage", input).Return(&sqs.ReceiveMessageOutput{Messages: messages}, nil)

		ch := r.ReceiveMarkets()

		<-ch

		assert.Equal(t, 0, len(ch))
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Equal(t, "Unable to marshal message into message struct, invalid character 'i' looking for beginning of value.", hook.LastEntry().Message)

		client.AssertExpectations(t)
	})
}

var messageBody = `
	{
	  "Type": "Notification",
	  "MessageId": "72b286fb-a288-5b04-9093-dee1c8e08a85",
	  "TopicArn": "arn:aws:",
	  "Message": "{\"id\":\"1.2818721\",\"eventId\":148192,\"name\":\"OVER_UNDER_25\",\"date\":\"2019-01-14T11:00:00Z\",\"exchange\":\"betfair\",\"runners\":[{\"id\":472671,\"name\":\"Over 2.5 Goals\",\"backPrices\":[{\"price\":1.95,\"size\":1461}],\"layPrices\":[{\"price\":1.95,\"size\":1461}]}],\"timestamp\":1583971200}",
	  "Timestamp": "2020-11-02T20:12:24.030Z",
	  "SignatureVersion": "1",
	  "Signature": "aMVOnhHOyvVg4JhJ1TfopQQ55Ow5EbqzA6A/Cbhxl+ZOhI9fTEogukCQAG4lMBReh0Xbtx2BIJx+j+pDgKW3FPEuZxP/CeKdLQU+KAP1J86Nlja1cAeNMk05tJE6P4IwR07P6+0hIsZmEE9bFfwV5zw5cip7TnbpD/o9QyPnEv8Dt16RDprQfkuuJa+XAUvpFOgX6l1SQRnf3AwmZeV9H6mWPLFSyrc2RKkRzlOhbNXt31qul7+fT4R23p90TB42UtGXsf73l40Pz6s4ibb9QzMhl0kjHW7qwsH0iRMYJFtznoX4YP/X4InVzSYl7vv201ih3Wiixu0gbNByM8OBFg==",
	  "SigningCertURL": "https://sns.eu-west-2.amazonaws.com",
	  "UnsubscribeURL": "https://sns.eu-west-2.amazonaws.com"
	}
`
