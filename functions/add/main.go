/**
 * Route: POST /note
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	util "serverless-notes/utils"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

var sess = session.Must(session.NewSession())
var svc = dynamodb.New(sess)
var tableName = os.Getenv("NOTES_TABLE")

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Note struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	NodeId    string `json:"node_id"`
	Timestamp int64  `json:"timestamp"`
	Expires   int64  `json:"expires"`
}

type DataRequest struct {
	Item Note `json:"Item"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer
	var requestData DataRequest

	err := json.Unmarshal([]byte(request.Body), &requestData)

	userId := util.GetUserId(request.Headers)

	requestData.Item.UserId = userId
	requestData.Item.UserName = util.GetUserName(request.Headers)
	requestData.Item.NodeId = userId + ":" + uuid.New().String()
	requestData.Item.Timestamp = time.Now().Unix()
	requestData.Item.Expires = time.Now().AddDate(0, 0, 90).Unix()

	if err != nil {
		return Response{
			Body:       err.Error(),
			StatusCode: 500,
			Headers:    util.GetResponseHeaders(),
		}, err
	}

	fmt.Println("TABLE NAME ************************************")
	fmt.Println(tableName)

	marshalledItem, err := dynamodbattribute.MarshalMap(requestData.Item)

	if err != nil {
		log.Fatalf("Got error marshalling new movie marshalledItem: %s", err)
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      marshalledItem,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(putItemInput)

	if err != nil {
		fmt.Println("ERROR ******************************")
		fmt.Println(err)
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	noteResponse := Note{}

	err = dynamodbattribute.UnmarshalMap(putItemInput.Item, &noteResponse)

	fmt.Println(noteResponse)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"userId":    noteResponse.UserId,
		"userName":  noteResponse.UserName,
		"expires":   noteResponse.Expires,
		"timestamp": noteResponse.Timestamp,
	})

	json.HTMLEscape(&buf, responseBody)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers:         util.GetResponseHeaders(),
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
