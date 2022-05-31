import os
import urllib

import boto3
import logging


logger = logging.getLogger()
logger.setLevel(logging.INFO)

s3_client = boto3.client("s3")
ses_client = boto3.client("ses")

logger.info("SMTP connection successfully")

EMAIL_BUCKET_NAME = os.getenv("EMAIL_BUCKET_NAME", "")
EMAIL_SOURCE = os.getenv("EMAIL_SOURCE", "")


def send_report(email, message) -> bool:
    try:
        ses_client.send_email(
            Source=EMAIL_SOURCE,
            Destination={
                'ToAddresses': [
                    email,
                ],
            },
            Message={
                'Subject': {
                    'Data': 'Transactions Report',
                    'Charset': 'utf-8'
                },
                'Body': {
                    'Html': {
                        'Data': message,
                        'Charset': 'utf-8'
                    }
                }
            }
        )
        return True
    except Exception as error:
        logger.error(error)
        return False


def lambda_handler(event, context):

    if EMAIL_BUCKET_NAME == "":
        raise Exception("There is not bucket target defined")
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    logger.info(f"Start processing file ${key} from bucket ${bucket}")
    try:
        message = s3_client.get_object(Bucket=bucket, Key=key)["Body"].read().decode('utf-8').splitlines()
        email = message[0]
        html_message = "\n".join(message[1:])
        email_status = send_report(email, html_message)
        return {
            "statusCode": 200 if email_status else 400,
        }
    except Exception as error:
        logger.error(error)
