package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"net/url"
	"encoding/json"
	"io/ioutil"
)

const API_URL = "https://api.pushover.net/1/messages.json"
const USER_ENV_NAME = "PUSHOVER_USER"
const TOKEN_ENV_NAME = "PUSHOVER_TOKEN"
const CURRENT_VERSION = "1.0.0"

func isValidSound(sound string) bool {
	// Sounds list from https://pushover.net/api#sounds (21 + none, as of Sept 2019)
	switch (sound) {
		case
			"alien",
			"bike", "bugle",
			"cashregister", "classical", "climb", "cosmic",
			"echo",
			"falling",
			"gamelan",
			"incoming", "intermission",
			"magic", "mechanical",
			"persistent", "pianobar", "pushover",
			"siren", "spacealarm",
			"tugboar",
			"updown",
			"none":
				return true
	}
	return false
}

func main() {
	var user, token, device, title, sound, messageUrl, messageUrlTitle string
	var priority, timestamp int
	var version bool

	// Required params. Default to environment variables for user & token
	flag.StringVar(&user, "user", os.Getenv(USER_ENV_NAME), "Your personal user ID to identify you as the sender. PUSHOVER_USER env variable also available. (required)")
	flag.StringVar(&token, "token", os.Getenv(TOKEN_ENV_NAME), "The application API token you want to send this message as. PUSHOVER_TOKEN env variable also available (required)")

	// Optional params
	flag.StringVar(&device, "device", "", "The name of the device you want to send the message to directly (optional)")
	flag.StringVar(&title, "title", "", "The title you want to give your message (optional)")
	flag.StringVar(&sound, "sound", "", "Name of the sound to play on the recipient device(s). See https://pushover.net/api#sounds (optional)")
	flag.StringVar(&messageUrl, "url", "", "A URL to show with your message (optional)")
	flag.StringVar(&messageUrlTitle, "url-title", "", "Text for the URL in --url to show (optional)")

	flag.IntVar(&priority, "priority", 0, "-2 no notification; -1 quiet notification; 1 high-priority; 2 require confirmation (optional)")
	flag.IntVar(&timestamp, "timestamp", 0, "Unix timestamp to show to the user (optional)")

	flag.BoolVar(&version, "version", false, "Display the version of this tool and exit (optional)")

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("      pushover-go [options] 'Message to send'")
		fmt.Println("      echo 'Message to send' | ./pushover-go [options]")
		fmt.Println("      pushover-go [options] < file-with-message.txt")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Println("Running version", CURRENT_VERSION)
		os.Exit(0)
	}

	if user == "" || token == "" {
		if user == "" { fmt.Println("No user provided") }
		if token == "" { fmt.Println("No application token provided") }

		flag.PrintDefaults()
		os.Exit(1)
	}

	// Sound name must be valid
	if sound != "" && !isValidSound(sound) {
		fmt.Println("Invalid sound name")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// URL title should only be given if a URL is also given
	if messageUrlTitle != "" && messageUrl == "" {
		fmt.Println("URL title provided but no URL given")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Priority must be -2, -1, 1, or 2 (or not given)
	if priority != 0 && (priority < -2 || priority > 2) {
		fmt.Println("Invalid priority given")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get message to send
	var message string
	switch flag.NArg() {
	case 0:
		// Message is passed as stdin
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading message", err)
			os.Exit(1)
		}
		message = string(input)
	case 1:
		// Message is an argument
		message = flag.Args()[0]
	default:
		fmt.Println("No message provided")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sendPushover(user, token, message, device, title, sound, messageUrl, messageUrlTitle, priority, timestamp)
}

func sendPushover(user, token, message string, device, title, sound, messageUrl, messageUrlTitle string, priority, timestamp int) {
	fmt.Println("Sending message:", message)

	form := url.Values{}
	form.Add("token", token)
	form.Add("user", user)
	form.Add("message", message)

	// Only want to add these optional fields if they were provided and are non-default
	switch {
	case device != "":
		form.Add("device", device)
	case title != "":
		form.Add("title", title)
	case messageUrl != "":
		form.Add("url", messageUrl)
	case messageUrlTitle != "":
		form.Add("url_title", messageUrlTitle)
	case priority != 0:
		form.Add("priority", string(priority))
	case sound != "":
		form.Add("sound", sound)
	case timestamp != 0:
		form.Add("timestamp", string(timestamp))
	}

	response, err := http.PostForm(API_URL, form)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Error sending POST request to Pushover")
		fmt.Println(err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	// `result` map stores values as interface{}, but we know `status` should be float64
	status := result["status"].(float64)

	if status == 1 && response.StatusCode == 200 {
		fmt.Println("Successfully sent message")

		header := response.Header
		fmt.Println("Total Quota:", header["X-Limit-App-Limit"][0])
		fmt.Println("Quota Remaining:", header["X-Limit-App-Remaining"][0])

		os.Exit(0)
	} else {
		fmt.Println("Error sending message, are your application and user tokens correct?")
		fmt.Println("Received data", result)
		os.Exit(1)
	}
}
