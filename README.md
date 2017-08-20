# Azure Table Storage Hooks for [Logrus](https://github.com/Sirupsen/logrus) 

## Usage

```go

package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/kpfaulkner/azuretablehook"
	)

func main() {

    // Azure Storage keys can either be provided in the NewHook method or it will be read from environment variables
    // ACCOUNT_NAME and ACCOUNT_KEY
    log.AddHook( atshook.NewHook("XXXXX Azure account name XXXX", "XXXX Azure account key XXXXX", "mylogtable", 
    log.DebugLevel) )
    log.SetLevel( log.DebugLevel)

    log.WithFields(log.Fields{
        "species": "cat",
        "name" :"fred",
        "number": 1,
    }).Info("A cat was here")
}

```