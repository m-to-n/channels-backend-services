package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/m-to-n/channels-backend-services/dapr"
	whatsapp "github.com/m-to-n/common/channels/whatsapp-twilio"
	common_dapr "github.com/m-to-n/common/dapr"
	"github.com/m-to-n/common/logging"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// quick & dirty for initial test only! do it properly! :)
// https://www.twilio.com/docs/whatsapp/tutorial
func sendTwilioResponse(request whatsapp.TwilioRequest, response string, accSid string, authToken string) (*string, error) {
	client := &http.Client{}
	twilioUrl := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json?From", accSid)
	v := url.Values{}
	v.Set("From", request.To)
	v.Set("To", request.From)
	v.Set("Body", response)

	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(accSid, authToken)
	// req.Header.Add("Authorization", "Basic ...")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("sending twilio request: %s: ", req)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := string(bodyText)
	return &s, nil

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

	// create the client
	client := common_dapr.DaprClient(dapr.DAPR_GRPC_PORT)
	if client == nil {
		log.Printf("dapr client init error")
		return nil, err
	}

	// TBD: this will be taken from tenant config!
	opt := map[string]string{}
	secretAccId, err := client.GetSecret(ctx, "channels-backend-services-secret-store", "twilioAccSid", opt)
	secretAuthToken, err := client.GetSecret(ctx, "channels-backend-services-secret-store", "twilioAuthToken", opt)

	log.Println(secretAccId)
	log.Println(secretAuthToken)

	resp, err := sendTwilioResponse(tReq, fmt.Sprintf("you said: %s", tReq.Body), secretAccId["twilioAccSid"], secretAuthToken["twilioAuthToken"])

	if err != nil {
		log.Println("unable to send twilio response")
		log.Println(err)
		return nil, err
	}

	log.Println("twilio response sent %s", *resp)
	return nil, nil
}

func main() {
	s := common_dapr.DaprService(dapr.DAPR_APP_GRPC_ADDR)

	// cron binding is used for quick debugging / troubleshooting only
	/* if err := s.AddBindingInvocationHandler("/run", cronHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	} */

	if err := s.AddBindingInvocationHandler(dapr.DAPR_BINDING_SQS_GRPC, sqsHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}
