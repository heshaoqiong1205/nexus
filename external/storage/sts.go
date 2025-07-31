package external_storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"

	ossApi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ossStsClient "github.com/alibabacloud-go/sts-20150401/v2/client"
	ossService "github.com/alibabacloud-go/tea-utils/v2/service"
	ossTea "github.com/alibabacloud-go/tea/tea"

	cosSts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

const (
	AccessKeyID     = "YourAccessKeyID"
	SecretAccessKey = "YourSecretAccessKey"
	SessionToken    = "YourSessionToken"
	RoleArn         = "arn:aws:iam::123456789012:role/YourRoleName"
	CosAppId        = "YourCosAppId"
)

type Action string

const (
	GetObject    Action = "GetObject"
	PutObject    Action = "PutObject"
	DeleteObject Action = "DeleteObject"
	ListBucket   Action = "ListBucket"
	HeadObject   Action = "HeadObject"
	AllActions   Action = "*"
)

type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)

type Credentials struct {
	AccessKeyId     *string
	SecretAccessKey *string
	SessionToken    *string
	Expiration      *time.Time
}

type Resource struct {
	Bucket string
	Path   string
}

type Statement struct {
	Effect   Effect
	Actions  []Action
	Resource Resource
}

type AwsStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

type AwsPolicy struct {
	Version    string         `json:"Version"`
	Statements []AwsStatement `json:"Statement"`
}

type OssStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

type OssPolicy struct {
	Version    string         `json:"Version"`
	Statements []OssStatement `json:"Statement"`
}

type IStsService interface {
	GenerateCredentials(provider string, region string, statement []Statement, duration int32) (Credentials, error)
}

type StsService struct{}

func NewStsService() *StsService {
	return &StsService{}
}

func (ss *StsService) GenerateCredentials(provider string, region string, statement []Statement, duration int32) (Credentials, error) {
	switch provider {
	case "aws":
		return generateAwsCredentials(region, statement, duration)
	case "oss":
		return generateOssCredentials(region, statement, duration)
	case "cos":
		return generateCosCredentials(region, statement, duration)
	case "mock":
		return generateMockCredentials(region, statement, duration)
	default:
		return Credentials{}, errors.New("unsupported provider")
	}
}

func generateAwsCredentials(region string, statement []Statement, duration int32) (Credentials, error) {
	// This function would typically generate AWS credentials for the specified bucket.
	// For simplicity, let's return a dummy set of credentials.
	// In a real application, this would involve using AWS SDK to generate temporary credentials.

	staticCredentialsProvider := credentials.NewStaticCredentialsProvider(AccessKeyID, SecretAccessKey, SessionToken)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(staticCredentialsProvider))
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	policy := AwsPolicy{
		Version:    "2012-10-17",
		Statements: newAwsStatement(statement),
	}
	value, err := json.Marshal(policy)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to marshal policy: %w", err)
	}
	stsClient := sts.NewFromConfig(cfg)
	input := sts.AssumeRoleInput{
		RoleArn:         aws.String(RoleArn),
		RoleSessionName: aws.String(uuid.New().String()),
		Policy:          aws.String(string(value)),
		DurationSeconds: aws.Int32(duration),
	}
	assumeRoleOutput, err := stsClient.AssumeRole(context.TODO(), &input)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to assume role: %w", err)
	}
	return Credentials{
		AccessKeyId:     assumeRoleOutput.Credentials.AccessKeyId,
		SecretAccessKey: assumeRoleOutput.Credentials.SecretAccessKey,
		SessionToken:    assumeRoleOutput.Credentials.SessionToken,
		Expiration:      assumeRoleOutput.Credentials.Expiration,
	}, nil
}

/*
see https://help.aliyun.com/zh/oss/developer-reference/use-temporary-access-credentials-provided-by-sts-to-access-oss
*/
func generateOssCredentials(region string, statement []Statement, duration int32) (Credentials, error) {
	// This function would typically generate OSS credentials for the specified bucket.
	// For simplicity, let's return a dummy set of credentials.
	// In a real application, this would involve using OSS SDK to generate temporary credentials.

	cfg := &ossApi.Config{
		AccessKeyId:     ossTea.String(AccessKeyID),
		AccessKeySecret: ossTea.String(SecretAccessKey),
		Endpoint:        ossTea.String(fmt.Sprintf("sts.%s.aliyuncs.com", region)),
	}

	client, err := ossStsClient.NewClient(cfg)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to create OSS client: %w", err)
	}

	policy := OssPolicy{
		Version:    "1",
		Statements: newOssStatement(statement),
	}

	value, err := json.Marshal(policy)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to marshal policy: %w", err)
	}

	input := &ossStsClient.AssumeRoleRequest{
		RoleArn:         ossTea.String(RoleArn),
		RoleSessionName: ossTea.String(uuid.New().String()),
		Policy:          ossTea.String(string(value)),
		DurationSeconds: ossTea.Int64(int64(duration)),
	}

	output, err := client.AssumeRoleWithOptions(input, &ossService.RuntimeOptions{})

	if err != nil {
		return Credentials{}, fmt.Errorf("failed to assume OSS role: %w", err)
	}
	expiration, err := time.Parse(time.RFC3339, ossTea.StringValue(output.Body.Credentials.Expiration))
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to parse expiration time: %w", err)
	}
	// Here you would call the OSS SDK to generate the temporary credentials
	// For now, we return dummy credentials
	return Credentials{
		AccessKeyId:     output.Body.Credentials.AccessKeyId,
		SecretAccessKey: output.Body.Credentials.AccessKeySecret,
		SessionToken:    output.Body.Credentials.SecurityToken,
		Expiration:      &expiration,
	}, nil
}

func generateCosCredentials(region string, statement []Statement, duration int32) (Credentials, error) {
	// This function would typically generate COS credentials for the specified bucket.
	// For simplicity, let's return a dummy set of credentials.
	// In a real application, this would involve using COS SDK to generate temporary credentials.

	cosClient := cosSts.NewClient(AccessKeyID, SecretAccessKey, nil)

	opt := &cosSts.CredentialOptions{
		DurationSeconds: int64(duration),
		Region:          region,
		Policy: &cosSts.CredentialPolicy{
			Statement: newCosStatement(CosAppId, region, statement),
		},
	}
	output, err := cosClient.GetCredential(opt)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to get COS session token: %w", err)
	}

	expiration, err := time.Parse(time.RFC3339, output.Expiration)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to parse expiration time: %w", err)
	}
	return Credentials{
		AccessKeyId:     &output.Credentials.TmpSecretID,
		SecretAccessKey: &output.Credentials.TmpSecretKey,
		SessionToken:    &output.Credentials.SessionToken,
		Expiration:      &expiration,
	}, nil
}

func generateMockCredentials(_ string, _ []Statement, duration int32) (Credentials, error) {
	return Credentials{
		AccessKeyId:     aws.String("mock-access-key-id"),
		SecretAccessKey: aws.String("mock-secret-access-key"),
		SessionToken:    aws.String("mock-session-token"),
		Expiration:      aws.Time(time.Now().Add(time.Duration(duration) * time.Second)),
	}, nil
}

func newAwsStatement(statement []Statement) []AwsStatement {
	awsStatements := make([]AwsStatement, len(statement))
	for i, stmt := range statement {
		awsStatements[i] = AwsStatement{
			Effect:   string(stmt.Effect),
			Action:   newAwsActions(stmt.Actions),
			Resource: fmt.Sprintf("arn:aws:s3:::%s/%s", stmt.Resource.Bucket, stmt.Resource.Path),
		}
	}
	return awsStatements
}

func newOssStatement(statement []Statement) []OssStatement {
	ossStatement := make([]OssStatement, len(statement))
	for i, stmt := range statement {
		ossStatement[i] = OssStatement{
			Effect:   string(stmt.Effect),
			Action:   newOssActions(stmt.Actions),
			Resource: []string{fmt.Sprintf("acs:oss:*:*:%s/%s", stmt.Resource.Bucket, stmt.Resource.Path)},
		}
	}
	return ossStatement
}

func newCosStatement(appid, region string, statement []Statement) []cosSts.CredentialPolicyStatement {
	cosStatements := make([]cosSts.CredentialPolicyStatement, len(statement))
	for i, stmt := range statement {
		cosStatements[i] = cosSts.CredentialPolicyStatement{
			Effect:   string(stmt.Effect),
			Action:   newCosActions(stmt.Actions),
			Resource: []string{fmt.Sprintf("qcs::cos:%s:uuid/%s:%s/%s", region, appid, stmt.Resource.Bucket, stmt.Resource.Path)},
		}
	}
	return cosStatements
}

func newAwsActions(actions []Action) []string {
	awsActions := make([]string, len(actions))
	for i, action := range actions {
		switch action {
		case GetObject:
			awsActions[i] = "s3:GetObject"
		case PutObject:
			awsActions[i] = "s3:PutObject"
		case DeleteObject:
			awsActions[i] = "s3:DeleteObject"
		case ListBucket:
			awsActions[i] = "s3:ListBucket"
		case HeadObject:
			awsActions[i] = "s3:HeadObject"
		case AllActions:
			awsActions[i] = "s3:*"
		default:
			awsActions[i] = string(action)
		}
	}
	return awsActions
}

func newOssActions(actions []Action) []string {
	ossActions := make([]string, len(actions))
	for i, action := range actions {
		switch action {
		case GetObject:
			ossActions[i] = "oss:GetObject"
		case PutObject:
			ossActions[i] = "oss:PutObject"
		case DeleteObject:
			ossActions[i] = "oss:DeleteObject"
		case ListBucket:
			ossActions[i] = "oss:ListBucket"
		case HeadObject:
			ossActions[i] = "oss:HeadObject"
		case AllActions:
			ossActions[i] = "oss:*"
		default:
			ossActions[i] = string(action)
		}
	}
	return ossActions
}

func newCosActions(actions []Action) []string {
	cosActions := make([]string, len(actions))
	for i, action := range actions {
		switch action {
		case GetObject:
			cosActions[i] = "cos:GetObject"
		case PutObject:
			cosActions[i] = "cos:PutObject"
		case DeleteObject:
			cosActions[i] = "cos:DeleteObject"
		case ListBucket:
			cosActions[i] = "cos:ListBucket"
		case HeadObject:
			cosActions[i] = "cos:HeadObject"
		case AllActions:
			cosActions[i] = "cos:*"
		default:
			cosActions[i] = string(action)
		}
	}
	return cosActions
}
