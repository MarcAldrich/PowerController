package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"net/http"
	"os"
	"strconv"
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
	userRequestedPumpID, ok := r.URL.Query()["pumpRelayId"]
	if ok {
		//If not found error out
		http.Error(w, "pumpRelayId not found in option strings. Example: `curl -X POST raspberrypi3-64.local:2080/pump?pumpRelayId[]=1`.", http.StatusBadRequest)
		return
	}

	// Loop over all values passed in by user
	for _, userInputPumpRelayId := range userRequestedPumpID {
		// Convert pumpID from string to int
		pumpRelayId, err := strconv.Atoi(userInputPumpRelayId) //TODO: Make this line not hardcoded. Handle list decomposition, maybe the user wants to control multiple pumps at once.
		if err != nil {
			http.Error(w, fmt.Sprintf("Expected Relay ID. Must specify which pumpRelayId to toggle. `%s` not recognized.", pumpRelayId), http.StatusBadRequest)
			return
		}

		// Verify pumpID is in range
		if pumpRelayId > len(pumpPowerRelayControlPins) {
			http.Error(w, "Pump Relay ID unknown.", http.StatusBadRequest)
		}

		switch hwState := pumpPowerRelayControlPins[pumpRelayId].Read(); hwState {
		case rpio.High:
			pumpPowerRelayControlPins[pumpRelayId].Low()
		case rpio.Low:
			pumpPowerRelayControlPins[pumpRelayId].High()
		}
	}

	w.WriteHeader(http.StatusOK)
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