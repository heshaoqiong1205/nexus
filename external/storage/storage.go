package external_storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	ossCredentials "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type IStorageService interface {
	GeneratePutPresignedURL(provider string, region, bucket, objectKey string, expiration time.Duration) (string, error)
	GenerateGetPresignedURL(provider string, region, bucket, objectKey string, expiration time.Duration) (string, error)
	ListObjects(provider, region, bucket string, prefix string, maxKeys int32) ([]string, error)
}

type StorageService struct {
}

func (ss *StorageService) GeneratePutPresignedURL(provider string, region, bucket, objectKey string, expiration time.Duration) (string, error) {
	// This function should generate a presigned URL based on the storage provider.
	// The implementation will vary depending on whether it's AWS S3, Alibaba Cloud OSS, or Tencent COS.
	switch provider {
	case "aws":
		return generateAwsPutPresignedURL(region, bucket, objectKey, expiration)
	case "oss":
		return generateOssPutPresignedURL(region, bucket, objectKey, expiration)
	case "cos":
		return generateCosPutPresignedURL(region, bucket, objectKey, expiration)
	case "mock":
		return generateMockPutPresignedURL(region, bucket, objectKey, expiration)
	default:
		return "", errors.New("unsupported provider")
	}
}

func (ss *StorageService) GenerateGetPresignedURL(provider string, region, bucket, objectKey string, expiration time.Duration) (string, error) {
	// This function should generate a presigned URL for getting an object from the storage provider.
	// The implementation will vary depending on whether it's AWS S3, Alibaba Cloud OSS, or Tencent COS.
	switch provider {
	case "aws":
		return generateAwsGetPresignedURL(region, bucket, objectKey, expiration)
	case "oss":
		return generateOssGetPresignedURL(region, bucket, objectKey, expiration)
	case "cos":
		return generateCosGetPresignedURL(region, bucket, objectKey, expiration)
	case "mock":
		return generateMockGetPresignedURL(region, bucket, objectKey, expiration)
	default:
		return "", errors.New("unsupported provider")
	}
}

func (ss *StorageService) ListObjects(provider, region, bucket string, prefix string, maxKeys int32) ([]string, error) {
	// This function should list objects in the specified bucket and prefix.
	// The implementation will vary depending on whether it's AWS S3, Alibaba Cloud OSS, or Tencent COS.
	switch provider {
	case "aws":
		return listAwsObjects(region, bucket, prefix, maxKeys)
	case "oss":
		return listOssObjects(region, bucket, prefix, maxKeys)
	case "cos":
		return listCosObjects(region, bucket, prefix, maxKeys)
	case "mock":
		return listMockObjects(region, bucket, prefix, maxKeys)
	default:
		return nil, errors.New("unsupported provider")
	}
}

func generateAwsPutPresignedURL(region, bucket, objectKey string, duration time.Duration) (string, error) {
	// This function should generate a presigned URL for AWS S3.
	presignClient, err := newAwsPresignClient(region)
	if err != nil {
		return "", err
	}

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &objectKey,
	}
	presignOutput, err := presignClient.PresignPutObject(context.TODO(), input, func(po *s3.PresignOptions) {
		po.Expires = duration
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}
	return presignOutput.URL, nil
}

func generateOssPutPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {
	// This function should generate a presigned URL for Alibaba Cloud OSS.
	// x := ossCredentials.NewStaticCredentialsProvider(AccessKeyID, SecretAccessKey, SessionToken)
	client := newOssClient(region)

	input := &oss.PutObjectRequest{
		Bucket: &bucket,
		Key:    &objectKey,
	}
	result, err := client.Presign(context.TODO(), &input, oss.PresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}
	return result.URL, nil
}

func generateCosPutPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {
	// This function should generate a presigned URL for Tencent COS.
	client, err := newCosClient(region, bucket)
	if err != nil {
		return "", err
	}

	presignedURL, err := client.Object.GetPresignedURL(context.TODO(), http.MethodPut, objectKey, client.GetCredential().SecretID, client.GetCredential().SecretKey, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}
	return presignedURL.String(), nil
}

func generateMockPutPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {
	return fmt.Sprintf("https://%s.mockstorage.com/%s/%s?expires=%d", region, bucket, objectKey, time.Now().Add(expiration).Unix()), nil
}

func generateAwsGetPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {
	// This function should generate a presigned URL for getting an object from AWS S3.
	client, err := newAwsPresignClient(region)
	if err != nil {
		return "", err
	}

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objectKey,
	}
	presignOutput, err := client.PresignGetObject(context.TODO(), input, func(po *s3.PresignOptions) {
		po.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign get object: %w", err)
	}
	return presignOutput.URL, nil
}

func generateOssGetPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {
	// This function should generate a presigned URL for getting an object from Alibaba Cloud OSS.
	client := newOssClient(region)
	input := &oss.GetObjectRequest{
		Bucket: &bucket,
		Key:    &objectKey,
	}
	result, err := client.Presign(context.TODO(), &input, oss.PresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to presign get object: %w", err)
	}
	return result.URL, nil
}

func generateCosGetPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {

	// This function should generate a presigned URL for getting an object from Tencent COS.
	client, err := newCosClient(region, bucket)
	if err != nil {
		return "", err
	}

	presignedURL, err := client.Object.GetPresignedURL(context.TODO(), http.MethodGet, objectKey, client.GetCredential().SecretID, client.GetCredential().SecretKey, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}
	return presignedURL.String(), nil
}

func generateMockGetPresignedURL(region, bucket, objectKey string, expiration time.Duration) (string, error) {
	return fmt.Sprintf("https://%s.mockstorage.com/%s/%s?expires=%d", region, bucket, objectKey, time.Now().Add(expiration).Unix()), nil
}

func listAwsObjects(region, bucket, prefix string, maxKeys int32) ([]string, error) {
	// This function should list objects in an AWS S3 bucket.
	client, err := newAwsClient(region)
	if err != nil {
		return nil, err
	}
	input := &s3.ListObjectsV2Input{
		Bucket:  &bucket,
		Prefix:  &prefix,
		MaxKeys: &maxKeys,
	}

	result, err := client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to list items in bucket %q, %v", bucket, err)
	}

	var objects []string
	for _, item := range result.Contents {
		objects = append(objects, *item.Key)
	}
	return objects, nil
}

func listOssObjects(region, bucket, prefix string, maxKeys int32) ([]string, error) {
	// This function should list objects in an Alibaba Cloud OSS bucket.
	client := newOssClient(region)

	input := &oss.ListObjectsV2Request{
		Bucket:  &bucket,
		Prefix:  &prefix,
		MaxKeys: maxKeys,
	}
	paginator := client.NewListObjectsV2Paginator(input)

	var objects []string
	for paginator.HasNext() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in bucket %q: %w", bucket, err)
		}
		for _, object := range page.Contents {
			objects = append(objects, *object.Key)
		}
	}
	return objects, nil
}

func listCosObjects(region, bucket string, prefix string, maxKeys int32) ([]string, error) {
	// This function should list objects in a Tencent COS bucket.
	client, err := newCosClient(region, bucket)
	if err != nil {
		return nil, err
	}

	opt := &cos.BucketGetOptions{
		Prefix:  prefix,
		MaxKeys: int(maxKeys),
	}
	result, _, err := client.Bucket.Get(context.TODO(), opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %q: %w", bucket, err)
	}

	var objects []string
	for _, item := range result.Contents {
		objects = append(objects, item.Key)
	}
	return objects, nil
}

func listMockObjects(_, _, _ string, maxKeys int32) ([]string, error) {
	var objects []string
	for i := 0; i < int(maxKeys); i++ {
		objects = append(objects, fmt.Sprintf("mock-object-%d", i))
	}
	return objects, nil
}

func newAwsPresignClient(region string) (*s3.PresignClient, error) {
	staticCredentialsProvider := credentials.NewStaticCredentialsProvider(AccessKeyID, SecretAccessKey, SessionToken)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(staticCredentialsProvider))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	return s3.NewPresignClient(client), nil
}

func newOssClient(region string) *oss.Client {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(ossCredentials.NewStaticCredentialsProvider(AccessKeyID, SecretAccessKey, SessionToken)).
		WithRegion(region)
	client := oss.NewClient(cfg)
	return client
}

func newCosClient(region, bucket string) (*cos.Client, error) {
	// This function should create a new Tencent COS presign client.
	// The implementation will depend on the Tencent COS SDK.
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	if err != nil {
		return nil, err
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     AccessKeyID,
			SecretKey:    SecretAccessKey,
			SessionToken: SessionToken,
		},
	})

	return client, nil
}

func newAwsClient(region string) (*s3.Client, error) {
	// This function should create a new AWS S3 client.
	staticCredentialsProvider := credentials.NewStaticCredentialsProvider(AccessKeyID, SecretAccessKey, SessionToken)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(staticCredentialsProvider))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return s3.NewFromConfig(cfg), nil
}
