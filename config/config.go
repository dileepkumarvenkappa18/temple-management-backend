package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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
}

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

		RedisAddr:     os.Getenv("REDIS_ADDR"),     // e.g., redis:6379
		RedisPassword: os.Getenv("REDIS_PASSWORD"), // blank if none
		RedisDB:       redisDB,                     // usually 0

			// ✅ Razorpay fields loaded
	RazorpayKey:    os.Getenv("RAZORPAY_KEY_ID"),
	RazorpaySecret: os.Getenv("RAZORPAY_KEY_SECRET"),
	}
}
