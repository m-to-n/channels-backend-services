package main

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/m-to-n/channels-backend-services/data"
	"log"
	"net/http"
)

func foo() data.ChannelType {
	return data.Unknown
}

func cronHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Printf("cronHandler binding - Data:%s, Meta:%v", in.Data, in.Metadata)
	return nil, nil
}

func sqsHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Printf("sqsHandler binding - Data:%s, Meta:%v", in.Data, in.Metadata)
	return nil, nil
}

func main() {
	s := daprd.NewService(":6002")

	// cron binding is used for quick debugging / troubleshooting only
	/* if err := s.AddBindingInvocationHandler("/run", cronHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	} */

	if err := s.AddBindingInvocationHandler("/channels-backend-services-sqs-wa-twilio", sqsHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}
