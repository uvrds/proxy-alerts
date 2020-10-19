package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type data struct {
	Url   string
	Token string
	Body  map[string]string
}

//Parsing response body
func parsingBodyCreateMessage(req *http.Request) string {
	//Read byte
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("err: %v\n", err)
	}

	//Parsing json Message
	var message Message
	err = json.Unmarshal(content, &message)
	if err != nil {
		log.Printf("warn: %v\n", err)
	}
	log.Printf("warn: %v\n", message)
	status := ""

	if message.Status == "firing" {
		status = "🔥"
	} else if message.Status == "resolved" {
		status = "✅"
	}
	severity := ""
	if message.Alerts[0].Labels.Severity == "critical" {
		severity = "❗"
	} else if message.Alerts[0].Labels.Severity == "warning" {
		severity = "⚠"
	} else if message.Alerts[0].Labels.Severity == "info" {
		severity = "ℹ"
	}

	Template := severity + status + "\n" +
		message.Alerts[0].Labels.AlertName +
		"\n---" +
		"\nВремя инцидента: " + message.Alerts[0].StartsAt.String() +
		"\nКластер:" + message.Alerts[0].Labels.ClusterName +
		"\nУзел: " + message.Alerts[0].Labels.Instance +
		"\nУровень инцидента: " + message.Alerts[0].Labels.Severity +
		"\n\nТекущие значение: " + message.Alerts[0].Annotations.CurrentValue +
		"\nПороговое значение: " + message.Alerts[0].Labels.ThresholdValue +
		"\nВыражение: " + message.Alerts[0].Labels.Expression

	return Template
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	res = nil
	bodyReq := parsingBodyCreateMessage(req)
	chatsId := []string{"246186171", "257434654"}
	for i, s := range chatsId {
		fmt.Println(i, s)

		options, err := json.Marshal(map[string]string{
			"chat_id": s,
			"text":    bodyReq,
		})
		data := data{
			Url:   os.Getenv("URL"),
			Token: os.Getenv("TOKEN"),
		}

		resp, err := http.Post(data.Url+"/"+data.Token+"/sendMessage", "Accept: application/json", bytes.NewBuffer(options))
		if err != nil {
			log.Printf("err: %v\n", err)
		}
		log.Printf("REPONSE: %v\n", resp)
		log.Printf("Template: %v\n", bodyReq)
	}
}

func main() {
	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
