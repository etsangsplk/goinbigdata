package kinesis

import (
	"fmt"
	"github.com/upitau/goinbigdata/examples/cloudformation/tags"

	awskinesis "github.com/aws/aws-sdk-go/service/kinesis"
)

const KinesisCfnTemplate = `
Description: This is an AWS CloudFormation template for creating kinesis.

Parameters:

    RoleArn:
        Type: String

    InputDataStream:
        Type: String

    DataApplicationName:
        Type: String

    RetentionPeriodHours:
        Type: Integer

    ShardCount:
        Type: Integer

    KeyId:
        Type: String

    KCLDynamoDBPrimaryKey:
        Type: String

Resources:

    KinesisStream:
        Type: AWS::Kinesis::Stream
        Properties:
            Name: !Ref InputDataStream
            RetentionPeriodHours: !Ref RetentionPeriodHours
            ShardCount: !Ref ShardCount
            {{- if .ShouldEncrypt }}
            StreamEncryption:
            EncryptionType: {{ .EncryptionType}}
            KeyId: {{ .Kms.KeyId}}
            {{end -}}
            Tags:
            - Key: K8sNamespace
              Value: !Ref Namespace
            - Key: K8sCluster
              Value: !Ref Cluster
            - Key: SpecVersion
              Value: !Ref SpecVersion
            {{range .Tags.Tags -}}
            - Key: {{ .TagKey }}
              Value: {{ .TagValue }}
            {{ end }}
    KCLDynamoDBTable:
        Type: "AWS::DynamoDB::Table"
        Properties:
            AttributeDefinitions:
            - AttributeName: !Ref KCLDynamoDBPrimaryKey
              AttributeType: "S"
            KeySchema:
            - AttributeName: !Ref KCLDynamoDBPrimaryKey
              KeyType: "HASH"
            ProvisionedThroughput:
              ReadCapacityUnits: "10"
              WriteCapacityUnits: "10"
            TableName: !Ref DataApplicationName

Outputs:
    KinesisArn:
        Description: The Kinesis ARN created
        Value: !GetAtt KinesisStream.Arn

`

type AES struct {
	Algorithm string
}

func (a *AES) GetType() string {
	// return awskinesis.EncryptionTypeNone
	return "AES256"
}

type KMS struct {
	KeyId string
}

func (a *KMS) GetType() string {
	return awskinesis.EncryptionTypeKms
}

type EncryptionType string

type KinesisStreamDataObject struct {
	Aes            AES
	Kms            KMS
	Tags           tags.Tags
	EncryptionType EncryptionType
	ShouldEncrypt  bool
}

func CreateDataObject(specTags *map[string]string) KinesisStreamDataObject {
	var t tags.Tags
	if specTags != nil {
		for k, v := range *specTags {
			t.Tags = append(t.Tags, tags.Tag{TagKey: k, TagValue: v})
		}
	}
	a := KMS{
		KeyId: "arn:aws:kms:us-west-2:801351377084:key/85aa336d-be55-4c7f-b183-85c2f5c0e51b",
	}

	return KinesisStreamDataObject{
		Kms:            a,
		Tags:           t,
		EncryptionType: EncryptionType(a.GetType()),
		ShouldEncrypt:  true,
	}
}
