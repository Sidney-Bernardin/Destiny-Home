package destinyhome

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)

	// Load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("couldn't load .env file: %v", err)
	}

	// Setup environment variables.
	apiKey = os.Getenv("BUNGIE_API_KEY")
	projectID = os.Getenv("PROJECT_ID")
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
}
