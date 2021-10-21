package entities

import (
	"fmt"
	"strings"
)

type SecretQuestionAndResponse struct {
	Question string `json:"question"`
	Response string `json:"response"`
}

type SecretQuestionsAndResponses [2]*SecretQuestionAndResponse

func (sqar SecretQuestionsAndResponses) Trim() {
	for _, sq := range sqar {
		fmt.Printf("before %p\n", &sq)
		sq = &SecretQuestionAndResponse{
			Question: strings.TrimSpace(sq.Question),
			Response: strings.TrimSpace(sq.Response),
		}
		fmt.Printf("after %p\n", &sq)
	}
}

func (sqar SecretQuestionsAndResponses) ToLowerResponses() {
	for _, sq := range sqar {
		fmt.Printf("before: %p\n", &sq)
		sq = &SecretQuestionAndResponse{
			Question: sq.Question,
			Response: strings.ToLower(sq.Response),
		}
		fmt.Printf("after: %p\n", &sq)
	}
}
