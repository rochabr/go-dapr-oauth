#!/bin/bash

# -----------------------------
# Azure AD App Registration Automation for Dapr OBO Flow
# -----------------------------

# Set your own values here
RESOURCE_GROUP="my-resource-group"
APP_NAME_PREFIX="dapr-obo-demo"
TENANT_ID=$(az account show --query tenantId -o tsv)

# Create client-app (frontend / Postman)
CLIENT_APP_ID=$(az ad app create \
  --display-name "$APP_NAME_PREFIX-client-app" \
  --reply-urls "https://oauth.pstmn.io/v1/callback" \
  --enable-public-client true \
  --query appId -o tsv)

az ad app update --id $CLIENT_APP_ID --set publicClient=true

# Create gateway-api (API that receives token A)
GATEWAY_APP_ID=$(az ad app create \
  --display-name "$APP_NAME_PREFIX-gateway-api" \
  --identifier-uris "api://$APP_NAME_PREFIX-gateway-api" \
  --query appId -o tsv)

# Create client secret for gateway-api
CLIENT_SECRET=$(az ad app credential reset \
  --id $GATEWAY_APP_ID \
  --append \
  --display-name "obo-secret" \
  --query password -o tsv)

# Expose scope: access_as_user on gateway-api
az ad app permission add --id $GATEWAY_APP_ID \
  --api $GATEWAY_APP_ID \
  --api-permissions $GATEWAY_APP_ID/access_as_user=Scope

# Grant client-app permission to call gateway-api
az ad app permission add --id $CLIENT_APP_ID \
  --api $GATEWAY_APP_ID \
  --api-permissions $GATEWAY_APP_ID/access_as_user=Scope

# Grant admin consent
az ad app permission grant --id $CLIENT_APP_ID --api $GATEWAY_APP_ID --scope access_as_user

# Print outputs
echo "\n✅ Azure App Setup Complete"
echo "Tenant ID: $TENANT_ID"
echo "client-app ID: $CLIENT_APP_ID"
echo "gateway-api ID: $GATEWAY_APP_ID"
echo "gateway-api Secret: $CLIENT_SECRET"
echo "Scope: api://$APP_NAME_PREFIX-gateway-api/access_as_user"
