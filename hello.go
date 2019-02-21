package main

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"os"
	"strings"
)

var request string

func main() {
	// selenium_test.ExampleFindElement()
	// fmt.Println("hello world")

	// credential := models.UserCredentials{
	// 	Token:    "CLVqwgGaqhnI0Psp3BSo0J5GasiGQwlMo3kajdOk-jZ",
	// 	Email:    "szhu@freelancer.com",
	// 	Password: "",
	// }
	// client := rest.NewClient(&url.URL{Host: "chat.tools.flnltd.com:8080"}, true)
	// err := client.Login(&credential)
	// fmt.Println(err)

	const (
		// These paths will be different on your system.
		seleniumPath    = "/home/billz/go/src/github.com/tebeka/selenium/vendor/selenium-server-standalone-3.14.0.jar"
		geckoDriverPath = "/home/billz/go/src/github.com/tebeka/selenium/vendor/chromedriver-linux64-2.42"
		port            = 8888
	)
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.SetLogLevel("performance", "INFO")
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	wd.Get("https://chat.tools.flnltd.com")
	err = wd.Wait(LoggedIn)
	err = wd.Wait(IsHeaderSet)

	// Get Headers
	var data map[string]interface{}
	json.Unmarshal([]byte(request), &data)
	mesg, _ := data["message"].(map[string]interface{})
	params, _ := mesg["params"].(map[string]interface{})
	req, _ := params["request"].(map[string]interface{})
	headers, _ := req["headers"].(map[string]interface{})
	fmt.Println(headers)

	err = wd.Wait(Halt)
}

func LoggedIn(wd selenium.WebDriver) (bool, error) {
	_, err2 := wd.GetCookie("rc_token")
	return err2 == nil, nil
}

func Halt(wd selenium.WebDriver) (bool, error) {
	return false, nil
}

func IsHeaderSet(wd selenium.WebDriver) (bool, error) {
	msg, _ := wd.Log("performance")
	for i := 0; i < len(msg); i++ {
		if strings.Contains(strings.ToLower(msg[i].Message), "\"x-auth-token\"") {
			request = msg[i].Message
			return true, nil
		}
	}
	return false, nil
}
