package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Vehicle Declaration */
type Vehicle struct {
	Make          string `json:"make"`
	Model         string `json:"model"`
	Owner         string `json:"owner"`
	ChasisNumber  int  `json:"chasisnumber"`
        EngineNumber  int  `json:"enginenumber"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getVehicleDetails" {
		return s.getVehicleDetails(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createVehicle" {
		return s.createVehicle(APIstub, args)
	} else if function == "getAllVehicles" {
		return s.getAllVehicles(APIstub)
	} else if function == "changeOwnerShip" {
		return s.changeOwnerShip(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getVehicleDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	vehicleAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(vehicleAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	vehicles := []Vehicle{
		Vehicle{Make: "Maruthi", Model: "Swift",   Owner: "MFR", ChasisNumber: 1234567, EngineNumber: 9866053729},
		Vehicle{Make: "Maruthi", Model: "Desire",  Owner: "MFR", ChasisNumber: 2345678, EngineNumber: 1111111111},
		Vehicle{Make: "Maruthi", Model: "Santro",  Owner: "MFR", ChasisNumber: 3456787, EngineNumber: 2222222229},
		Vehicle{Make: "Maruthi", Model: "Celerio", Owner: "MFR", ChasisNumber: 3456767, EngineNumber: 3333333339},
	}

	i := 0
	for i < len(vehicles) {
		fmt.Println("i is ", i)
		vehicleAsBytes, _ := json.Marshal(vehicles[i])
		APIstub.PutState("VEHICLE"+strconv.Itoa(i), vehicleAsBytes)
		fmt.Println("Added", vehicles[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createVehicle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

        chasisnumber,_ := strconv.Atoi(args[3]);
        enginenumber,_ := strconv.Atoi(args[4]);
	vehicle := Vehicle{Make: args[0], Model: args[1], Owner: args[2], ChasisNumber: chasisnumber, EngineNumber: enginenumber}

	vehicleAsBytes, _ := json.Marshal(vehicle)
	APIstub.PutState(args[3], vehicleAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getAllVehicles(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "1111111"
	endKey := "9999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getAllVehicles:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeOwnerShip(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	vehicleAsBytes, _ := APIstub.GetState(args[0])
	vehicle := Vehicle{}

	json.Unmarshal(vehicleAsBytes, &vehicle)
	vehicle.Owner = args[1]

	vehicleAsBytes, _ = json.Marshal(vehicle)
	APIstub.PutState(args[0], vehicleAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
