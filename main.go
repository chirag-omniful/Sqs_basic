package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	region := "us-west-2"
	queue, err := CreateQueue("meri_kyu", region)
	if err != nil {
		return
	}

	fmt.Println("Queue Created:", *queue.QueueUrl)

	SendMessageOutput, err := SendMessage("qwerty", queue)
	if err != nil {
		fmt.Println("message bhejne me dikkat hai", err)
		return
	}

	fmt.Println("Message Sent", SendMessageOutput)
	result, err := RecieveMessage(*queue.QueueUrl)
	fmt.Println(result[0].Body)
}

var sess = session.Must(session.NewSession())

func CreateQueue(queueName string, region string) (*sqs.CreateQueueOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	})
	svc := sqs.New(sess)

	params := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}

	result, err := svc.CreateQueue(params)
	if err != nil {
		fmt.Println("queue creation me dikkat hai")
		return nil, err
	}

	return result, nil
}

func SendMessage(message string, queue *sqs.CreateQueueOutput) (*sqs.SendMessageOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	})
	if err != nil {
		return nil, err
	}

	svc := sqs.New(sess)
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(message),
		QueueUrl:     queue.QueueUrl,
		DelaySeconds: aws.Int64(0),
	}

	return svc.SendMessage(params)
}

func RecieveMessage(queueUrl string) ([]*sqs.Message, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	})
	if err != nil {
		return nil, err
	}

	svc := sqs.New(sess)
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(60),
	}

	result, err := svc.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}

	return result.Messages, nil

}
