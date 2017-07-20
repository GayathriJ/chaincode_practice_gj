package main

import (
	"errors"
	"fmt"
    "time"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// KYCChainCode example simple Chaincode implementation
type KYCChainCode struct {
}

type KYCData struct{
    GCI string `json:"gci"`
    NAME string `json:"name"`
    UPDATED_TIME string `json:"updated_time"`
}

func main() {
	err := shim.Start(new(KYCChainCode))
	if err != nil {
		fmt.Printf("Error starting KYCChainCode: %s", err)
	}
}

// Init resets all the things
func (t *KYCChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	value,_ := json.Marshal(true)
    err := stub.PutState("initStatus", value)
    if err != nil {
    	return nil, errors.New("Failed to Initialize chaincode")
    }
    return nil, err
}

// Invoke isur entry point to invoke a chaincode function
func (t *KYCChainCode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "set" {
		return t.set(stub, args)
	}
          

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *KYCChainCode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "get" { //read a variable
		return t.get(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func (t *SimpleChaincode) set(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    if len(args) != 2 {
            return "", fmt.Errorf("Incorrect arguments. Expecting arg1 = gci and arg2 = name")
    }

    var currentKYCDataObj KYCData

    currentKYCDataObj.GCI = args[0]
    currentKYCDataObj.NAME = args[1]
    currentKYCDataObj.UPDATED_TIME = time.Now().Format("Mon Jan _2 15:04:05 2006")

    var KYCDataObjBlocks []KYCData

    KYCDataObjBlocksAsBytes,_ := stub.GetState(args[0])

    if KYCDataObjBlocksAsBytes != nil{
        json.Unmarshal(KYCDataObjBlocksAsBytes, &KYCDataObjBlocks)
        KYCDataObjBlocks = append(KYCDataObjBlocks,currentKYCDataObj)
    }else {
        KYCDataObjBlocks = make([]KYCData, 1)
        KYCDataObjBlocks[0] = currentKYCDataObj
    }

    updatedKYCDataObjBlocksAsBytes,_ := json.Marshal(KYCDataObjBlocks)

    err := stub.PutState(args[0], updatedKYCDataObjBlocksAsBytes)
    if err != nil {
            return "", fmt.Errorf("Failed to add KYC details for GCI : %s", args[0])
    }
    return "success", nil
}

// Get returns the value of the specified asset key
func (t *KYCChainCode) get(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    if len(args) != 1 {
            return "", fmt.Errorf("Incorrect arguments. Expecting a key")
    }

    KYCDataObjBlocksAsBytes, err := stub.GetState(args[0])
    if err != nil {
            return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
    }
    if KYCDataObjBlocksAsBytes == nil {
            return "", fmt.Errorf("Asset not found: %s", args[0])
    }

    return KYCDataObjBlocksAsBytes, nil
}