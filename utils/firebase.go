package utils

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var (
	FirebaseApp    *firebase.App
	FirebaseClient *messaging.Client
)

// InitFirebase initializes Firebase Admin SDK and FCM client
func InitFirebase() error {
	ctx := context.Background()

	// Get credentials path from environment
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		credentialsPath = "./serviceAccountKey.json"
	}

	// Check if file exists
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		log.Printf("⚠️ Firebase credentials file not found at: %s", credentialsPath)
		log.Println("ℹ️ Continuing without Firebase (push notifications will be disabled)")
		return nil
	}

	// Create Firebase config with explicit project ID
	config := &firebase.Config{
		ProjectID: "tms-app-38fc7",
	}

	// Initialize Firebase app
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Printf("❌ Error initializing Firebase app: %v", err)
		return err
	}

	log.Printf("✅ Firebase app initialized successfully for project: tms-app-38fc7")

	// Try to get FCM client
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("❌ Error getting FCM client: %v", err)
		log.Println("ℹ️ Continuing without FCM (push notifications will be disabled)")
		// Store app but continue without FCM
		FirebaseApp = app
		FirebaseClient = nil
		return nil
	}

	log.Println("✅ FCM client initialized successfully")

	// Store globally
	FirebaseApp = app
	FirebaseClient = fcmClient

	return nil
}

// GetFCMClient returns the FCM client instance
func GetFCMClient() *messaging.Client {
	return FirebaseClient
}

// IsFCMEnabled checks if FCM is available
func IsFCMEnabled() bool {
	return FirebaseClient != nil
}

// GetFirebaseApp returns the Firebase app instance
func GetFirebaseApp() *firebase.App {
	return FirebaseApp
}