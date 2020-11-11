#!/bin/bash

gcloud functions deploy DestinyHomeWebhook \
    --runtime go111 \
    --entry-point Webhook \
    --set-env-vars BUNGIE_API_KEY=${BUNGIE_API_KEY} \
    --set-env-vars PROJECT_ID=${PROJECT_ID} \
    --trigger-http

