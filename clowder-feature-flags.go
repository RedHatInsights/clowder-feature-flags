package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/Unleash/unleash-client-go/v3/context"
	client "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Args: %v", os.Args)
	fmt.Fprintln(w, msg)
}

func ffHandler(w http.ResponseWriter, r *http.Request) {
	features := make(map[string]bool)
	ctx := context.Context{
		UserId: r.Header.Get("UserId"),
	}
	for _, flag := range []string{"normalOption", "adminOption"} {
		features[flag] = unleash.IsEnabled(flag, unleash.WithContext(ctx))
		fmt.Fprintf(w, "%v:%v\n", flag, unleash.IsEnabled(flag, unleash.WithContext(ctx)))
	}
}

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Hi")
	} else {
		port := client.LoadedConfig.WebPort

		unleash.Initialize(
			unleash.WithListener(&unleash.DebugListener{}),
			unleash.WithAppName("clowder-hello"),
			unleash.WithUrl(fmt.Sprintf("http://%v:%v/api", client.LoadedConfig.FeatureFlags.Hostname, client.LoadedConfig.FeatureFlags.Port)),
		)

		http.HandleFunc("/", ffHandler)
		http.HandleFunc("/healthz", helloHandler)

		address := fmt.Sprintf(":%d", port)

		fmt.Printf("Started, serving at %s\n", address)
		err := http.ListenAndServe(address, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}
}
