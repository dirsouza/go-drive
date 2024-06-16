package bucket

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

type AwsConfig struct {
	Config     *aws.Config
	BucketDown string
	BucketUp   string
}

type awsSession struct {
	session    *session.Session
	bucketDown string
	bucketUp   string
}

func newAwsSession(cfg *AwsConfig) *awsSession {
	sessionMust := session.Must(session.NewSession(cfg.Config))

	return &awsSession{
		session:    sessionMust,
		bucketDown: cfg.BucketDown,
		bucketUp:   cfg.BucketUp,
	}
}

func (awsSess *awsSession) Upload(file io.Reader, key string) error {
	uploader := s3manager.NewUploader(awsSess.session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsSess.bucketUp),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

func (awsSess *awsSession) Download(src, dest string) (file *os.File, err error) {
	file, err = os.Create(dest)
	if err != nil {
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(awsSess.session)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(awsSess.bucketDown),
		Key:    aws.String(src),
	})

	return
}

func (awsSess *awsSession) Delete(key string) error {
	svc := s3.New(awsSess.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(awsSess.bucketDown),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsSess.bucketDown),
		Key:    aws.String(key),
	})
}
