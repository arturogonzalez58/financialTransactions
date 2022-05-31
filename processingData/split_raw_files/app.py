import logging
import os
from typing import List
import boto3
import urllib
import csv
import uuid

logger = logging.getLogger()
logger.setLevel(logging.INFO)

LINES_PER_FILE = os.getenv("LINES_PER_FILE", "100")
SPLIT_BUCKET_NAME = os.getenv("SPLIT_BUCKET_NAME", "")
TEMP_FILE = '/tmp/output.csv'

s3 = boto3.client('s3')


def write_lines_to_files(lines_to_write: List[str]) -> str:
    output_file_name: str = f"{uuid.uuid4()}.csv"
    with open(TEMP_FILE, 'w', newline='') as outputFile:
        writer = csv.writer(outputFile)
        for output_line in lines_to_write:  # reverse order
            writer.writerow(output_line)
    s3.upload_file(TEMP_FILE, SPLIT_BUCKET_NAME, output_file_name)
    return output_file_name


def lambda_handler(event, context):
    if SPLIT_BUCKET_NAME == "":
        raise Exception("There is not bucket target defined")
    lines_per_file: int = int(LINES_PER_FILE)
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    logger.info(f"Start processing file ${key} from bucket ${bucket}")
    try:
        file_content = s3.get_object(Bucket=bucket, Key=key)["Body"].read().decode('utf-8').splitlines()
        lines_to_write: List[str] = []
        lines = csv.reader(file_content)
        logger.info(f"Files to process in bucket {SPLIT_BUCKET_NAME}:")
        for line in lines:
            lines_to_write.append(line)
            if len(lines_to_write) >= lines_per_file:
                output_file = write_lines_to_files(lines_to_write)
                logger.info(output_file)
                lines_to_write.clear()
        if len(lines_to_write) > 0:
            write_lines_to_files(lines_to_write)

        return {
            "message": "success",
        }
    except Exception as error:
        logger.error(error)
        raise error
