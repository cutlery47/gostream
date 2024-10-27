package obj

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cutlery47/gostream/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage interface {
	Store(file InFile) (location string, err error)
	StoreMultiple(files ...InFile) (locations []string, err error)
	Get(filename string) (file *InFile, err error)
	Delete(filename string) (file *InFile, err error)
}

type MinioS3 struct {
	cl *minio.Client

	conf config.S3Config
}

func NewS3(conf config.S3Config) (*MinioS3, error) {
	url := fmt.Sprintf("%v:%v", conf.Host, conf.Port)

	credentials := credentials.NewEnvMinio()

	minioClient, err := minio.New(url, &minio.Options{Creds: credentials})
	if err != nil {
		return nil, err
	}

	s3 := &MinioS3{
		cl:   minioClient,
		conf: conf,
	}

	if err := s3.createBuckets(conf.VidBucket, conf.ChunkBucket, conf.ManBucket); err != nil {
		return nil, err
	}

	return s3, nil
}

func (s3 MinioS3) Store(file InFile) (key string, err error) {
	ctx := context.Background()

	opts := minio.PutObjectOptions{}

	bucket := s3.determineBucket(file.Name)
	info, err := s3.cl.PutObject(ctx, bucket, file.Name, file.Raw, int64(file.Size), opts)
	if err != nil {
		return key, err
	}

	return fmt.Sprintf("%v/%v", info.Bucket, info.Key), err
}

func (s3 MinioS3) StoreMultiple(files ...InFile) (keys []string, err error) {
	for _, file := range files {
		key, err := s3.Store(file)
		if err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (s3 MinioS3) Get(filename string) (file *InFile, err error) {
	return nil, fmt.Errorf("sdafasdf")
}

func (s3 MinioS3) Delete(filename string) (file *InFile, err error) {
	return nil, fmt.Errorf("xyu yxu")
}

func (s3 MinioS3) determineBucket(filename string) (bucket string) {
	if strings.HasSuffix(filename, ".mp4") {
		return s3.conf.VidBucket
	}

	if strings.HasSuffix(filename, ".m3u8") {
		return s3.conf.ManBucket
	}

	if strings.HasSuffix(filename, ".ts") {
		return s3.conf.ChunkBucket
	}

	panic("couldn't determine bucket")
}

func (s3 MinioS3) createBuckets(buckets ...string) error {
	ctx := context.Background()

	for _, bucket := range buckets {
		// create a bucket by bucket name
		err := s3.cl.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			// check to see if we already own this bucket
			exists, errExists := s3.cl.BucketExists(ctx, bucket)
			if errExists == nil && exists {
				log.Printf("Bucket %s already exitsts\n", bucket)
			} else {
				return err
			}
		} else {
			log.Printf("Successfully created bucket %s\n", bucket)
		}
	}

	return nil
}
