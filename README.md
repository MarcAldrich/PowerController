# PowerController
Automates electrical loading of battery storage system in PowerHive. Written in a microservices architecture to expirement with implementations in low-power/edge environments.

# Project Status
Early proof-of-concept.

# Example Usage
Deploy via docker-compose on RPi.
`curl -X POST raspberrypi3-64.local:2080/pump?pumpRelayId=1` <-- Inverts the state of the relay on relay 1
- NOTE: PumpIds are current hardcoded and 0-indexed. This should be moved to enviornment-var based config in near future.

# Hardware
Assumes Raspberry Pi 3, as that's what I have. Should easily support other RPi generations.
Relay Board

