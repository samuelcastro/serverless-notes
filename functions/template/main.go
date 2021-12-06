/**
 * Route: POST /note
 */

package add

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	util "serverless-notes/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess = session.Must(session.NewSession())
var svc = dynamodb.New(sess)
var tableName = os.Getenv("TABLE_NAME")

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	// https://bitfieldconsulting.com/golang/map-string-interface
	errorBody, err := json.Marshal(map[string]interface{}{
		"error":   "Exception",
		"message": "Unknown Error",
	})

	responseBody, _ := json.Marshal(map[string]interface{}{
		"success": "ok",
	})

	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       string(errorBody),
			Headers:    util.GetResponseHeaders(),
		}, err
	}

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
