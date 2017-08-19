package atshook

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/Azure/azure-sdk-for-go/storage"
	"os"
)

// AtsHook to handle writing to Azure Table Storage
type AtsHook struct {

	// Azure specifics
	accountName string
	accountKey string
	tableName string
	
	// azure table client
	tableCli storage.TableServiceClient

	levels    []logrus.Level
	formatter logrus.Formatter

}

// NewHook creates a new instance of atsHook.
// The accountName, accountKey and tableName for Azure are required.
func NewHook(accountName string, accountKey string, tableName string, level logrus.Level) *AtsHook {
	levels := []logrus.Level{}
	for _, lev := range logrus.AllLevels {
		if lev <= level {
			levels = append(levels, lev)
		}
	}
	
	
	hook := &AtsHook{}	
	client, err  := createTableClient(accountName, accountKey)
	if err != nil {
		// unsure what to do with error.....?
		fmt.Printf("Unable to create Azure Table Storage hook %s", err)
		return nil // is nil valid?
	}
	
	hook.tableCli = client.GetTableService()
	hook.accountName = accountName
	hook.accountKey = accountKey
	hook.tableName = tableName
	hook.levels = levels
	return hook
}

func createTableClient( accountName string, accountKey string ) (*storage.Client, error) {
	// use parameters if passed in otherwise fall back to env vars.
	if accountName == "" || accountKey == "" {

		accountName = os.Getenv("ACCOUNT_NAME")
		accountKey = os.Getenv("ACCOUNT_KEY")
		
		client, err := storage.NewBasicClient(accountName, accountKey)
		if err != nil {
			return nil, err
		}

		return &client, nil
	}

	return nil, errors.New("Unable to create Azure Table Storage clent")
}

func (hook *AtsHook) Fire(entry *logrus.Entry) error {

	return nil
}


// Levels returns configured log levels
func (hook *AtsHook) Levels() []logrus.Level {
	return hook.levels
}
