package s3

import (
	"fmt"
	"github.com/upitau/goinbigdata/examples/cloudformation/tags"
)

const S3CfnTemplate = `
Description: This is an AWS CloudFormation template for creating an S3 bucket.

Parameters:
    S3BucketName:
        Type: String

    Namespace:
        Type: String

    Cluster:
        Type: String

    SpecVersion:
        Type: String

    VersioningConfiguration:
        Type: String

Resources:
    S3Bucket:
        Type: AWS::S3::Bucket
        Properties:
            BucketName: !Ref S3BucketName
            VersioningConfiguration:
                Status: !Ref VersioningConfiguration
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
            {{end}}
            {{- if .ShouldEncrypt }}
            BucketEncryption:
                ServerSideEncryptionConfiguration:
                - ServerSideEncryptionByDefault:
                     SSEAlgorithm: {{ .Type}}
            {{- else }}
            {{- end }}

    {{- if .ShouldEncrypt }}
    S3BucketPolicy:
        Type: AWS::S3::BucketPolicy
        Properties:
            Bucket: !Ref S3BucketName
            PolicyDocument:
                Version: '2012-10-17'
                Statement:
                    Effect: Deny
                    Principal:
                         AWS: !Ref RoleArn
                    Action:
                    - s3:PutObject
                    Resource: !Join ["", ["arn:aws:s3:::", !Ref S3BucketName, "/*"]]
                    Condition:
                        StringNotEquals:
                            s3:x-amz-server-side-encryption: {{ .Type}}
    {{- else }}
    {{- end }}

Outputs:
    S3Bucket:
        Description: The name of the S3 bucket
        Value: !GetAtt S3Bucket.Arn
`

type AES struct {
	Algorithm string
}

func (a *AES) GetType() string {
	return "AES256"
}

type KMS struct {
	KeyId string
}

func (a *KMS) GetType() string {
	return "aws.kms"
}

type EncryptionType string

type S3CfnDataObject struct {
	//Aes           AES
	//Kms           KMS
	Tags          tags.Tags
	Type          EncryptionType
	ShouldEncrypt bool
}

func CreateDataObject(specTags *map[string]string) S3CfnDataObject {
	var t tags.Tags
	if specTags != nil {
		for k, v := range *specTags {
			t.Tags = append(t.Tags, tags.Tag{TagKey: k, TagValue: v})
		}
		fmt.Printf("tags : %v", t)
	}
	a := AES{
		Algorithm: "AES256",
	}
	o := S3CfnDataObject{
		Tags:          t,
		Type:          EncryptionType(a.GetType()),
		ShouldEncrypt: true,
	}
	fmt.Printf("S3CfnDataObject: %v \n", o)
	return o
}
