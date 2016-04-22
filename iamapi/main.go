package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/smithatlanta/cloudthings/nonprovider/iamapi/internalAccount"
	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt/aws"
	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt/cloudhealth"
	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt/internalValidation"
)

//Users -
type Users struct {
	Users []User `json:"users"`
}

//User -
type User struct {
	Name   string  `json:"name"`
	Groups []Group `json:"groups"`
}

//Group -
type Group struct {
	GroupName string `json:"groupname"`
}

var intPort string
var intURL string
var chURL string
var intAuthKey string
var chAuthKey string
var region string
var role string

var users Users

func main() {
	flag.StringVar(&intPort, "intPort", "3000", "port for http server")
	flag.StringVar(&intURL, "intURL", "", "internal url for authorization")
	flag.StringVar(&chURL, "chURL", "", "cloudhealthurl for authorization")
	flag.StringVar(&intAuthKey, "intAuthKey", "", "internal authorization key")
	flag.StringVar(&chAuthKey, "chAuthKey", "", "cloud health authorization key")
	flag.StringVar(&region, "region", "", "aws region")
	flag.StringVar(&role, "role", "", "role to assume")
	flag.Parse()

	var missingStuff string

	if len(intURL) == 0 {
		missingStuff += "Missing -intURL (internal auth URL) parameter.\n"
	}
	if len(chURL) == 0 {
		missingStuff += "Missing -chURL (CloudHealth URL) parameter.\n"
	}

	if len(intAuthKey) == 0 {
		missingStuff += "Missing -intAuthKey (internal auth key) parameter.\n"
	}

	if len(chAuthKey) == 0 {
		missingStuff += "Missing -chAuthKey (CloudHealth auth key) parameter.\n"
	}

	if len(region) == 0 {
		missingStuff += "Missing -region (AWS region) parameter.\n"
	}

	if len(role) == 0 {
		missingStuff += "Missing -role (AWS role) parameter.\n"
	}

	if len(missingStuff) != 0 {
		fmt.Println(missingStuff)
		os.Exit(0)
	}

	fmt.Println("Listening on port: " + intPort)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("healthcheck", HealthCheckIndex).Methods("GET")
	router.HandleFunc("/chaccounts", CloudhealthAccountsIndex).Methods("GET")
	router.HandleFunc("/chaccount/{id}", CloudhealthAccountIndex).Methods("GET")
	router.HandleFunc("/awsaccount/all/{id}", AWSAllAccountIndex).Methods("GET")
	router.HandleFunc("/awsaccount/invalid/{id}", AWSInvalidAccountIndex).Methods("GET")
	router.HandleFunc("/intaccount/{id}", IntAccountGet).Methods("GET")
	router.HandleFunc("/intaccount", IntAccountPut).Methods("PUT")
	router.HandleFunc("/intaccount", IntAccountPost).Methods("POST")
	router.HandleFunc("/intaccount/{id}", IntAccountDelete).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+intPort, router))
}

// Index -
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nothing Here"))
}

// HealthCheckIndex - lb check
func HealthCheckIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// CloudhealthAccountsIndex -
func CloudhealthAccountsIndex(w http.ResponseWriter, r *http.Request) {
	account := account.NewSimpleAccount(chURL, chAuthKey)
	accounts, err := account.GetAccounts()
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusOK)
	} else {
		json.NewEncoder(w).Encode(accounts)
	}
}

// CloudhealthAccountIndex -
func CloudhealthAccountIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	account := account.NewSimpleAccount(chURL, chAuthKey)
	accounts, err := account.GetAccount(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(accounts)
}

// AWSAllAccountIndex -
func AWSAllAccountIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var accessKey *string
	var secretKey *string
	var sessionToken *string
	var err error

	//Assume Role in each account
	sts := aws.NewRole(region, nil, nil, nil)
	accessKey, secretKey, sessionToken, err = sts.AssumeRole("arn:aws:iam::"+id+":role/"+role, "IAMCHECKER")
	if err != nil {
		fmt.Println(err.Error())
		var users []string
		json.NewEncoder(w).Encode(users)
	} else {
		//AWS user scan showing name, groups
		iam := aws.NewIdentity(region, accessKey, secretKey, sessionToken)
		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
			var emptyusers []string
			json.NewEncoder(w).Encode(emptyusers)
		} else {
			json.NewEncoder(w).Encode(users)
		}
	}
}

// AWSInvalidAccountIndex -
func AWSInvalidAccountIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var userList []string

	var accessKey *string
	var secretKey *string
	var sessionToken *string
	var err error

	//Assume Role in each account
	sts := aws.NewRole(region, nil, nil, nil)
	accessKey, secretKey, sessionToken, err = sts.AssumeRole("arn:aws:iam::"+id+":role/"+role, "IAMCHECKER")
	if err != nil {
		fmt.Println(err.Error())
		var users []string
		json.NewEncoder(w).Encode(users)
	} else {
		//AWS user scan showing name, groups
		iam := aws.NewIdentity(region, accessKey, secretKey, sessionToken)
		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
			var emptyusers []string
			json.NewEncoder(w).Encode(emptyusers)
		}

		for a := 0; a < len(users); a++ {
			internal := internalValidation.NewSimpleIdentity(intURL, intAuthKey)
			isValid, err := internal.Validate(users[a])
			if err != nil && err.Error() != "No records found" {
				continue
			}
			if isValid == false {
				userList = append(userList, users[a])
			}
		}
		json.NewEncoder(w).Encode(userList)
	}
}

// IntAccountGet -
func IntAccountGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	store := internalAccount.NewAccountStore(region, true)
	acc, err := store.GetByAccountID(id)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNoContent)
	}

	json.NewEncoder(w).Encode(acc)
}

// IntAccountPut -
func IntAccountPut(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t internalAccount.Account
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	store := internalAccount.NewAccountStore(region, true)
	err = store.Add(t)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// IntAccountPost -
func IntAccountPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t internalAccount.Account
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	store := internalAccount.NewAccountStore(region, true)
	err = store.Add(t)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// IntAccountDelete -
func IntAccountDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	store := internalAccount.NewAccountStore(region, true)
	err := store.Delete(id)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
