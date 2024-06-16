package worker

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/dirsouza/go-driver/internal/bucket"
	"github.com/dirsouza/go-driver/internal/queue"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	queueCfg := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBITMQ_URL"),
		TopicName: os.Getenv("RABBITMQ_TOPIC_NAME"),
		Timeout:   time.Second * 30,
	}
	queueConn, err := queue.New(queue.RabbitMQ, queueCfg)
	if err != nil {
		panic(err)
	}

	chanel := make(chan queue.MessageDto)
	err = queueConn.Consume(chanel)
	if err != nil {
		panic(err)
	}

	bucketCfg := bucket.AwsConfig{
		Config: &aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		},
		BucketDown: os.Getenv("AWS_BUCKET_DOWN"),
		BucketUp:   os.Getenv("AWS_BUCKET_UP"),
	}

	bucketConn, err := bucket.New(bucket.Aws, bucketCfg)
	if err != nil {
		panic(err)
	}

	for msg := range chanel {
		src := fmt.Sprintf("%s/%s", msg.Path, msg.Filename)
		dest := fmt.Sprintf("%d_%s", msg.ID, msg.Filename)

		file, err := bucketConn.Download(src, dest)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		var buffer bytes.Buffer
		zipWriter := gzip.NewWriter(&buffer)
		_, err = zipWriter.Write(body)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		if err := zipWriter.Close(); err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		zipReader, err := gzip.NewReader(&buffer)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		err = bucketConn.Upload(zipReader, src)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		err = os.Remove(dest)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}
	}
}
