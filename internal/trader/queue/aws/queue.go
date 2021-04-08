package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-trader/internal/trader/queue"
)

type Message struct {
	Type           string `json:"type"`
	MessageID      string `json:"messageId"`
	TopicArn       string `json:"topicArn"`
	Message        string `json:"message"`
	Timestamp      string `json:"timestamp"`
	Signature      string `json:"signature"`
	SigningCertURL string `json:"signingCertUrl"`
	UnsubscribeURL string `json:"unsubscribeUrl"`
}

type marketQueue struct {
	client   sqsiface.SQSAPI
	logger   *logrus.Logger
	queueUrl string
	timeout  int64
}

func (q *marketQueue) ReceiveMarkets() <-chan *queue.EventMarket {
	ch := make(chan *queue.EventMarket, 100)

	go q.receiveMessages(ch)

	return ch
}

func (q *marketQueue) receiveMessages(ch chan<- *queue.EventMarket) {
	defer close(ch)

	input := &sqs.ReceiveMessageInput{
		QueueUrl: &q.queueUrl,
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		WaitTimeSeconds: &q.timeout,
	}

	result, err := q.client.ReceiveMessage(input)

	if err != nil {
		q.logger.Errorf("Unable to receive messages from queue %q, %v.", q.queueUrl, err)
		return
	}

	for _, message := range result.Messages {
		q.parseMessage(message, ch)
	}
}

func (q *marketQueue) parseMessage(ms *sqs.Message, ch chan<- *queue.EventMarket) {
	var message Message
	err := json.Unmarshal([]byte(*ms.Body), &message)

	if err != nil {
		q.logger.Errorf("Unable to marshal message into message struct, %v.", err)
		return
	}

	var mk *queue.EventMarket
	err = json.Unmarshal([]byte(message.Message), &mk)

	if err != nil {
		q.logger.Errorf("Unable to marshal message into market struct, %v.", err)
		return
	}

	ch <- mk

	q.deleteMessage(ms.ReceiptHandle)
}

func (q *marketQueue) deleteMessage(handle *string) {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueUrl,
		ReceiptHandle: handle,
	}

	_, err := q.client.DeleteMessage(input)

	if err != nil {
		q.logger.Errorf("Error deleting message from queue %q", err)
	}
}

func NewMarketQueue(c sqsiface.SQSAPI, l *logrus.Logger, queue string, timeout int64) queue.MarketQueue {
	return &marketQueue{
		client:   c,
		logger:   l,
		queueUrl: queue,
		timeout:  timeout,
	}
}
