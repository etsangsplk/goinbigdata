package s3

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
            {{range .Tags -}}
            - Key: {{ .TagKey }}
              Value: {{ .TagValue }}
            {{end}}
Outputs:
    S3Bucket:
        Description: The name of the S3 bucket
        Value: !GetAtt S3Bucket.Arn
`

type EncryptionType struct {
	Type string
}

type AES struct {
	EncryptionType
	Algorithm string
}

type Kms struct {
	EncryptionType
	KeyId string
}

type Encryption struct {
	Type EncryptionType
}

type S3CfnDataObject struct {
	Encryption Encryption
}
