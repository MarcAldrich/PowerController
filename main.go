package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"net/http"
	"os"
)

var (

	pumpPowerRelayControlPins = [2]rpio.Pin{rpio.Pin(14), rpio.Pin(10)} //TODO: Give pumps a struct. They should be searchable. {Name, hardware pin, unique id, power relay status} UID -> make them globally unique from the start
	httpAddrAndListenPort = ":2080"
)

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "PowerHive by RedBeardRanch\n")
	fmt.Println("Endpoint Hit: home")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/pump", handlePumpRequest)
	log.Fatal(http.ListenAndServe(httpAddrAndListenPort, nil))
}

// Alternates pump status on/off.
func handlePumpRequest(w http.ResponseWriter, r *http.Request) {
	// Only handling a post method
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotImplemented)
	}

	// Get userRequestedPumpID ID to cycle
	userRequestedPumpID := r.URL.Query().Get("pumpRelayId")

	// Loop over all values passed in by user
	// TODO: START HERE. Don't think the loop is running at all
	for _, userInputPumpRelayId := range userRequestedPumpID {
		//Make sure we have that pin in our config
		if int(userInputPumpRelayId) < len(pumpPowerRelayControlPins) {

			switch hwState := pumpPowerRelayControlPins[userInputPumpRelayId].Read(); hwState {
			case rpio.High:
				pumpPowerRelayControlPins[userInputPumpRelayId].Low()
			case rpio.Low:
				pumpPowerRelayControlPins[userInputPumpRelayId].High()
			}
			// Request completed
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Failed to set pin `%d` because it isn't configured on this device.", http.StatusBadRequest)
		}
	}
	return
}

func setupGpio() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pumpPowerRelayControlPin to output mode
	for _, controlPin := range pumpPowerRelayControlPins {
		controlPin.Output()
	}
}

func main() {
	// Initialize Gpio
	setupGpio()

	// Start REST server
	handleRequests()
}