package entities

import (
	"btradoc/helpers"
	"fmt"
	"testing"
)

func TestStruct(t *testing.T) {
	pwdBody := PasswordUpdateBody{
		Token: helpers.GenerateID(22),
		SecretQuestionsAndResponses: SecretQuestionsAndResponses{
			{
				Question: " PDjejqcbvf",
				Response: " ERZAcDe ",
			},
			{
				Question: " bcnDNEZAZz",
				Response: "efRteAZEC ",
			},
		},
		Password: "12345678910",
	}

	t.Logf("%+v\n", pwdBody)

	pwdBody.SecretQuestionsAndResponses.Trim()
	fmt.Println()
	pwdBody.SecretQuestionsAndResponses.ToLowerResponses()

	for _, pwdbody := range pwdBody.SecretQuestionsAndResponses {
		t.Log(pwdbody.Question, pwdbody.Response)
	}
	t.Logf("%+v\n", pwdBody)
}
