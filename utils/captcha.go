package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// TurnstileResponse represents the response from Cloudflare Turnstile API
type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
}

// CaptchaService handles CAPTCHA verification
type CaptchaService struct {
	secretKey  string
	enabled    bool
	httpClient *http.Client
}

// CaptchaConfig holds CAPTCHA configuration
type CaptchaConfig struct {
	SecretKey string
	Enabled   bool
}

// NewCaptchaService creates a new CAPTCHA service
func NewCaptchaService(config CaptchaConfig) *CaptchaService {
	return &CaptchaService{
		secretKey: config.SecretKey,
		enabled:   config.Enabled,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// LoadCaptchaFromEnv loads CAPTCHA configuration from environment variables
func LoadCaptchaFromEnv() *CaptchaService {
	enabled := os.Getenv("VITE_ENABLE_TURNSTILE") != "false"
	secretKey := os.Getenv("TURNSTILE_SECRET_KEY")

	if enabled && secretKey == "" {
		fmt.Println("❌ CAPTCHA enabled but TURNSTILE_SECRET_KEY is NOT set")
	}

	if !enabled {
		fmt.Println("ℹ️ CAPTCHA verification disabled via VITE_ENABLE_TURNSTILE=false")
	}

	return NewCaptchaService(CaptchaConfig{
		SecretKey: secretKey,
		Enabled:   enabled,
	})
}

// VerifyToken verifies a Cloudflare Turnstile token
func (s *CaptchaService) VerifyToken(token string, remoteIP string) error {

	// Skip verification if disabled
	if !s.enabled {
		fmt.Println("⚠️ CAPTCHA verification skipped (disabled)")
		return nil
	}

	// Validate inputs
	if token == "" {
		return errors.New("CAPTCHA token is required")
	}

	if s.secretKey == "" {
		return errors.New("CAPTCHA secret key is not configured")
	}

	// Prepare form-urlencoded payload (REQUIRED by Cloudflare)
	data := url.Values{}
	data.Set("secret", s.secretKey)
	data.Set("response", token)

	if remoteIP != "" {
		data.Set("remoteip", remoteIP)
	}

	// Create request
	req, err := http.NewRequest(
		"POST",
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create CAPTCHA verification request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("CAPTCHA verification request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read CAPTCHA verification response: %w", err)
	}

	// Parse response
	var result TurnstileResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse CAPTCHA verification response: %w", err)
	}

	// Check result
	if !result.Success {
		if len(result.ErrorCodes) > 0 {
			return fmt.Errorf("CAPTCHA verification failed: %v", result.ErrorCodes)
		}
		return errors.New("CAPTCHA verification failed")
	}

	fmt.Printf("✅ CAPTCHA verified successfully (hostname=%s)\n", result.Hostname)
	return nil
}

// IsEnabled returns whether CAPTCHA verification is enabled
func (s *CaptchaService) IsEnabled() bool {
	return s.enabled
}
