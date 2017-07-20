package main

import (
    "fmt"
    "time"
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

// KYCChainCode implements a simple chaincode to manage an asset
type KYCChainCode struct {
}

type KYCData struct{
    GCI string `json:"gci"`
    NAME string `json:"name"`
    UPDATED_TIME string `json:"updated_time"`
}

// main function starts up the chaincode in the container during instantiate
func main() {
    if err := shim.Start(new(KYCChainCode)); err != nil {
            fmt.Printf("Error starting KYCChainCode chaincode: %s", err)
    }
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *KYCChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
    // Get the args from the transaction proposal
    //args := stub.GetStringArgs()
    // We store the key and the value on the ledger
    value,_ := json.Marshal(true)
    err := stub.PutState("initStatus", value)
    if err != nil {
            return shim.Error(fmt.Sprintf("Failed to Initialize chaincode"))
    }
    return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *KYCChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
    // Extract the function and args from the transaction proposal
    fn, args := stub.GetFunctionAndParameters()

    var result string
    var err error
    if fn == "set" {
            result, err = set(stub, args)
    } else if fn == "get" { 
            result, err = get(stub, args)
    } else {
            result, err = "", fmt.Errorf("Invoke doesn't support function : %s", args[0])
    }
    if err != nil {
            return shim.Error(err.Error())
    }

    // Return the result as success payload
    return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
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
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
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

    return string(KYCDataObjBlocksAsBytes), nil
}

