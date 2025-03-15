# MIT License
# Copyright (c) 2025 Toni Liesche
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

import json
import logging
import os

logger = logging.getLogger()
logger.setLevel(logging.INFO)

FAIL_ALL_MESSAGES = os.getenv("FAIL_ALL_MESSAGES", "false").lower() == "true"

def lambda_handler(event, context):
    failed_messages = []

    for record in event.get("Records", []):
        if "messageId" in record:
            if FAIL_ALL_MESSAGES:
                failed_messages.append({"itemIdentifier": record["messageId"]})

    logger.info("Received event: %s", json.dumps(event, indent=2))

    return {
        "batchItemFailures": failed_messages
    }