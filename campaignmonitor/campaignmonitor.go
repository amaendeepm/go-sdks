package campaignmonitor-sdk


import (
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)


// Two parameters that need defining, or need to come from global context
// 1. CM_SignupListID : Subscriber List ID
// 2. CM_AuthHeader : Actual Authoriation Header for your Campaign Monitor Account


//PUBLIC: 1
//Refer: https://www.campaignmonitor.com/api/subscribers/#adding_a_subscriber

func addToSignUpList(jsonStr []byte, r *http.Request) (string, error) {
	return add2CMList(buildEndpointForSubscriberAdding(CM_SignupListID), jsonStr, r)
}

//PUBLIC: 2
//Refer: https://www.campaignmonitor.com/api/subscribers/#updating_a_subscriber

func updateSubscriber(jsonStr []byte, r *http.Request, email string) (string, error) {

	putEp:= buildEndpointForSubscriberUpdating(CM_SignupListID, email)

	//Start of Code - Required if working in context of AppEngine
	ctx := appengine.NewContext(r)
	client := &http.Client{
		Transport: &urlfetch.Transport{Context: ctx},
	}
	//End of Code - Required if working in context of AppEngine

	req, err := http.NewRequest("PUT", putEp, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Basic " + CM_AuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(ctx, "error-1: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(ctx, "error-2: %v", err)
		return "",err
	}

	log.Debugf(ctx,"Response Received: %d",string(body))

	js, err := simplejson.NewJson(body)

	if err != nil {
		log.Errorf(ctx,"error-3: %v", err)
		return "",err
	}

	respStr,_ := js.String()

	log.Debugf(ctx, "Response Body: %s", string(respStr))

	return respStr, nil
}

//PUBLIC : 3 : Sending Transaction Email
//Refer: https://www.campaignmonitor.com/api/transactional/

func sendCMTxnEmail(templateID string, jsonStr []byte, r *http.Request) (string, error) {

	return invokeCM_POST(buildEndpointForTxnEmail(templateID), jsonStr, r)
}


//INTERNAL USE METHODS

func buildEndpointForTxnEmail(templateID string) string {
	// Following this Spec: https://www.campaignmonitor.com/api/transactional/#send_a_smart_email
	return  "https://api.createsend.com/api/v3.1/transactional/smartemail/" + string(templateID)+ "/send"
}

func buildEndpointForSubscriberAdding(listID string) string { //
	// Following this Spec: https://www.campaignmonitor.com/api/subscribers/
	return "https://api.createsend.com/api/v3.1/subscribers/"+listID+".json"
}

func buildEndpointForSubscriberUpdating(listID string, email string) string { //
	// Following this Spec: https://www.campaignmonitor.com/api/subscribers/
	return "https://api.createsend.com/api/v3.1/subscribers/"+listID+".json?email="+email
}


func add2CMList(endpoint string, jsonStr []byte, r *http.Request) (string, error) {

	return invokeCM_POST(endpoint, jsonStr, r)
}


func invokeCM_POST(endpoint string, payload []byte, r *http.Request) (string, error) {

	//Start of Code - Required if working in context of AppEngine
	ctx := appengine.NewContext(r)
	client := &http.Client{
		Transport: &urlfetch.Transport{Context: ctx},
	}
	//End of Code - Required if working in context of AppEngine

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Basic " + CM_AuthHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(ctx, "error-1: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(ctx, "error-2: %v", err)
		return "",err
	}

	log.Debugf(ctx,"Response Received: %d",string(body))

	js, err := simplejson.NewJson(body)

	if err != nil {
		log.Errorf(ctx,"error-3: %v", err)
		return "",err
	}

	respStr,_ := js.String()

	log.Debugf(ctx, "Response Body: %s", string(respStr))

	return respStr, nil

}
