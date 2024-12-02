package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aldotp/golang-login-with-google/config"
	"github.com/aldotp/golang-login-with-google/internal/model"
	"github.com/aldotp/golang-login-with-google/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	// "google.golang.org/api/oauth2/v2"
	// "google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthService defines authentication logic
type AuthService interface {
	AuthenticateGoogle(token string) (*model.User, error)
	RegisterGoogle(token string) (*model.User, error)
	GenerateGoogleAuthURL() string
	ExchangeGoogleCodeForToken(code string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}
func (s *authService) AuthenticateGoogle(token string) (*model.User, error) {
	data, err := s.extractEmailFromGoogleToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindUserByEmail(data.Email)
	if err == nil {
		return user, nil
	}

	return nil, errors.New("user not found")
}

// GenerateGoogleAuthURL generates the URL for Google OAuth2 login

func (s *authService) GenerateGoogleAuthURL() string {
	clientID := config.GetGoogleClientID()
	redirectURI := url.QueryEscape(config.GetGoogleRedirectURI())

	// Generate URL ensuring the query parameters are correctly formatted
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email+profile",
		clientID,
		redirectURI,
	)

	fmt.Println("Generated Google Auth URL:", authURL) // Debugging: Print the full URL
	return authURL
}

func (s *authService) ExchangeGoogleCodeForToken(code string) (string, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     config.GetGoogleClientID(),
		ClientSecret: config.GetGoogleClientSecret(),
		RedirectURL:  config.GetGoogleRedirectURI(),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	ctx := context.Background()
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Extract id_token from the token response
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("id_token not found in token response")
	}

	return idToken, nil
}

// // ExchangeGoogleCodeForToken exchanges an authorization code for an access token
// func (s *authService) ExchangeGoogleCodeForToken(code string) (string, error) {
// 	oauth2Config := &oauth2.Config{
// 		ClientID:     config.GetGoogleClientID(),
// 		ClientSecret: config.GetGoogleClientSecret(),
// 		RedirectURL:  config.GetGoogleRedirectURI(),
// 		Scopes:       []string{"email", "profile"},
// 		Endpoint:     google.Endpoint,
// 	}

// 	ctx := context.Background()
// 	fmt.Println("Exchanging code for token:")
// 	fmt.Println("Code:", code)
// 	fmt.Println("ClientID:", oauth2Config.ClientID)
// 	fmt.Println("RedirectURL:", oauth2Config.RedirectURL)

// 	token, err := oauth2Config.Exchange(ctx, code)
// 	if err != nil {
// 		fmt.Printf("Error during token exchange: %v\n", err)
// 		return "", err
// 	}

// 	fmt.Println("Access Token:", token.AccessToken)
// 	return token.AccessToken, nil
// }

// RegisterGoogle registers a new user using a Google OAuth token
func (s *authService) RegisterGoogle(token string) (*model.User, error) {
	data, err := s.extractEmailFromGoogleToken(token)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	_, err = s.repo.FindUserByEmail(data.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	tNow := time.Now()
	// Create a new user
	newUser := &model.User{
		ID:        uuid.New().String(),
		Email:     data.Email,
		Name:      data.Name,
		Password:  "",
		Role:      "user",
		Provider:  "google",
		CreatedAt: tNow,
		UpdatedAt: tNow,
	}

	err = s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) extractEmailFromGoogleToken(idToken string) (*model.User, error) {
	certs, err := fetchGooglePublicKeys()
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}

		key, ok := certs[kid]
		if !ok {
			return nil, errors.New("key not found for kid")
		}

		return jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, _ := claims["email"].(string)
		// if !ok {
		// 	return nil, errors.New("email not found in token claims")
		// }

		name, _ := claims["name"].(string)

		fmt.Println("email:", email)
		fmt.Println("name:", name)
		// if !ok {
		// 	return nil, errors.New("name not found in token claims")
		// }

		data := &model.User{
			Email: email,
			Name:  name,
		}

		return data, nil
	}

	return nil, errors.New("invalid token")
}

func fetchGooglePublicKeys() (map[string]string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return nil, errors.New("unable to fetch Google public keys")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var keysData map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&keysData); err != nil {
		return nil, errors.New("unable to parse Google public keys")
	}

	return keysData, nil
}

// fetchGooglePublicKeys retrieves Google's public keys for JWT validation
// func fetchGooglePublicKeys() (map[string]string, error) {
// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/certs")
// 	if err != nil {
// 		return nil, errors.New("unable to fetch Google public keys")
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	var certs map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
// 		return nil, errors.New("unable to parse Google public keys")
// 	}

// 	keys := make(map[string]string)
// 	for kid, key := range certs {
// 		keys[kid] = fmt.Sprintf("%v", key)
// 	}

// 	return keys, nil
// }
