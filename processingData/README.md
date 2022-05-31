## Requeriments
- Python 3.9
- Docker
- SAM CLI
- AWS Credentials

## Setup
- Set your aws credentials at `~/.aws/credentials`
- Set the secrets at `template.yaml`
- Run the command to build and deploy
```bash
sam build
sam deploy
```

Note: Your AWS user should have permissions to access to:
- Lambda functions
- S3
- SES
- RDS
- IAM
- VPC