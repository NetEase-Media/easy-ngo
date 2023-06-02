package chatgpt

import (
	"fmt"
	"testing"
)

func TestChat(t *testing.T) {
	chatgpt := &ChatGPT{"sk-xxxx"}
	res, err := chatgpt.Chat("What's your name?")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res=" + res)
}
