package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

<<<<<<< HEAD
// ✅ Global constants (accessible from other packages)
var UploadPath = "./uploads"
var BaseURL   = "http://localhost:8080"

=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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
}

<<<<<<< HEAD
// Load reads environment variables and returns a Config object
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using environment variables")
	}

	accessTTL, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL_HOURS"))
	refreshTTL, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL_HOURS"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

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

<<<<<<< HEAD
=======
		// ✅ SMTP fields
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFromName:  os.Getenv("SMTP_FROM_NAME"),
		SMTPFromEmail: os.Getenv("SMTP_FROM_EMAIL"),
	}
<<<<<<< HEAD
}
=======
}
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
