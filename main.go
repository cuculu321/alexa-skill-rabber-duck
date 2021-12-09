package main

import (
	"errors"
	"fmt"

	"rubber-duck/alexa"
	//	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	// ErrInvalidIntent is error-object
	ErrInvalidIntent = errors.New("Invalid intent")
)

/*
 * Functions that control the skill's behavior
 */

// Helpを求められた時の応答
func GetHelperResponse() alexa.Response {
	sessionAttributes := make(map[string]interface{})
	cardTitle := "Helper"
	speechOutput := "あなたの悩みや課題を一緒に考え、考えを進めるための質問をします。"
	repromptText := "また悩み事があれば相談してください。"
	shouldEndSession := false
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// GetFinishRequest is function-type
func GetFinishRequest() alexa.Response {
	sessionAttributes := make(map[string]interface{})
	cardTitle := "Session Ended"
	speechOutput := "解決できていれば何よりです"
	repromptText := ""
	shouldEndSession := true
	fmt.Println(speechOutput)
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// GetNoEntityResponse is function-type
func GetNoEntityResponse() alexa.Response {
	cardTitle := ""
	sessionAttributes := make(map[string]interface{})
	shouldEndSession := false
	speechOutput := "なんでも相談してください"
	repromptText := ""
	fmt.Println(speechOutput)
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

/*
 * Events
 */

// OnSessionStarted is function-type
func OnSessionStarted(sessionStartedRequest map[string]string, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnSessionStarted requestId=" + sessionStartedRequest["requestId"] + ", sessionId=" + session.SessionID)
	return GetNoEntityResponse(), nil
}

// OnLaunch is function-type
func OnLaunch(launchRequest alexa.RequestDetail, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnLaunch requestId=" + launchRequest.RequestID + ", sessionId=" + session.SessionID)
	return GetHelperResponse(), nil
}

func GetIntentResponse() alexa.Response {
	sessionAttributes := make(map[string]interface{})
	cardTitle := "Response"
	speechOutput := "ここの文章に適当な質問を仕込む"
	repromptText := "また悩み事があれば相談してください。"
	shouldEndSession := false
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// OnIntent is function-type
func OnIntent(intentRequest alexa.RequestDetail, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnIntent requestId=" + intentRequest.RequestID + ", sessionId=" + session.SessionID)
	intentName := intentRequest.Intent.Name
	if intentName == "HelloWorldIntent" {
		return GetIntentResponse(), nil
	} else if intentName == "AMAZON.HelpIntent" {
		return GetHelperResponse(), nil
	} else if intentName == "AMAZON.StopIntent" || intentName == "AMAZON.CancelIntent" {
		return GetFinishRequest(), nil
	}
	return alexa.Response{}, ErrInvalidIntent
}

// OnSessionEnded is function-type
func OnSessionEnded(sessionEndedRequest alexa.RequestDetail, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnSessionEnded requestId=" + sessionEndedRequest.RequestID + ", sessionId=" + session.SessionID)
	return GetNoEntityResponse(), nil
}

//起動時の方法
func Handler(event alexa.Request) (alexa.Response, error) {
	fmt.Println("event.session.application.applicationId=" + event.Session.Application.ApplicationID)
	getQuestions()

	eventRequestType := event.Request.Type
	fmt.Println(eventRequestType)
	if event.Session.New {
		return OnSessionStarted(map[string]string{"requestId": event.Request.RequestID}, event.Session)
	} else if eventRequestType == "LaunchRequest" {
		return OnLaunch(event.Request, event.Session)
	} else if eventRequestType == "IntentRequest" {
		//インテント起動フレーズでの起動
		return OnIntent(event.Request, event.Session)
	} else if eventRequestType == "SessionEndedRequest" {
		return OnSessionEnded(event.Request, event.Session)
	}
	return alexa.Response{}, ErrInvalidIntent
}

type Question struct {
	MessageId int    `json:"message_id"`
	Question  string `json:"question"`
	Action    string
	Points    int
	Hidden    bool `dynamo:"-"`
}

func getQuestions() {
	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-northeast-1"))

	input := &dynamodb.ScanInput{
		TableName: aws.String("coaching_words"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		fmt.Println("[GetItem Error]", err)
		return
	}
	fmt.Println(result)
}
func main() {
	//see https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-handler.html
	fmt.Println("hoge")
	getQuestions()
	lambda.Start(Handler)
}

func var_dump(v ...interface{}) {
	for _, vv := range v {
		fmt.Printf("%#v\n", vv)
	}
}
