package ollamaorcalite

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Res struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func ChatAI(prompt string) string {
	url := "http://localhost:11434/api/generate"
	data := fmt.Sprintf(`{
		"model": "orca-mini",
		"prompt": "%s"
	}`, prompt)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	//fmt.Println("Response Body:", string(body))

	s := string(body)

	scanner := bufio.NewScanner(strings.NewReader(s))

	var jsonResponses []string

	for scanner.Scan() {
		jsonResponses = append(jsonResponses, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	var res []Res

	for _, item := range jsonResponses {
		var r Res
		err := json.Unmarshal([]byte(item), &r)
		if err != nil {
			fmt.Println(err.Error())
			return err.Error()
		}

		res = append(res, r)
	}

	r := ""

	for _, item := range res {
		if item.Done {
			break
		}

		fmt.Print(item.Response)
		r += item.Response
	}
	fmt.Println()

	return r
}
