import json
import os
import uuid
from typing import List, Dict

import boto3
import pymysql
import logging
import sys

RDS_HOST = os.getenv("RDS_HOST")
RDS_USER_NAME = os.getenv("RDS_USER_NAME")
RDS_PASSWORD = os.getenv("RDS_PASSWORD")
DB_NAME = os.getenv("DB_NAME")
AUTH = os.getenv("AUTH")
EMAIL_BUCKET_NAME = os.getenv("EMAIL_BUCKET_NAME", "")

logger = logging.getLogger()
logger.setLevel(logging.INFO)

try:
    conn = pymysql.connect(host=RDS_HOST, user=RDS_USER_NAME, passwd=RDS_PASSWORD, db=DB_NAME, connect_timeout=5)
except pymysql.MySQLError as e:
    logger.error("ERROR: Unexpected error: Could not connect to MySQL instance.")
    logger.error(e)
    sys.exit()

logger.info("Database connection successfully")

s3_client = boto3.client("s3")

logger.info("SMTP connection successfully")


class MonthTransactions:
    __month: str
    __count: int

    def __init__(self, month: str, count: int):
        self.__month = month
        self.__count = count

    def __str__(self):
        return f"Number of transactions in {self.__month}: {self.__count}"


class Summary:
    __balance: float
    __debit: float
    __credit: float

    def __init__(self, balance: float, debit: float, credit: float):
        self.__balance = balance
        self.__debit = debit
        self.__credit = credit

    def __str__(self):
        return "Total Balance {:.2f}".format(self.__balance) + "\n Average debit amount: {:.2f}".format(self.__debit) + \
               "\n Average credit amount: {:.2f}".format(self.__credit)


def get_transactions() -> List:
    month_transactions: List[MonthTransactions] = []
    query: str = "select date_format(transaction_date, \'%M\') as month, " \
                 "COUNT(date_format(transaction_date, \'%M\')) as total_transactions " \
                 "from transactions group by date_format(transaction_date, \'%M\');"

    cur = conn.cursor()
    cur.execute(query)
    data = cur.fetchall()
    for months in data:
        month_transactions.append(MonthTransactions(months[0], int(months[1])))
    return month_transactions


def get_summary() -> Summary:
    query: str = "select sum(amount) as balance, " \
                 "AVG(case when amount < 0 then amount else 0 end) as credit, " \
                 "AVG(case when amount > 0 then amount else 0 end) as debit " \
                 "from transactions;"
    cur = conn.cursor()
    cur.execute(query)
    data = cur.fetchone()
    return Summary(float(data[0]), float(data[1]), float(data[2]))


def write_lines_to_files(message: str, email: str) -> str:
    output_file_name: str = f"{uuid.uuid4()}.html"
    with open(f"/tmp/{output_file_name}", 'w', newline='') as outputFile:
        outputFile.write(f"{email}\n")
        outputFile.write(message)
    s3_client.upload_file(f"/tmp/{output_file_name}", EMAIL_BUCKET_NAME, output_file_name)
    return output_file_name


def lambda_handler(event, context):
    """Reports transactions to an email

    Parameters
    ----------
    event: dict, required
        S3 bucket PutObject Event

    context: object, required
        Lambda Context runtime methods and attributes

    Returns
    ------
    API Gateway Lambda Proxy Output Format: dict

    """
    logger.info(event)
    auth_header = event["headers"].get("auth-header", "invalid")
    if auth_header != AUTH:
        return {
            "statusCode": 401,
        }
    try:
        body_request: Dict = json.loads(event["body"])
    except Exception as error:
        logger.error(error)
        return {
            "statusCode": 404,
        }
    email = body_request.get("email")
    if not email:
        return {
            "statusCode": 404,
        }

    transactions = get_transactions()
    summary = get_summary()

    output_transaction = " <br>".join([f"<li>{str(transaction)}</li>" for transaction in transactions])
    html_message = f"""
           <h1>  
           <img src="https://blog.storicard.com/wp-content/uploads/2019/07/Stori-horizontal-11.jpg" />
           </h1>
           <h1> YOUR TRANSACTION REPORT </h1>
           <ul>
           {output_transaction}
           </ul>
           <h2> The summary </h2>
           {summary}
       """
    write_lines_to_files(html_message, email)
    return {
        "statusCode": 201,
    }
