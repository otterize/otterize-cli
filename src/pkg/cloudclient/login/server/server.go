package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/auth_api"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/authenticator"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

const (
	host = "localhost"
	port = 52744
)

type AuthResult struct {
	AccessToken string
	Expiry      time.Time
	Profile     map[string]interface{}
}

// A random state should be sent with the authorize request, and verified when getting a callback.
// this prevents an attack in which a client will be logged into an attacker's account by clicking a malicious link,
// that may cause an information leak (if for example the client applied intents into the attacker account)
// More details: https://stackoverflow.com/a/35988614/10574201
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)
	return state, nil
}

type LoginServer struct {
	auth0Conf      auth_api.Auth0Config
	state          string
	auth           *authenticator.Authenticator
	authResultChan chan *AuthResult
}

func NewLoginServer(auth0Conf auth_api.Auth0Config) LoginServer {
	callbackUrl := fmt.Sprintf("http://%s:%d/callback", host, port)
	return LoginServer{
		auth0Conf:      auth0Conf,
		state:          must.MustRet(generateRandomState()),
		auth:           must.MustRet(authenticator.New(auth0Conf.ClientId, auth0Conf.Domain, callbackUrl)),
		authResultChan: make(chan *AuthResult, 1),
	}
}

func (l *LoginServer) GetLoginUrl(forceRelogin bool) string {
	options := []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("audience", l.auth0Conf.Audience)}
	if forceRelogin {
		options = append(options, oauth2.SetAuthURLParam("max_age", "0"))
	}
	return l.auth.AuthCodeURL(l.state, options...)
}

func (l *LoginServer) callback(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("state") != l.state {
		http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
		return
	}

	if req.URL.Query().Get("code") == "" {
		logrus.Errorf("Request is missing auth code. Request params: %s", req.URL.Query())
		queryJson, err := json.MarshalIndent(req.URL.Query(), "", "    ")
		must.Must(err)
		http.Error(w, fmt.Sprintf("Invalid response code: %s", queryJson), http.StatusBadRequest)
		return
	}

	// Exchange an authorization code for a token.
	token, err := l.auth.Exchange(req.Context(), req.URL.Query().Get("code"))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to convert an authorization code into a token.", http.StatusUnauthorized)
		return
	}

	idToken, err := l.auth.VerifyIDToken(req.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify id Token.", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, "Couldn't extract profile from idToken", http.StatusInternalServerError)
		return
	}

	l.authResultChan <- &AuthResult{
		AccessToken: token.AccessToken,
		Expiry:      token.Expiry,
		Profile:     profile,
	}
	_, _ = io.WriteString(w, "Login completed successfully. You can close this window now.")
}

func (l *LoginServer) Start() {
	http.HandleFunc("/callback", l.callback)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil); err != nil {
			logrus.Fatalf("There was an error with the http server: %v", err)
		}
	}()
}

func (l *LoginServer) GetAuthResultChannel() <-chan *AuthResult {
	return l.authResultChan
}
