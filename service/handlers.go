package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/600lyy/accountservice/dbclient"
	"github.com/600lyy/accountservice/model"
)

// DBClient is a global variable
var DBClient dbclient.IBoltClient
var RedisClient dbclient.IRedisClient

// GetAccount returns an account
func GetAccount(w http.ResponseWriter, r *http.Request) {
	var accountId = mux.Vars(r)["accountId"]
	account, err := DBClient.QueryAccount(accountId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// account.ServedBy = getIP()

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// HealthCheck is the handlerfunc for endpoint "check"
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Since we're here, we already know that HTTP service is up. Let's just check the state of the boltdb connection
	dbUp := DBClient.Check()
	if dbUp {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJSONResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccessible"})
		writeJSONResponse(w, http.StatusServiceUnavailable, data)
	}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := model.Account{}
	body, err := ioutil.ReadAll(&io.LimitedReader{r.Body, 10485760})
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	if err = json.Unmarshal(body, &account); err != nil {
		fmt.Errorf("cannnot Unmarshal json: %v", err)
	}
	DBClient.CreateAccount(&account)
}

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
