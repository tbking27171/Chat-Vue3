/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chatcat/backend/model"
	"chatcat/backend/pkg/cgpt"
	"chatcat/backend/service"
	"chatcat/backend/service/setting"
	"fmt"
	"github.com/sashabaranov/go-openai"

	"github.com/spf13/cobra"
)

var (
	prompt    string
	modelName string
	token     string
)

// chatgptCmd represents the chatgpt command
var chatgptCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "chatgpt接口测试",
	Long:  `chatgpt gpt3,gpt4 接口测试`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("chatgpt called")
		handleChat()
	},
}

func init() {
	rootCmd.AddCommand(chatgptCmd)
	chatgptCmd.PersistentFlags().StringVar(&modelName, "m", openai.GPT3Dot5Turbo, "模型 gpt3.5 gpt4")
	chatgptCmd.Flags().StringVar(&prompt, "p", "", "提示词")
	chatgptCmd.MarkFlagRequired("p")
	chatgptCmd.MarkFlagRequired("m")
}

func handleChat() {
	app := service.NewApp()
	fmt.Println("OK app")
	generalInfo := setting.New(app).GetGeneralInfo()
	data := generalInfo.Data.(model.Setting)
	token = data.ApiKey
	gpt := cgpt.New(token).WithProxy("http://127.0.0.1:7890")

	// /completions no stream
	//stream, err := gpt.WithModel(openai.GPT3TextDavinci003).
	//	WithPrompt(prompt).
	//	WithStream(false).
	//	WithMaxTokens(0).
	//	WithCompletionRequest().
	//	CompletionNoStream()
	//fmt.Println(stream, err)

	// /completions stream
	gpt.WithModel(openai.GPT3TextDavinci003).
		WithPrompt(prompt).
		WithStream(false).
		WithMaxTokens(0).
		WithCompletionRequest().
		CompletionStream()

	// /chat/completions
	//gpt.WithModel(openai.GPT3Dot5Turbo).WithMessages([]openai.ChatCompletionMessage{
	//	{
	//		Role:    openai.ChatMessageRoleUser,
	//		Content: prompt,
	//	},
	//}).WithChatCompletionRequest().ChatCompletionStream()

}
