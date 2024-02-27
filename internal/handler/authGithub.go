package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/internal/models"
	"io/ioutil"
	"log"
	"net/http"
)

var githubOAuthConfig = oauth2{
	ClientID:     "961979141895b40b95d0",
	ClientSecret: "356cfa2cedd31ffae057c1e1c243f01e27513636",
	RedirectURL:  "https://localhost:8080/auth/github/callback",
	AuthURL:      "https://github.com/login/oauth/authorize",
	TokenURL:     "https://github.com/login/oauth/access_token",
	UserInfoURL:  "https://api.github.com/user",
}

func (h *Handler) githubLogin(w http.ResponseWriter, r *http.Request) {
	// Get the environment variable
	githubClientID := "961979141895b40b95d0"

	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"https://localhost:8080/login/github/callback")

	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (h *Handler) githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Set return type JSON
	// w.Header().Set("Content-type", "application/json")

	var user_data models.GithubUserData
	if err := json.Unmarshal([]byte(githubData), &user_data); err != nil {
		log.Panic("JSON parse error:", err)
	}
	if len(user_data.NodeID) < 1 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	models.InfoLog.Printf("\n        Login: %s\n        Nodel: %s\n        ID:    %d\n        Status:%s\n",
		user_data.Login, user_data.NodeID, user_data.ID, "OAuth Github")
	token, expired, err := h.Service.Auth.CreateOrLoginByGithub(user_data)
	if err != nil {
		info := models.InfoSign{
			Error:    err.Error(),
			Username: user_data.Login,
			Password: user_data.NodeID,
		}
		if err := h.Temp.ExecuteTemplate(w, "signin.html", info); err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: expired,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getGithubAccessToken(code string) string {

	clientID := githubOAuthConfig.ClientID
	clientSecret := githubOAuthConfig.ClientSecret

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		githubOAuthConfig.TokenURL,
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Response body converted to stringified JSON
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		githubOAuthConfig.UserInfoURL,
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}
