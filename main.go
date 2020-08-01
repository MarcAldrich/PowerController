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
		return
	}

	rawUserRequestedPumpId := r.URL.Query().Get("pumpRelayId")
	// Get userRequestedPumpID ID to cycle
	userRequestedPumpID, err := strconv.Atoi(rawUserRequestedPumpId)
	if err != nil {
		http.Error(w, fmt.Sprintf("pumpRelayId: %s not recognized.", rawUserRequestedPumpId), http.StatusBadRequest)
		return
	}

	if userRequestedPumpID < len(pumpPowerRelayControlPins)-1 {
		switch hwState := pumpPowerRelayControlPins[userRequestedPumpID].Read(); hwState {
		case rpio.High:
			pumpPowerRelayControlPins[userRequestedPumpID].Low()
		case rpio.Low:
			pumpPowerRelayControlPins[userRequestedPumpID].High()
		}
	}
	// Request completed
	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	// Initialize Gpio
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pumpPowerRelayControlPin to output mode
	for _, pumpRelayControlPin := range pumpPowerRelayControlPins {
		log.Printf(fmt.Sprintf("Setting GPIO Pin %d"))
		pumpRelayControlPin.Output()
	}

	// Start REST server
	handleRequests()
}