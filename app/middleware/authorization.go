package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	httpErrors "app/errors"
	"github.com/Nerzal/gocloak/v13"
	_ "github.com/gorilla/mux"
)

const clientCredentialsGrantType = "client_credentials"

var (
	clientId     string
	clientSecret string
	realm        string
)

var client *gocloak.GoCloak

func InitializeOauthServer(cfg *AuthorizationConfig) {
	client = gocloak.NewClient(cfg.Hostname)
	client.RestyClient().Debug = true
	clientId = cfg.ClientId
	clientSecret = cfg.ClientSecret
	realm = cfg.Realm
}

type AuthorizationConfig struct {
	ClientId     string
	ClientSecret string
	Realm        string
	Hostname     string
}

func Authorization(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//
		// Step 1: Inspect user access token
		//
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 1 {
			w.WriteHeader(401)
			_ = json.NewEncoder(w).Encode(httpErrors.UnauthorizedError())
			return
		}
		accessToken := strings.Split(authHeader, " ")[1]

		introSpectTokenResult, err := client.RetrospectToken(r.Context(), accessToken, clientId, clientSecret, realm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(httpErrors.BadRequestError(err.Error()))
			return
		}
		if !*introSpectTokenResult.Active {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(httpErrors.TokenNotActive())
			return
		}

		//
		// Step 2: Get client token
		//
		clientToken, err := client.GetToken(r.Context(), realm, gocloak.TokenOptions{
			ClientID:     &clientId,
			ClientSecret: &clientSecret,
			GrantType:    gocloak.StringP(clientCredentialsGrantType),
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(httpErrors.BadRequestError(err.Error()))
			return
		}

		//
		// Step 3: Get resource by URL
		//
		resources, err := client.GetResourcesClient(r.Context(), clientToken.AccessToken, realm, gocloak.GetResourceParams{
			URI: gocloak.StringP(r.RequestURI),
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(httpErrors.BadRequestError(err.Error()))
			return
		}
		if len(resources) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(httpErrors.UnauthorizedError())
			return
		}

		//
		// Step 4: Get decision on resource access from the authorization server
		//
		decision, err := client.GetRequestingPartyPermissionDecision(r.Context(), accessToken, realm, gocloak.RequestingPartyTokenOptions{
			Audience:    gocloak.StringP(clientId),
			Permissions: &[]string{*resources[0].Name},
		})
		apiErr := &gocloak.APIError{}
		if errors.As(err, &apiErr) {
			if apiErr.Code == http.StatusForbidden {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(httpErrors.UnauthorizedError())
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(httpErrors.BadRequestError(err.Error()))
				return
			}
		}
		if decision.Result != nil && *decision.Result {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(httpErrors.UnauthorizedError())
		return
	})

}
