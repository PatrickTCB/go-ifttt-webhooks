package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type arguments struct {
	val1  string
	val2  string
	val3  string
	event string
	key   string
}

func args(rawArgs []string) arguments {
	var a arguments
	for _, raw := range rawArgs {
		if strings.Contains(raw, "val1=") {
			a.val1 = strings.ReplaceAll(raw, "val1=", "")
		} else if strings.Contains(raw, "val2=") {
			a.val2 = strings.ReplaceAll(raw, "val2=", "")
		} else if strings.Contains(raw, "val3=") {
			a.val3 = strings.ReplaceAll(raw, "val3=", "")
		} else if strings.Contains(raw, "key=") {
			a.key = strings.ReplaceAll(raw, "key=", "")
		} else if strings.Contains(raw, "event=") {
			a.event = strings.ReplaceAll(raw, "event=", "")
		} else {
			if !strings.Contains(raw, "ifttt-webhook") {
				fmt.Printf("%s isn't a valid argument. This program accepts val1, val2, val3, key & event as valid arguments.\n", raw)
			}
		}
	}
	if a.key == "" && os.Getenv("IFTTT_WEBHOOK_KEY") != "" {
		a.key = os.Getenv("IFTTT_WEBHOOK_KEY")
	}
	if a.event == "" && os.Getenv("IFTTT_WEBHOOK_DEFAULT_EVENT") != "" {
		a.event = os.Getenv("IFTTT_WEBHOOK_DEFAULT_EVENT")
	}
	return a
}

func main() {
	a := args(os.Args)
	if a.key == "" {
		fmt.Println("Oops. You need to specify a key either as an argument 'key=foo' or environment variable IFTTT_WEBHOOK_KEY")
		os.Exit(1)
	} else if a.event == "" {
		fmt.Println("Oops. You need to specify an event. This can be done either as an argument 'event=bar' or environment variable IFTTT_WEBHOOK_DEFAULT_EVENT")
		os.Exit(2)
	} else {
		ifttt_url := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", a.event, a.key)
		ifttt_params := url.Values{}
		if a.val1 != "" {
			ifttt_params.Add("value1", a.val1)
		}
		if a.val2 != "" {
			ifttt_params.Add("value2", a.val2)
		}
		if a.val3 != "" {
			ifttt_params.Add("value3", a.val3)
		}
		response, rerr := http.PostForm(ifttt_url, ifttt_params)
		if rerr != nil {
			fmt.Println(rerr.Error())
			os.Exit(3)
		}
		fmt.Println("Response received. Status " + strconv.Itoa(response.StatusCode))
		defer response.Body.Close()
		body, berr := ioutil.ReadAll(response.Body)
		if berr != nil {
			fmt.Println(berr.Error())
			os.Exit(4)
		}
		fmt.Printf("Response content: %s\n", string(body))
	}
}
