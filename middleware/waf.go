package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// WAFMiddleware inspects request input for malicious patterns (SQLi, XSS, etc.)
func WAFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("----- WAF Middleware Started -----")

		// Use CF-Connecting-IP header if present (real client IP behind Cloudflare)
		clientIP := c.GetHeader("CF-Connecting-IP")
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		fmt.Printf("Request Path: %s, Method: %s, ClientIP: %s\n", c.Request.URL.Path, c.Request.Method, clientIP)

		// 1️⃣ Check query parameters
		for key, values := range c.Request.URL.Query() {
			for _, value := range values {
				if isMalicious(value) {
					blockRequest(c, "query:"+key, value)
					return
				}
			}
		}

		// 2️⃣ Check POST form data
		if strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded") ||
			strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			if err := c.Request.ParseForm(); err == nil {
				for key, values := range c.Request.PostForm {
					for _, value := range values {
						if isMalicious(value) {
							blockRequest(c, "form:"+key, value)
							return
						}
					}
				}
			}
		}

		// 3️⃣ Check JSON body
		if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			var body map[string]interface{}

			// Read body bytes
			buf, err := io.ReadAll(c.Request.Body)
			if err != nil {
				fmt.Println("Failed to read body:", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			// Restore body so handlers can read it later
			c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

			if len(buf) > 0 {
				if err := json.Unmarshal(buf, &body); err == nil {
					if checkJSON(body) {
						blockRequest(c, "JSON body", "")
						return
					}
				} else {
					fmt.Println("JSON unmarshal ignored:", err)
				}
			}
		}

		fmt.Println("----- WAF Middleware Passed, Proceeding to Handler -----")
		c.Next()
	}
}

// isMalicious checks a single string against common attack patterns
func isMalicious(input string) bool {
	patterns := []string{
		`(?i)select\b`, `(?i)union\b`, `(?i)insert\b`, `(?i)update\b`, `(?i)delete\b`, `(?i)drop\b`,
		`<script>`, `</script>`, `--`, `;`, `\bOR\b`, `\bAND\b`,
	}
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			fmt.Printf("Matched malicious pattern: %s in value: %s\n", pattern, input)
			return true
		}
	}
	return false
}

// checkJSON recursively checks all JSON values
func checkJSON(body map[string]interface{}) bool {
	for key, v := range body {
		switch val := v.(type) {
		case string:
			if isMalicious(val) {
				fmt.Printf("Malicious string found in JSON key '%s': %s\n", key, val)
				return true
			}
		case map[string]interface{}:
			if checkJSON(val) {
				return true
			}
		case []interface{}:
			for _, item := range val {
				switch i := item.(type) {
				case map[string]interface{}:
					if checkJSON(i) {
						return true
					}
				case string:
					if isMalicious(i) {
						fmt.Printf("Malicious string in JSON array key '%s': %s\n", key, i)
						return true
					}
				}
			}
		}
	}
	return false
}

// blockRequest aborts the request with 403 Forbidden
func blockRequest(c *gin.Context, key, value string) {
	fmt.Printf("WAF Blocking request! Param: %s, Value: %s\n", key, value)
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error":   "Malicious request blocked",
		"param":   key,
		"value":   value,
		"message": "Request contains forbidden input patterns",
	})
}
