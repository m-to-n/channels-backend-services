package main

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/m-to-n/common/channels"
	whatsapp "github.com/m-to-n/common/channels/whatsapp-twilio"
	"github.com/m-to-n/common/logging"
	"log"
	"net/http"
)

func foo() channels.ChannelType {
	return channels.Unknown
}

func cronHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Printf("cronHandler binding - Data:%s, Meta:%v", in.Data, in.Metadata)
	return nil, nil
}

func sqsHandler(ctx context.Context, in *common.BindingEvent) ([]byte, error) {
	log.Printf("sqsHandler binding - Data:%s, Meta:%v", in.Data, in.Metadata)

	var tReq whatsapp.TwilioRequest

	err := json.Unmarshal(in.Data, &tReq)
	if err != nil {
		log.Printf("sqsHandler: error when unamrshaling sqs payload: #{err}")
		return nil, err
	}

	// move struct pretty print to common
	var structStr *string
	structStr, err = logging.StructToPrettyString(tReq)
	if err != nil {
		return nil, err
	}
	log.Printf("twilio request: %s: ", *structStr)

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
