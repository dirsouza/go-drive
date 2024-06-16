package bucket

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

type TypeBucket int

const (
	Aws TypeBucket = iota
)

type IBucket interface {
	Upload(io.Reader, string) error
	Download(string, string) (*os.File, error)
	Delete(string) error
}

type Bucket struct {
	provider IBucket
}

func New(typeBucket TypeBucket, cfg any) (bucket *Bucket, err error) {
	rt := reflect.TypeOf(cfg)

	switch typeBucket {
	case Aws:
		if rt.Name() != "AwsConfig" {
			return nil, fmt.Errorf("config need's to be of type AwsConfig")
		}

		bucket.provider = newAwsSession(cfg.(*AwsConfig))
	default:
		log.Fatal("invalid type bucket")
	}

	return
}

func (bucket *Bucket) Upload(file io.Reader, key string) error {
	return bucket.provider.Upload(file, key)
}

func (bucket *Bucket) Download(src, dest string) (file *os.File, err error) {
	return bucket.provider.Download(src, dest)
}

func (bucket *Bucket) Delete(key string) error {
	return bucket.provider.Delete(key)
}
