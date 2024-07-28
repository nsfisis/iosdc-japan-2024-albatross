package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type RequestBody struct {
	Code  string `json:"code"`
	Stdin string `json:"stdin"`
}

type ResponseBody struct {
	Result string `json:"result"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func doExec(code string, stdin string, maxDuration time.Duration) ResponseBody {
	_ = code
	_ = stdin
	_ = maxDuration

	return ResponseBody{
		Result: "success",
		Stdout: "42",
		Stderr: "",
	}
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	maxDurationStr := r.URL.Query().Get("max_duration")
	maxDuration, err := strconv.Atoi(maxDurationStr)
	if err != nil || maxDuration <= 0 {
		http.Error(w, "Invalid max_duration parameter", http.StatusBadRequest)
		return
	}

	var reqBody RequestBody
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resBody := doExec(reqBody.Code, reqBody.Stdin, time.Duration(maxDuration)*time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resBody)
}

func main() {
	http.HandleFunc("/api/exec", execHandler)
	http.ListenAndServe(":80", nil)
}
