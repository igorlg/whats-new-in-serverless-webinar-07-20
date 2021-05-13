import json

import requests
from aws_lambda_powertools import Logger, Tracer, Metrics


logger = Logger()
tracer = Tracer()
metrics = Metrics()

CHECKIP_URL = "https://checkip.amazonaws.com"

@logger.inject_lambda_context
@tracer.capture_lambda_handler
@metrics.log_metrics(capture_cold_start_metric=True)
def handler(event, context):
    logger.info("Starting function")

    r = requests.get(CHECKIP_URL)
    logger.debug({'checkip': r.text})
    ip = r.text

    return {
        'body': json.dumps(f'Hello {ip}'),
        'statusCode': 200
    }
