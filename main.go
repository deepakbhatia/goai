package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"log"
)

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	model   string
	input   string
	output  string
	apiType string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "ai-client",
		Short: "CLI client for OpenAI and Hugging Face APIs",
		Run:   execute,
	}

	rootCmd.PersistentFlags().StringVar(&model, "model", "", "Name of the model")
	rootCmd.PersistentFlags().StringVar(&input, "input", "", "Input text or image")
	rootCmd.PersistentFlags().StringVar(&output, "output", "", "Output directory")
	rootCmd.PersistentFlags().StringVar(&apiType, "api", "", "API type: openai or huggingface")

	rootCmd.MarkPersistentFlagRequired("model")
	rootCmd.MarkPersistentFlagRequired("input")
	rootCmd.MarkPersistentFlagRequired("output")
	rootCmd.MarkPersistentFlagRequired("api")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(cmd *cobra.Command, args []string) {
	client := resty.New()

	var response string
	var err error

	switch apiType {
	case "openai":
		response, err = callOpenAI(client, model, input)
	case "huggingface":
		response, err = callHuggingFace(client, model, input)
	default:
		fmt.Println("Invalid API type. Use 'openai' or 'huggingface'.")
		return
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := ioutil.WriteFile(filepath.Join(output, "response.txt"), []byte(response), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Response written to", filepath.Join(output, "response.txt"))
}

func callOpenAI(client *resty.Client, model, prompt string) (string, error) {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")

	messages := make([]map[string]string, 0)

	promptMap := map[string]string{
		"role":    "user",
		"content": prompt,
	}

	messages = append(messages, promptMap)
	log.Println((messages[0]))
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openaiAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":    model,
			"messages": messages,
		}).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func callHuggingFace(client *resty.Client, model, image string) (string, error) {
	huggingFaceAPIKey := os.Getenv("HUGGING_FACE_TOKEN")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+huggingFaceAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": model,
			"image": image,
		}).
		Post("https://api-inference.huggingface.co/models/" + model)

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
