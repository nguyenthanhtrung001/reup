package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Function gửi yêu cầu quét Douyin videos
func (uc implUseCase) SentScanDouyinVideos(spaceArr []ArrChannel, isFull bool) error {
	apiURL := "http://10.10.10.189:6080/api/v1/post/scan"

	apiKey := "za1eabkdme138d37e1k76koetza51"

	log.Printf("Sending scan request | Channels count: %d | Full: %v", len(spaceArr), isFull)

	requestBody := ScanRequest{
		Spaces: spaceArr,
		IsFull: isFull,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Scan API failed: %v", err)
		return fmt.Errorf("scan API failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var responseBody ScanResponse
		// Giải mã JSON phản hồi
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			log.Printf("Error decoding response: %v", err)
			return fmt.Errorf("error decoding response: %v", err)
		}

		log.Printf("Scan successful | Message: %s | Queued count: %d", responseBody.Message, responseBody.QueuedCount)
	} else {
		log.Printf("Scan API responded with unexpected status: %d", resp.StatusCode)
		return fmt.Errorf("scan API responded with unexpected status: %d", resp.StatusCode)
	}

	return nil
}
