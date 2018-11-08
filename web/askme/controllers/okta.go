package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	verifier "github.com/okta/okta-jwt-verifier-golang"
)

// OktaController handles OKTA login functionality
type OktaController struct {
	framework.Router
	framework.Renderer
	config *framework.Config
	smgr   framework.SessionManager
	user.AuthUsecase
}

var (
	state = "ApplicationState"
	nonce = "NonceNotSetYet"
)

// NewOktaController creates OKTA Controller
func NewOktaController(
	rtr framework.Router,
	rndr framework.Renderer,
	config *framework.Config,
	smgr framework.SessionManager,
	authUC user.AuthUsecase) *OktaController {

	c := &OktaController{
		Router:      rtr,
		Renderer:    rndr,
		config:      config,
		smgr:        smgr,
		AuthUsecase: authUC,
	}

	c.Get("/login", c.oktaLogin)
	c.Get("/authorization-code/callback", c.callback)
	c.Get("/logout", c.logout)
	return c
}

func (o *OktaController) oktaLogin(cxt framework.Context) {
	if isAuthenticated(cxt.Session()) {
		redir := fmt.Sprintf("/u/%s", cxt.User().ID)
		cxt.Redirect(redir, http.StatusFound)
	}

	cxt.Session().Set("redir", cxt.Request().URL.Query().Get("redir"))

	issuerParts, _ := url.Parse(os.Getenv("OKTA_ISSUER"))
	baseUrl := issuerParts.Scheme + "://" + issuerParts.Hostname()

	o.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "login.body",
			Title:    "Login",
			HeadTmpl: "login.head",
			Bag: framework.Map{
				"BaseUrl":  baseUrl,
				"ClientId": os.Getenv("OKTA_CLIENT_ID"),
				"Issuer":   os.Getenv("OKTA_ISSUER"),
				"State":    state,
				"Nonce":    nonce,
			}})
}

func (o *OktaController) callback(cxt framework.Context) {
	// Check the state that was returned in the query string is the same as the above state
	if cxt.Request().URL.Query().Get("state") != state {
		cxt.ResponseWriter().Write([]byte("The state was not as expected"))
		return
	}
	// Make sure the code was provided
	if cxt.Request().URL.Query().Get("code") == "" {
		cxt.ResponseWriter().Write([]byte("The code was not returned or is not accessible"))
		return
	}

	exchange := exchangeCode(cxt.Request().URL.Query().Get("code"), cxt.Request())

	_, verificationError := verifyToken(exchange.IdToken)

	if verificationError != nil {
		fmt.Println(verificationError)
	}

	if verificationError == nil {
		cxt.Session().Set("id_token", exchange.IdToken)
		cxt.Session().Set("access_token", exchange.AccessToken)
	}

	user := getProfileData(cxt.Session())

	userObj := &models.User{Email: user["email"], Name: user["name"]}

	cxt.Session().Set("user", userObj)
	cxt.SetUser(&framework.User{ID: string(user["email"]), Name: user["name"]})

	redir, _ := cxt.Session().Get("redir").(string)
	if len(redir) == 0 {
		redir = fmt.Sprintf("/u/%s", user["email"])
	}

	searchResult, err := o.FindUserByEmail(user["email"])
	if err != nil && searchResult == nil {
		_, err := o.Signup(user["email"], "defaultPassword", user["name"])
		if err != nil {
			cxt.ResponseWriter().Write([]byte(fmt.Sprintf("Fail to create user: %v", user["email"])))
			return
		}
	}

	cxt.Redirect(redir, http.StatusFound)
}

func (o *OktaController) logout(cxt framework.Context) {
	o.smgr.Abandon(cxt)
	cxt.Redirect("/", http.StatusFound)
}

func exchangeCode(code string, r *http.Request) Exchange {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(os.Getenv("OKTA_CLIENT_ID") + ":" + os.Getenv("OKTA_CLIENT_SECRET")))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", "http://localhost:8080/authorization-code/callback")

	oktaURL := os.Getenv("OKTA_ISSUER") + "/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", oktaURL, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var exchange Exchange
	json.Unmarshal(body, &exchange)

	return exchange

}

func isAuthenticated(session *framework.Session) bool {
	if session.Get("id_token") == nil || session.Get("id_token") == "" {
		return false
	}

	return true
}

func getProfileData(session *framework.Session) map[string]string {
	m := make(map[string]string)

	reqUrl := os.Getenv("OKTA_ISSUER") + "/v1/userinfo"

	req, _ := http.NewRequest("GET", reqUrl, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+session.Get("access_token").(string))
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(body, &m)

	return m
}

func verifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = os.Getenv("OKTA_CLIENT_ID")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("OKTA_ISSUER"),
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}

type Exchange struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	IdToken          string `json:"id_token,omitempty"`
}
