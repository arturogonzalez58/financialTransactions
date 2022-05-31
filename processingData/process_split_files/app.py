import csv
import datetime
import logging
import os
import sys
import urllib
import uuid
from enum import Enum
from typing import List, Tuple

import pymysql
import boto3

RDS_HOST = os.getenv("RDS_HOST")
RDS_USER_NAME = os.getenv("RDS_USER_NAME")
RDS_PASSWORD = os.getenv("RDS_PASSWORD")
DB_NAME = os.getenv("DB_NAME")

logger = logging.getLogger()
logger.setLevel(logging.INFO)

s3 = boto3.client('s3')

try:
    conn = pymysql.connect(host=RDS_HOST, user=RDS_USER_NAME, passwd=RDS_PASSWORD, db=DB_NAME, connect_timeout=5)
except pymysql.MySQLError as e:
    logger.error("ERROR: Unexpected error: Could not connect to MySQL instance.")
    logger.error(e)
    sys.exit()

logger.info("Database connection successfully")


def parse_date(date: str) -> datetime.datetime:
    return datetime.datetime.strptime(date, '%Y-%m-%d %H:%M:%S %z %Z')


class TransactionState(Enum):
    VALID = "valid"
    INVALID = "invalid"


class TransactionModel:
    __id_transaction: str
    __date: datetime.date
    __amount: float
    __job_id: str
    __state: TransactionState

    def is_valid(self):
        return self.__state

    def __init__(self, id_transaction: str, date: str, amount: str, job_id: str):
        self.__id_transaction = id_transaction
        self.__date = parse_date(date)
        self.__amount = float(amount)
        self.__job_id = job_id
        try:
            uuid.UUID(id_transaction)
            self.__state = TransactionState.VALID
        except Exception as error:
            logger.info(f"Not valid data ${error}")
            self.__state = TransactionState.INVALID

    def __iter__(self):
        yield f"transaction-{self.__id_transaction}"
        yield self.__date
        yield self.__amount
        yield str(self.__job_id)
        yield str(self.__state)


def lambda_handler(event, context):
    job_id: str = f"process-job-{uuid.uuid4().hex}"
    logger.info(f"JOB ID: ${job_id} start")
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    try:
        file_content = s3.get_object(Bucket=bucket, Key=key)["Body"].read().decode('utf-8').splitlines()
        transactions: List[Tuple[TransactionModel]] = []
        lines = csv.reader(file_content)
        for line in lines:
            logger.info(line)
            transaction = TransactionModel(line[0], line[1], line[2], str(job_id))
            if transaction.is_valid() == TransactionState.VALID:
                transactions.append(tuple(transaction))
        cursor = conn.cursor()
        query = 'INSERT INTO transactions(transaction_id,transaction_date,amount,job_id,status)' \
                'VALUES( %s, %s, %s, %s, %s) '

        cursor.executemany(query, transactions)
        conn.commit()
        cursor.close()
        logger.info(f"JOB ID: ${job_id} completed")
        conn.close()
    except Exception as error:
        logger.error(error)
        raise error
