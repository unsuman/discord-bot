package ollamaorcalite

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ChatAI(prompt string) []string {
	url := "http://localhost:11434/api/generate"
	data := fmt.Sprintf(`{
		"model": "orca-mini",
		"prompt": "%s"
	}`, prompt)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s := string(body)

	scanner := bufio.NewScanner(strings.NewReader(s))

	var jsonResponses []string

	for scanner.Scan() {
		jsonResponses = append(jsonResponses, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	return jsonResponses
}
