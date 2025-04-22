package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Configuration for Azure AD OBO token exchange
var (
	tenantID     = os.Getenv("AZURE_TENANT_ID")
	clientID     = os.Getenv("AZURE_CLIENT_ID")
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	scope        = os.Getenv("AZURE_OBO_SCOPE")
	tokenURL     = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	resourceURL  = os.Getenv("RESOURCE_API_URL") // e.g., http://localhost:6002/protected
)

type oboResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	print("[gateway-api] Starting Gateway API...\n")
	print("[gateway-api] Environment Variables:\n")
	print("[gateway-api] AZURE_TENANT_ID: ", tenantID, "\n")
	print("[gateway-api] AZURE_CLIENT_ID: ", clientID, "\n")
	print("[gateway-api] AZURE_CLIENT_SECRET: ", clientSecret, "\n")
	print("[gateway-api] AZURE_OBO_SCOPE: ", scope, "\n")
	print("[gateway-api] RESOURCE_API_URL: ", resourceURL, "\n")
	print("[gateway-api] Token URL: ", tokenURL, "\n")
	print("[gateway-api] Resource URL: ", resourceURL, "\n")
	print("[gateway-api] Starting HTTP server...\n")
	http.HandleFunc("/call-resource", handleRequest)
	log.Println("[gateway-api] Listening on :6001")
	log.Fatal(http.ListenAndServe(":6001", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	tokenA := r.Header.Get("Authorization")
	if tokenA == "" {
		http.Error(w, "missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenB, err := exchangeToken(tokenA)
	if err != nil {
		http.Error(w, "failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := callDownstreamAPI(tokenB)
	if err != nil {
		http.Error(w, "failed to call downstream API: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func exchangeToken(tokenA string) (string, error) {
	data := map[string]string{
		"client_id":           clientID,
		"client_secret":       clientSecret,
		"grant_type":          "urn:ietf:params:oauth:grant-type:jwt-bearer",
		"requested_token_use": "on_behalf_of",
		"scope":               scope,
		"assertion":           tokenA[len("Bearer "):],
	}

	form := make([]byte, 0)
	for k, v := range data {
		form = append(form, []byte(fmt.Sprintf("%s=%s&", k, v))...)
	}
	form = form[:len(form)-1] // remove last &

	req, _ := http.NewRequest("POST", tokenURL, bytes.NewBuffer(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return "", fmt.Errorf("azure ad error: %s", string(body))
	}

	var obo oboResponse
	if err := json.Unmarshal(body, &obo); err != nil {
		return "", err
	}
	return obo.AccessToken, nil
}

func callDownstreamAPI(token string) ([]byte, error) {
	req, _ := http.NewRequest("GET", resourceURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
