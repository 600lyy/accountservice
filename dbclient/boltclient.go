package dbclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/600lyy/accountservice/model"
	"github.com/boltdb/bolt"
)

//IBoltClient interface
type IBoltClient interface {
	OpenBoltDb()
	QueryAccount(username string) (model.Account, error)
	CreateAccount(account *model.Account) (err error)
	QueryAllDemoAccounts() (accounts []model.Account)
	Seed()
	Check() bool
}

//BoltClient uses *bolt.DB
type BoltClient struct {
	boltDB *bolt.DB
}

// OpenBoltDb opens and connnect to a db
func (bc *BoltClient) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open("account.db", 0600, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// QueryAccount returns a account indexed by id
func (bc *BoltClient) QueryAccount(username string) (account model.Account, err error) {

	err = bc.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		accountBytes := b.Get([]byte(username))
		if accountBytes == nil {
			return fmt.Errorf("No username found for " + username)
		}

		json.Unmarshal(accountBytes, &account)
		return nil
	})

	return  //bare return
}

//QueryAllDemoAccounts iterats over key/value pairs in a bucket
func (bc *BoltClient) QueryAllDemoAccounts() (accounts []model.Account) {
	bc.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		b.ForEach(func(k, v []byte) error {
			account := model.Account{}
			json.Unmarshal(v, &account)
			accounts = append(accounts, account)
			return nil
		})
		
		return nil
	})
	return accounts
}

func (bc *BoltClient) CreateAccount(account *model.Account) (err error) {

	err = bc.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		accountBytes := b.Get([]byte(account.UserName))
		if accountBytes == nil {
			return fmt.Errorf("No username found for " + account.UserName)
		}
		return nil
	})

	// return if user already exists
	if err == nil {
		return fmt.Errorf("User already exists: " + account.UserName)
	}

	err = bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		id, _ := b.NextSequence()
		account.ID = (id)
		jsonBytes, _ := json.Marshal(account)
		err := b.Put([]byte(account.UserName), jsonBytes)
		return err
	})

	return err
}

// Seed starts seeding accounts
func (bc *BoltClient) Seed() {
	bc.initializeBucket()
	bc.seedDemoAccounts()
}

// Check db connection
func (bc *BoltClient) Check() bool {
	return bc.boltDB != nil
}

func (bc *BoltClient) initializeBucket() error {
	return bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("AccountBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// Seed (n) make-believe account objects into the AcountBucket bucket.
func (bc *BoltClient) seedDemoAccounts() {

	total := 10
	for i := 1; i <= total; i++ {

		acc := model.Account {
			UserName: 	"user_" + strconv.Itoa(i),
			Name:		"DemoUser",
			Passwd:		"123456",
		}

		// Write the data to the AccountBucket
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AccountBucket"))
			acc.ID = 0
			jsonBytes, _ := json.Marshal(acc)
			fmt.Printf("Account: %s\n", jsonBytes)
			err := b.Put([]byte(acc.UserName), jsonBytes)
			return err
		})
	}
	fmt.Printf("Seeded %v demo accounts...\n", total)
}
