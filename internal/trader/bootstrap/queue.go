package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/statistico/statistico-trader/internal/trader/queue"
	qa "github.com/statistico/statistico-trader/internal/trader/queue/aws"
	"github.com/statistico/statistico-trader/internal/trader/queue/log"
)

func (c Container) Queue() queue.MarketQueue {
	if c.Config.QueueDriver == "aws" {
		key := c.Config.AWS.Key
		secret := c.Config.AWS.Secret

		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, ""),
			Region:      aws.String(c.Config.AWS.Region),
		})

		if err != nil {
			panic(err)
		}

		return qa.NewMarketQueue(
			sqs.New(sess),
			c.Logger,
			c.Config.AWS.QueueUrl,
			10,
		)
	}

	if c.Config.QueueDriver == "log" {
		return log.NewMarketQueue(c.Logger)
	}

	panic("Queue driver provided is not supported")
}
