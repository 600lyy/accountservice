package service

import (
	"log"
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"html/template"
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
	var username = mux.Vars(r)["username"]
	account, err := DBClient.QueryAccount(username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// account.ServedBy = getIP()

	// If found, marshal into JSON, write headers and content
	account.Passwd = ""
	data, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetAllDemoAccounts returns an account
func GetAllDemoAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := DBClient.QueryAllDemoAccounts()
	size := len(accounts)
	if size == 0 {
		log.Println("Account database is empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//datas := make([]byte, size, size)

	datas, _ := json.MarshalIndent(accounts, "", "    ")
		//datas = append(datas, data...)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(datas)))
	w.WriteHeader(http.StatusOK)
	w.Write(datas)
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

	if r.Method != http.MethodPost {
		t := template.Must(template.ParseFiles(
			"/home/autotest/go/src/github.com/600lyy/accountservice/templates/html/login/register.html"))
		t.Execute(w, nil)
		return
	}

	account := model.Account{}
	body, err := ioutil.ReadAll(&io.LimitedReader{r.Body, 10485760})
	
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &account)
	if err != nil {
		log.Println(err)
	}

	if err = DBClient.CreateAccount(&account); err != nil {
		log.Printf("cannnot create account [%s]: %v", account.Name, err)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len("User Already Exists")))
		w.Write([]byte(err.Error()))
		return
	}
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

//UserLogin handles user login
//Method GET to receive the form of /login/home.html 
func UserLogin(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.Method)
	var err error
	if r.Method != http.MethodPost {
		t := template.Must(template.ParseFiles(
			"/home/autotest/go/src/github.com/600lyy/accountservice/templates/html/login/home.html"))
		t.Execute(w, nil)
		return
	}

	r.ParseForm()
	log.Println(r.Form)
	username := r.Form["username"][0]
	
	if username != "" {
		account, _ := DBClient.QueryAccount(username)
		password := r.Form["password"][0]
		if password == account.Passwd {
			err = nil
		} else {
			err = fmt.Errorf("authentication failed, check user password")
		}
	}
	
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(err.Error())))
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login/index", http.StatusFound)
		log.Println("Redirect request to /login")
		return
	}
	t, err := template.ParseFiles("/home/autotest/go/src/github.com/600lyy/accountservice/templates/html/404.html")
	if (err != nil) {
		log.Println(err)
	}
	t.Execute(w, nil)
}
