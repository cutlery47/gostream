package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cutlery47/gostream/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage interface {
	// stores file in the object storage
	Store(ctx context.Context, file File) (Location, error)
	// stores multiple files in the object storage
	StoreMultiple(ctx context.Context, files ...File) ([]Location, error)
	// retrieves object from the object storage
	Get(ctx context.Context, location Location) (io.ReadCloser, error)
	// deletes object from the object storage
	Delete(ctx context.Context, location Location) error
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

	ctx := context.Background()

	if err := s3.createBuckets(ctx, conf.VidBucket, conf.ChunkBucket, conf.ManBucket); err != nil {
		return nil, err
	}

	return s3, nil
}

func (s3 MinioS3) Store(ctx context.Context, file File) (Location, error) {
	bucket := s3.determineBucket(file.ObjectName)

	info, err := s3.cl.PutObject(ctx, bucket, file.ObjectName, file.Raw, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return Location{}, err
	}

	return Location{Bucket: info.Bucket, Object: info.Key}, nil
}

func (s3 MinioS3) StoreMultiple(ctx context.Context, files ...File) ([]Location, error) {
	var locs []Location

	// TODO: also try concurrent insertions

	for _, file := range files {
		loc, err := s3.Store(ctx, file)
		if err != nil {
			return locs, err
		}
		locs = append(locs, loc)
	}

	return locs, nil
}

func (s3 MinioS3) Get(ctx context.Context, loc Location) (file io.ReadCloser, err error) {
	return s3.cl.GetObject(ctx, loc.Bucket, loc.Object, minio.GetObjectOptions{})
}

func (s3 MinioS3) Delete(ctx context.Context, loc Location) error {
	return s3.cl.RemoveObject(ctx, loc.Bucket, loc.Object, minio.RemoveObjectOptions{})
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

func (s3 MinioS3) createBuckets(ctx context.Context, buckets ...string) error {
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
