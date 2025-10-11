package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ✅ Global constants (accessible from other packages)
var UploadPath = "./uploads"
var BaseURL = "https://localhost:8080"

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTAccessSecret    string
	JWTRefreshSecret   string
	JWTAccessTTLHours  int
	JWTRefreshTTLHours int

	// ✅ Redis Config
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// ✅ Razorpay Keys
	RazorpayKey    string
	RazorpaySecret string

	// ✅ SMTP Config
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromName  string
	SMTPFromEmail string

	// ✅ WAF Config
	WafEnable             string
	ProcessRequestHeaders string
	ProcessRequestBody    string

	// ✅ HTTPS/TLS Config
	EnableHTTPS string
	TLSCertPath string
	TLSKeyPath  string
}

// Load reads environment variables and returns a Config object
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using environment variables")
	}

	accessTTL, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL_HOURS"))
	refreshTTL, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL_HOURS"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	// Default values for HTTPS
	enableHTTPS := os.Getenv("ENABLE_HTTPS")
	if enableHTTPS == "" {
		enableHTTPS = "false"
	}

	tlsCertPath := os.Getenv("TLS_CERT_PATH")
	if tlsCertPath == "" {
		tlsCertPath = "./certs/server.crt"
	}

	tlsKeyPath := os.Getenv("TLS_KEY_PATH")
	if tlsKeyPath == "" {
		tlsKeyPath = "./certs/server.key"
	}

	return &Config{
		Port: os.Getenv("PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		JWTAccessSecret:    os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:   os.Getenv("JWT_REFRESH_SECRET"),
		JWTAccessTTLHours:  accessTTL,
		JWTRefreshTTLHours: refreshTTL,

		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,

		RazorpayKey:    os.Getenv("RAZORPAY_KEY_ID"),
		RazorpaySecret: os.Getenv("RAZORPAY_KEY_SECRET"),

		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFromName:  os.Getenv("SMTP_FROM_NAME"),
		SMTPFromEmail: os.Getenv("SMTP_FROM_EMAIL"),

		WafEnable:             os.Getenv("WAF_ENABLE"),
		ProcessRequestHeaders: os.Getenv("PROCESS_REQUEST_HEADERS"),
		ProcessRequestBody:    os.Getenv("PROCESS_REQUEST_BODY"),

		EnableHTTPS: enableHTTPS,
		TLSCertPath: tlsCertPath,
		TLSKeyPath:  tlsKeyPath,
	}
}