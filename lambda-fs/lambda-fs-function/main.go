/*
  Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
  Permission is hereby granted, free of charge, to any person obtaining a copy of this
  software and associated documentation files (the "Software"), to deal in the Software
  without restriction, including without limitation the rights to use, copy, modify,
  merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so.
  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
  INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
  PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
  HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
  SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	efsPath = os.Getenv("EFS_PATH") + "/content"
)

func check(err error) {
	if err != nil {
		println(err)
		panic(err)
	}
}

// getMessages returns the messages in the file
func getMessages() (message string) {
	messages, err := ioutil.ReadFile(efsPath)
	if err != nil {
		return "No messages yet."
	}

	return string(messages)
}

// addMessages checks to see if the file exists, if not
// then it creates a new one and appends the message to
// the file as a new line
func addMessages(m string) {

	f, err := os.OpenFile(efsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	if _, err := f.Write([]byte(m + "\n")); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// deleteMessages deletes the file
func deleteMessages() {
	err := os.Remove(efsPath)
	check(err)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var messages string
	method := request.RequestContext.HTTPMethod

	switch method {
	case "GET":
		messages = getMessages()
	case "POST":
		addMessages(request.Body)
		messages = getMessages()
	case "DELETE":
		deleteMessages()
	default:
		messages = "Unsupported method"
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       messages,
	}, nil
}

func main() {
	lambda.Start(handler)
}
