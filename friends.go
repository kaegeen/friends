package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Function to get Facebook friends count
func getFacebookFriendsCount(accessToken string) (int, error) {
	// Facebook Graph API endpoint to get friends list
	url := fmt.Sprintf("https://graph.facebook.com/me/friends?access_token=%s", accessToken)

	// Send HTTP GET request to Facebook API
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Facebook API returns the friends count under 'summary'
	if summary, ok := result["summary"].(map[string]interface{}); ok {
		if friendsCount, ok := summary["total_count"].(float64); ok {
			return int(friendsCount), nil
		}
	}
	return 0, fmt.Errorf("could not parse friends count from Facebook API response")
}

// Function to get Twitter followers count
func getTwitterFollowersCount(bearerToken string) (int, error) {
	// Twitter API endpoint to get followers count
	url := "https://api.twitter.com/2/users/by/username/YOUR_TWITTER_USERNAME?user.fields=public_metrics"

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Twitter API returns the followers count under 'public_metrics'
	if data, ok := result["data"].(map[string]interface{}); ok {
		if metrics, ok := data["public_metrics"].(map[string]interface{}); ok {
			if followersCount, ok := metrics["followers_count"].(float64); ok {
				return int(followersCount), nil
			}
		}
	}
	return 0, fmt.Errorf("could not parse followers count from Twitter API response")
}

func main() {
	// Set your access tokens here
	facebookAccessToken := os.Getenv("FACEBOOK_ACCESS_TOKEN")
	twitterBearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	// Get Facebook friends count
	fbCount, err := getFacebookFriendsCount(facebookAccessToken)
	if err != nil {
		log.Fatalf("Error getting Facebook friends count: %v", err)
	}
	fmt.Printf("You have %d friends on Facebook.\n", fbCount)

	// Get Twitter followers count
	twitterCount, err := getTwitterFollowersCount(twitterBearerToken)
	if err != nil {
		log.Fatalf("Error getting Twitter followers count: %v", err)
	}
	fmt.Printf("You have %d followers on Twitter.\n", twitterCount)
}
