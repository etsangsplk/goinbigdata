package kinesis

const kinesisCfnTemplate = `
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
            {{- if .Encryption.ShouldEncrypt }}
            StreamEncryption:
            EncryptionType: !Ref EncryptionType
            KeyId: !Ref KeyId
            {{end -}}
            Tags:
            - Key: K8sNamespace
              Value: !Ref Namespace
            - Key: K8sCluster
              Value: !Ref Cluster
            - Key: SpecVersion
              Value: !Ref SpecVersion
            {{range .Tags -}}
            -
            Key: {{ .TagKey }}
            Value: {{ .TagValue }}
            {{end -}}

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
