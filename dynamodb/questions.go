package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Question struct {
	MessageId int    `json:"message_id"`
	Question  string `json:"question"`
}

//DynamoDBから文字列を取得する
func getQuestions() []Question {
	var questions []Question = []Question{}

	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-northeast-1"))

	input := &dynamodb.ScanInput{
		TableName: aws.String("coaching_words"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		fmt.Println("[GetItem Error]", err)
		questions = append(questions, Question{MessageId: 0, Question: "DBにアクセスできませんでした"})
	}

	for _, question := range result.Items {
		var questionTmp Question
		_ = dynamodbattribute.UnmarshalMap(question, &questionTmp)
		questions = append(questions, questionTmp)
	}
	return questions
}
