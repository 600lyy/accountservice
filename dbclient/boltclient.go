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
	QueryAccount(accoundId string) (model.Account, error)
	CreateAccount(account *model.Account) (err error)
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
func (bc *BoltClient) QueryAccount(accoundID string) (account model.Account, err error) {

	err = bc.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		accountBytes := b.Get([]byte(accoundID))
		if accountBytes == nil {
			return fmt.Errorf("No account found for " + accoundID)
		}

		json.Unmarshal(accountBytes, &account)
		return nil
	})

	return  //bare return
}

func (bc *BoltClient) CreateAccount(account *model.Account) (err error) {
	var jsonByte []byte
	if jsonByte, err = json.Marshal(account); err !=nil {
		fmt.Errorf("cannnot marshal json: %v", err)
		return err
	}
	err = bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("AccountBucket"))
		err := b.Put([]byte(account.Id), jsonByte)
		return err
	})

	return nil
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
		_, err := tx.CreateBucket([]byte("AccountBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// Seed (n) make-believe account objects into the AcountBucket bucket.
func (bc *BoltClient) seedDemoAccounts() {

	total := 10
	for i := 0; i < total; i++ {

		// Generate a key 10000 or larger
		key := strconv.Itoa(10000 + i)

		// Create an instance of our Account struct
		acc := model.Account{
			Id:   	key,
			Name: 	"user_" + strconv.Itoa(i),
			Passwd:	"123456",
		}

		// Serialize the struct to JSON
		jsonBytes, _ := json.Marshal(acc)

		// Write the data to the AccountBucket
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AccountBucket"))
			err := b.Put([]byte(key), jsonBytes)
			return err
		})
	}
	fmt.Printf("Seeded %v demo accounts...\n", total)
}
