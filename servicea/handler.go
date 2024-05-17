package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	var req CEPRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil || len(req.CEP) != 8 {
		http.Error(w, `{"error":"invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	resp, err := http.Post("http://serviceb:8081/cep", "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, `{"error":"failed to connect to service B"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(respBody))
		http.Error(w, `{"error":"service B error"}`, resp.StatusCode)
		return
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read response from service B"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
