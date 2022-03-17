package main

import (
	"consulwatcherpoc/config"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func main() {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	// Lookup the pair
	pair, _, err := kv.Get("bigbluebutton/secret", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)

	config.Load(string(pair.Value))

	params := map[string]interface{}{
		"type": "key",
		"key":  "bigbluebutton/secret",
	}

	plan, err := watch.Parse(params)
	if err != nil {
		panic(err)
	}

	plan.Handler = func(u uint64, raw interface{}) {
		fmt.Println("In handler")
		var pair *api.KVPair
		if raw == nil {
			pair = nil
		} else {
			var ok bool
			if pair, ok = raw.(*api.KVPair); !ok {
				return // ignore
			}
		}

		if pair == nil {
			panic(errors.New("pair is nil"))
		} else {
			config.Load(string(pair.Value))
			fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
		}
	}

	go func() {
		fmt.Println("in go routine")
		if err := plan.Run(api.DefaultConfig().Address); err != nil {
			panic(fmt.Errorf("err: %v", err))
		}
	}()

	for true { // A simple infinite loop to keep process running
		time.Sleep(1 * time.Second)
		fmt.Println(config.Bbb.Secret)
	}
}
