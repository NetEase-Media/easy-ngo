package chatgpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ChatGPT struct {
	APIKey string
}

type ChatGPTRequest struct {
	Prompt string `json:"prompt"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func (c *ChatGPT) Chat(prompt string) (string, error) {
	// 1. 构造请求
	req := ChatGPTRequest{
		Prompt: prompt,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// 2. 发送请求
	url := "https://api.openai.com/v1/engines/davinci/completions"
	reqReader := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 3. 解析响应
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respData ChatGPTResponse
	err = json.Unmarshal(respBytes, &respData)
	if err != nil {
		return "", err
	}
	return respData.Choices[0].Text, nil
}

func (c *ChatGPT) ChatWithHistory(prompt string, history []string) (string, error) {
	// 1. 构造请求
	req := ChatGPTRequest{
		Prompt: prompt,
	}
	for _, h := range history {
		req.Prompt += "\nHuman: " + h
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// 2. 发送请求
	url := "https://api.openai.com/v1/engines/davinci/completions"
	reqReader := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 3. 解析响应
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respData ChatGPTResponse
	err = json.Unmarshal(respBytes, &respData)
	if err != nil {
		return "", err
	}
	return respData.Choices[0].Text, nil
}

func (c *ChatGPT) ChatWithTemperature(prompt string, temperature float64) (string, error) {
	// 1. 构造请求
	req := ChatGPTRequest{
		Prompt: prompt,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// 2. 发送请求
	url := "https://api.openai.com/v1/engines/davinci/completions"
	reqReader := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 3. 解析响应
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respData ChatGPTResponse
	err = json.Unmarshal(respBytes, &respData)
	if err != nil {
		return "", err
	}
	return respData.Choices[0].Text, nil
}
