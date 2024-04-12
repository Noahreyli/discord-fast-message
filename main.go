package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	tokens, err := readTokens("token.txt")
	if err != nil {
		fmt.Println("トークンファイルを読み込めません:", err)
		return
	}

	var channelID string
	var messageContent string
	var addRandomString bool
        
	fmt.Print("Created by rucykun\n")

	fmt.Print("チャンネルIDを入力してください: ")
	_, _ = fmt.Scanln(&channelID)

	fmt.Print("送信するメッセージを入力してください: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	messageContent = scanner.Text()

	fmt.Print("メッセージの語尾に3文字のランダムな英数字を追加しますか？ (y/n): ")
	var consent string
	_, _ = fmt.Scanln(&consent)
	addRandomString = strings.ToLower(consent) == "y"

	var messageCount int
	fmt.Print("送信回数を入力してください: ")
	_, _ = fmt.Scanln(&messageCount)

	var wg sync.WaitGroup

	for _, token := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			for i := 0; i < messageCount; i++ {
				finalMessage := messageContent
				if addRandomString {
					randomString := generateRandomString(3)
					finalMessage += " " + randomString
				}
				success, err := postMessage(token, channelID, finalMessage)
				if success {
					fmt.Printf("[✓]Success %s\n", token[:10])
				} else {
					fmt.Printf("[✗]Failed %s - %s\n", token[:10], err)
				}
			}
		}(token)
	}

	wg.Wait()
}

func readTokens(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token := scanner.Text()
		tokens = append(tokens, token)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func postMessage(token, channelID, content string) (bool, error) {
	url := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", channelID)
	payload := fmt.Sprintf(`{"content":"%s"}`, content)
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
