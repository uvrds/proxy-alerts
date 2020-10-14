package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			AlertName      string `json:"alert_name"`
			AlertType      string `json:"alert_type"`
			Alertname      string `json:"alertname"`
			ClusterName    string `json:"cluster_name"`
			Comparison     string `json:"comparison"`
			Duration       string `json:"duration"`
			Expression     string `json:"expression"`
			GroupID        string `json:"group_id"`
			Instance       string `json:"instance"`
			Prometheus     string `json:"prometheus"`
			PrometheusFrom string `json:"prometheus_from"`
			RuleID         string `json:"rule_id"`
			Severity       string `json:"severity"`
			ThresholdValue string `json:"threshold_value"`
		} `json:"labels"`
		Annotations struct {
			CurrentValue string `json:"current_value"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		RuleID string `json:"rule_id"`
	} `json:"groupLabels"`
	CommonLabels struct {
		AlertName      string `json:"alert_name"`
		AlertType      string `json:"alert_type"`
		Alertname      string `json:"alertname"`
		ClusterName    string `json:"cluster_name"`
		Comparison     string `json:"comparison"`
		Duration       string `json:"duration"`
		Expression     string `json:"expression"`
		GroupID        string `json:"group_id"`
		Instance       string `json:"instance"`
		Prometheus     string `json:"prometheus"`
		PrometheusFrom string `json:"prometheus_from"`
		RuleID         string `json:"rule_id"`
		Severity       string `json:"severity"`
		ThresholdValue string `json:"threshold_value"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		CurrentValue string `json:"current_value"`
	} `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
	Version     string `json:"version"`
	GroupKey    string `json:"groupKey"`
}

type TestMessage []struct {
	Labels struct {
		TestMsg string `json:"test_msg"`
	} `json:"labels"`
	Annotations  interface{} `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
}

type data struct {
	Url   string
	Token string
	Body  string
}

//Parsing response body
func parsingBody(req *http.Request) string {

	//Read byte
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("err: %v\n", err)
	}

	//Parsing json test Message
	var respRancher TestMessage
	err = json.Unmarshal(content, &respRancher)
	if err != nil {
		log.Printf("warn: %v\n", err)
	}
	return respRancher[0].Labels.TestMsg
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	body := parsingBody(req)

	data := data{
		Url:   os.Getenv("URL"),
		Token: os.Getenv("TOKEN"),
		Body:  body,
	}

	resp, err := http.Get(data.Url + "/" + data.Token + "/sendMessage?" + "chat_id=246186171&parse_mode=markdown&text=" + data.Body)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	log.Printf("res: %v\n", resp)
	log.Printf("BODY: %v\n", data.Body)
}

func main() {
	// start server
	http.HandleFunc("/", handleRequestAndRedirect)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
