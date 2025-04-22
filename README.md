# Go OAuth2 OBO Demo with Dapr and Azure AD

This demo showcases a Go microservice architecture using [Dapr](https://dapr.io/) and [Azure Active Directory](https://learn.microsoft.com/en-us/entra/identity-platform/) to implement the **OAuth2 On-Behalf-Of (OBO) flow**.

## 🧱 Services

- **gateway-api**
  - Validates token A (user access token) using Dapr OAuth2 middleware
  - Exchanges it for token B via Azure AD OBO flow
  - Calls the downstream `resource-api` with token B

- **resource-api**
  - Simulates a protected backend service that requires token B

## ⚙️ Running Locally

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- [Dapr CLI](https://docs.dapr.io/getting-started/install-dapr/)
- [Azure AD App Registrations](https://learn.microsoft.com/en-us/entra/identity-platform/howto-create-app-registrations)

### Configure Azure AD (Mock or Real)

Follow instructions in [AZURE.md](AZURE.md)

### Environment Variables

```bash
export AZURE_TENANT_ID=<your-tenant-id>
export AZURE_CLIENT_ID=<your-client-id>
export AZURE_CLIENT_SECRET=<your-client-secret>
export AZURE_OBO_SCOPE=https://graph.microsoft.com/.default
export RESOURCE_API_URL=http://localhost:6002/protected
```

### Run Locally

```bash
dapr run --app-id gateway-api --app-port 6001 --dapr-http-port 3500 \
  --config ./config/config.yaml --components-path ./components \
  go run ./gateway-api/main.go

# In another terminal:
dapr run --app-id resource-api --app-port 6002 \
  go run ./resource-api/main.go
```

---

## 🧪 Testing the Flow

```bash
# Replace <tokenA> with an actual Azure AD access token
curl -H "Authorization: Bearer <tokenA>" http://localhost:6001/call-resource
```

Expected output:

```json
{
  "message": "Access to protected resource granted",
  "user": "token claims would be decoded here"
}
```

---

## 🚀 Kubernetes Deployment

1. Ensure Dapr is installed: `dapr init -k`

2. Apply manifests in `k8s/` folder:

```bash
kubectl apply -f k8s/
```

3. Update `dapr-components.yaml` with your Azure credentials and JWKS URL.

---

## 📄 License
MIT © Diagrid
