# Azure AD Setup Guide for Dapr OAuth2 On-Behalf-Of Flow

This guide walks you through creating and configuring two Azure AD App Registrations needed for the OAuth2 On-Behalf-Of (OBO) flow with Dapr and Go microservices.

---

## 🎯 Goal

| App             | Purpose                                                                 |
|----------------|-------------------------------------------------------------------------|
| **client-app**  | Represents Postman or a frontend app that initiates login and gets **token A** |
| **gateway-api** | Used by the Go API to exchange token A for token B using OBO flow       |

---

## ✅ Step-by-Step Instructions

### 1. Create the `client-app`

1. Go to [Azure Portal – App registrations](https://portal.azure.com/#view/Microsoft_AAD_RegisteredApps/ApplicationsListBlade)
2. Click **"+ New registration"**
3. Fill in:
   - **Name**: `client-app`
   - **Supported account types**: *Accounts in this organizational directory only*
   - **Redirect URI**:
     - Platform: **Web**
     - URI: `https://oauth.pstmn.io/v1/callback`
4. Click **Register**

#### 🔐 Enable Public Client Flows (for Postman)

1. In the `client-app`, go to **Authentication**
2. Scroll to **Advanced Settings**
3. ✅ Check `Allow public client flows`
4. Click **Save**

---

### 2. Create the `gateway-api`

1. Go back to **App registrations**
2. Click **"+ New registration"**
3. Fill in:
   - **Name**: `gateway-api`
   - **Supported account types**: *Accounts in this organizational directory only*
4. Click **Register**

#### 🔓 Expose an API (to define a scope)

1. In the `gateway-api`, go to **Expose an API**
2. Click **Set** to define the **Application ID URI**
   - Example: `api://<gateway-client-id>`
3. Click **"Add a scope"**:
   - Scope name: `access_as_user`
   - Admin consent display name: `Access Gateway API`
   - Admin consent description: `Allows the app to access this API on behalf of the signed-in user`
   - ✅ Enabled: Yes
4. Click **Add scope**

#### 🔑 Create a Client Secret (used by your Go app)

1. Go to **Certificates & secrets**
2. Click **New client secret**
3. Add a name like `obo-flow-secret`
4. Set expiry (e.g., 6 months or 12 months)
5. Click **Add** and immediately copy the **Value** (not the ID!)

---

### 3. Grant Permission for `client-app` to call `gateway-api`

1. Go to the `client-app`
2. Click **API permissions**
3. Click **Add a permission** → **My APIs**
4. Select `gateway-api`
5. Choose **Delegated permissions**, check `access_as_user`, and click **Add permissions**
6. Click **Grant admin consent** (you may need to be an admin)

---

## 🔐 Values You’ll Use in Your .env

```env
AZURE_TENANT_ID=<Directory (tenant) ID>
AZURE_CLIENT_ID=<gateway-api client ID>
AZURE_CLIENT_SECRET=<gateway-api secret VALUE>
AZURE_OBO_SCOPE=api://<gateway-api client ID>/access_as_user
```

You’ll also need the `client-app` client ID to request `token A` in Postman.

---

## ✅ Done!

You can now:

- Use Postman to request `token A` from `client-app`
- Send that to `gateway-api` for an OBO exchange to `token B`
- Validate token B using Dapr middleware on `resource-api`