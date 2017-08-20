package atshook

import (
	"github.com/sirupsen/logrus"
	"github.com/Azure/azure-sdk-for-go/storage"
	"os"
	"strconv"
	"errors"
)

const (
	// TableAlreadyExists indicates table already exists in Azure.
	tableAlreadyExists string = "TableAlreadyExists"
	timestampID string = "LogTimestamp"
	levelID string = "Level"
	messageID string = "Message"
)

// AtsHook to handle writing to Azure Table Storage
type AtsHook struct {

	// Azure specifics
	accountName string
	accountKey string
	tableName string
	
	// azure table client
	tableCli storage.TableServiceClient
	table *storage.Table
	
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
		//fmt.Printf("Unable to create client for Azure Table Storage hook %s", err)
		return nil // is nil valid?
	}
	
	hook.tableCli = client.GetTableService()
	table, err := createTable(hook.tableCli, tableName)
	if err != nil {

		// cant log...   but return no hook!
		return nil
	}

	hook.table = table
	hook.accountName = accountName
	hook.accountKey = accountKey
	hook.tableName = tableName
	hook.levels = levels
	return hook
}

func createTable( tableCli storage.TableServiceClient , tableName string) (*storage.Table, error) {
	table := tableCli.GetTableReference(tableName)
	err := table.Create( 30, storage.EmptyPayload, nil )
	if err != nil {
		azureErr, ok := err.(storage.AzureStorageServiceError)
		if !ok {
			// error... what to do?  Cant log it can we?
			return nil, err
		} 
		
		if azureErr.Code != tableAlreadyExists {
			// we are ok if the table already exists. Otherwise return nil
			return nil, errors.New("Unable to create log table")
		}
	
	}

	return table, nil
}

func createTableClient( accountName string, accountKey string ) (*storage.Client, error) {
	// use parameters if passed in otherwise fall back to env vars.
	if accountName == "" || accountKey == "" {

		accountName = os.Getenv("ACCOUNT_NAME")
		accountKey = os.Getenv("ACCOUNT_KEY")
	}	
	client, err := storage.NewBasicClient(accountName, accountKey)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// Fire adds the logrus entry to Azure Table Storage
func (hook *AtsHook) Fire(entry *logrus.Entry) error {

	rowKey := strconv.FormatInt(int64(entry.Time.UnixNano()), 10)
	tableEntry := hook.table.GetEntityReference("logrus",rowKey )
	props := make(map[string]interface{})

	// technically dont need to make since entry.Data is already a map to interface. But will keep mapping here incase it changes.
	for k,v := range entry.Data {
		props[k] = v
	}
	props[timestampID] = entry.Time.UTC()
	props[levelID] = entry.Level.String()
	props[messageID] = entry.Message	
	tableEntry.Properties = props
	err := tableEntry.Insert(storage.EmptyPayload, nil)
	if err != nil {
		return err
	}
	
	return nil
}


// Levels returns configured log levels
func (hook *AtsHook) Levels() []logrus.Level {
	return hook.levels
}
