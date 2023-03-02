package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Contract: captures contract high-level information, such as its identification number, duration, and status.
type Contract struct {
	ID               string `json:"ID"`     // This is the ID that will identify this contract in the ledger.
	Status           string `json:"Status"` // Can only be Pending, Active, and Terminated. In next version we will add Rejected, and Completed.
	Notes            string `json:"Notes"`
	StartDate        string `json:"Start date"`
	EndDate          string `json:"End date"`
	ExtensionDetails string `json:"Extension details"`
	Employer         Employer
	Employee         Employee
	Job              Job
	Benefits         Benefits
	Disputes         []Dispute
}

// Employer: provides data about the employer, such as name and address details
type Employer struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	EmployerAC string `json:"Employer address and contact details"`
	Country    string `json:"Country"`
}

// Employee: provides data about the employee, such as name and contact details
type Employee struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	EmployeeAC string `json:"Employee address and contact details"`
	Country    string `json:"Country"`
}

// Job: describes the job details, such as position and task description
type Job struct {
	Position    string `json:"Position"`
	Level       string `json:"Level"`
	Description string `json:"Description"`
}

// Benefits: states the job benefits, such salary, allowances, and annual increase
type Benefits struct {
	Currency       string `json:"Currency"`
	Salary         int    `json:"Salary"`
	AnnualIncrease string `json:"Annual increase"`
	AnnualLeave    string `json:"Annual leave"`
	Housing        int    `json:"Housing"`
	Allowances     int    `json:"Allowances"`
	OtherBenefits  string `json:"Other benefits"`
}

// Disputes: lists the disputes, if any, that is raised by the employee with their content and the last update dates.
type Dispute struct {
	ID              string `json:"ID"`
	Status          string `json:"Status"` // Can only be Active or Closed.
	LastUpdatedDate string `json:"Last updated date"`
	Content         string `json:"Content"`
	Responses       []Response
}

// Responses: list the employer responses, if any, for disputes raised by the employee.
type Response struct {
	ID              string `json:"ID"`
	LastUpdatedDate string `json:"Last updated date"`
	Content         string `json:"Content"`
}

// ReadContract returns the contract stored in the world state with given id. will return nil if nothing is found.
func (s *SmartContract) ReadContract(ctx contractapi.TransactionContextInterface, ID string) (*Contract, error) {

	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if contractJSON == nil {
		return nil, fmt.Errorf("the contract %s does not exist", ID)
	}
	var contract Contract
	err = json.Unmarshal(contractJSON, &contract)
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

/*
* If the jsonString is a valid contract, with matching id in the blockchain this method will return true.
* This method will not change old disputes, responses, and status.
* @param jsonString represents The new contract.
 */
func (s *SmartContract) UpdateContract(ctx contractapi.TransactionContextInterface, jsonString string) (bool, error) {

	//Parsing jsonString
	jsonString = strings.ReplaceAll(jsonString, "'", "\"")
	var contract Contract
	err := json.Unmarshal([]byte(jsonString), &contract)
	if err != nil {
		return false, fmt.Errorf("Error Unmarshaling JSON: %s, \n %s", err, jsonString)
	}

	// If we don't find the contract in the blockchain we stop.
	exists, err := s.ContractExist(ctx, contract.ID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("the Contract %s doesn't exists", contract.ID)
	}

	// We need the old contract because it the has  old disputes
	contractJSON, err := ctx.GetStub().GetState(contract.ID)
	if err != nil {
		return false, err
	}
	var oldContract Contract
	json.Unmarshal(contractJSON, &oldContract)
	if oldContract.Status == "Terminated" {
		return false, fmt.Errorf("You can not update a terminated contract.")
	}

	// Each create method will check for the necessary information and conditions.
	employer, flag1 := createEmployer(ctx, contract.Employer.ID, contract.Employer.Name, contract.Employer.EmployerAC, contract.Employer.Country)
	employee, flag2 := createEmployee(ctx, contract.Employee.ID, contract.Employee.Name, contract.Employee.EmployeeAC, contract.Employee.Country)
	job, flag3 := createJob(ctx, contract.Job.Position, contract.Job.Level, contract.Job.Description)
	benefits, flag4 := createBenefits(ctx, contract.Benefits.Currency, contract.Benefits.Salary, contract.Benefits.AnnualIncrease, contract.Benefits.AnnualLeave,
		contract.Benefits.Housing, contract.Benefits.Allowances, contract.Benefits.OtherBenefits)

	if !(flag1 && flag2 && flag3 && flag4) {
		return false, fmt.Errorf("Failed to initialize one of the structs.")
	}

	newContract := Contract{
		ID:               contract.ID,
		Status:           oldContract.Status, // Status only changes using TerminateContract and ApproveContract.
		Notes:            contract.Notes,
		StartDate:        contract.StartDate,
		EndDate:          contract.EndDate,
		ExtensionDetails: contract.ExtensionDetails,
		Employer:         employer,
		Employee:         employee,
		Job:              job,
		Benefits:         benefits,
		Disputes:         oldContract.Disputes, // Updating Disputes & Responses is outside the scope of this method.
	}

	contractJson, err := json.Marshal(newContract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(newContract.ID, contractJson)
	if err != nil {
		return false, err
	}

	return true, nil

}

/*
* Given an existing contract this method will append a new dispute into []Disputes, and return true.
* @Param Content must not be empty.
 */
func (s *SmartContract) IssueDispute(ctx contractapi.TransactionContextInterface, ID string, Content string) (bool, error) {
	curDate := time.Now()
	dispute := Dispute{
		ID:              "1", // should be modified later
		Status:          "Active",
		LastUpdatedDate: curDate.Format("01/02/2006"),
		Content:         Content,
		Responses:       []Response{},
	}

	// If the given dispute is faulty return false.
	if !(checkDispute(ctx, dispute)) {
		return false, fmt.Errorf("The given dispute doesn't meet all proper conditions. ")
	}

	// if the contract doesn't exist return false.
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil || contractJSON == nil {
		return false, fmt.Errorf("The given ID doesn't match any contract in the blockchain. ")
	}

	// If any old dispute shares the ID of the new dispute return false.
	var oldContract Contract
	json.Unmarshal(contractJSON, &oldContract)

	// Now the smart contract will give a dispute id by itself without user input.
	dispute.ID = strconv.Itoa(len(oldContract.Disputes))
	oldContract.Disputes = append(oldContract.Disputes, dispute)

	contractJson, err := json.Marshal(oldContract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(ID, contractJson)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
* Given an existing contract and a valid unique dispute this method will update the dispute, and return true.
* Only dispute.ID must be unique. @Param Content must not be empty.
 */
func (s *SmartContract) UpdateDispute(ctx contractapi.TransactionContextInterface, ID string, DID string, Content string) (bool, error) {

	curDate := time.Now()
	dispute := Dispute{
		ID:              DID,
		Status:          "Active",
		LastUpdatedDate: curDate.Format("01/02/2006"),
		Content:         Content,
		Responses:       []Response{},
	}

	// If the given dispute is faulty return false.
	if !(checkDispute(ctx, dispute)) {
		return false, fmt.Errorf("the given dispute doesn't meet all proper conditions")
	}

	// if the contract doesn't exist return false.
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil || contractJSON == nil {
		return false, fmt.Errorf("The given ID doesn't match any contract in the blockchain. ")
	}

	// We will verify if there is a matching dispute and if found we will modify it.
	flag := false
	var oldContract Contract
	json.Unmarshal(contractJSON, &oldContract)
	for i := 0; i < len(oldContract.Disputes); i++ {
		if oldContract.Disputes[i].ID == dispute.ID {
			if oldContract.Disputes[i].Status == "Closed" {
				return false, fmt.Errorf("You can't update a closed dispute.")
			}
			flag = true
			oldResponse := oldContract.Disputes[i].Responses // because ldContract.Disputes[i] = dispute will override responses.
			oldContract.Disputes[i] = dispute
			oldContract.Disputes[i].Responses = oldResponse
			break
		}
	}

	// If no matching dispute is found we leave.
	if !(flag) {
		return false, fmt.Errorf("The given DisputeID doesn't match any dispute on this contract.")
	}

	contractJson, err := json.Marshal(oldContract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(ID, contractJson)
	if err != nil {
		return false, err
	}
	return true, nil

}

// This method will return true if the ID and dispute.ID match an existing contract. Will return false if Dispute.status == "Closed"
func (s *SmartContract) CloseDispute(ctx contractapi.TransactionContextInterface, ID string, DisputeID string) (bool, error) {

	// if the contract doesn't exist return false.
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil || contractJSON == nil {
		return false, fmt.Errorf("The given ID doesn't match any contract in the blockchain. ")
	}

	var oldContract Contract
	json.Unmarshal(contractJSON, &oldContract)
	// Will search the Contract.Disputes for a matching dispute ID. If there is a match it must be active otherwise will return false.
	flag := false
	for i := 0; i < len(oldContract.Disputes); i++ {
		if oldContract.Disputes[i].ID == DisputeID {
			if oldContract.Disputes[i].Status == "Active" {
				curDate := time.Now()
				oldContract.Disputes[i].Status = "Closed"
				oldContract.Disputes[i].LastUpdatedDate = curDate.Format("01/02/2006")
				flag = true
				break
			} else {
				return false, fmt.Errorf("This contract is already closed.")
			}
		}
	}

	// Flag will be set to true if we find an active dispute with a matching id.
	if !flag {
		return false, fmt.Errorf("There is no matching disputes with the given ID: %s ", DisputeID)
	}

	contractJson, err := json.Marshal(oldContract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(ID, contractJson)
	if err != nil {
		return false, err
	}
	return true, nil
}

// This method responds to an open dispute.
// Given an existing ID and DisputeID, and a unique RespondID this method will return true. @Param Content must not be empty.
func (s *SmartContract) RespondToDispute(ctx contractapi.TransactionContextInterface, ID string, DisputeID string, Content string) (bool, error) {
	curDate := time.Now()
	response := Response{
		ID:              "Empty",
		LastUpdatedDate: curDate.Format("01/02/2006"),
		Content:         Content,
	}

	var flag bool
	response, flag = createResponse(ctx, response) // If the response structure is valid flag will be set to true
	if !flag {
		return false, fmt.Errorf("The response structure is invalid.")
	}
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil || contractJSON == nil {
		return false, fmt.Errorf("The given ID doesn't match any contract in the blockchain. ")
	}

	var contract Contract
	json.Unmarshal(contractJSON, &contract)

	// Will search the Contract.Disputes for a matching dispute ID. If there is a match it must be active otherwise will return false.
	flag = false
	var disputePos int
	for i := 0; i < len(contract.Disputes); i++ {
		if contract.Disputes[i].ID == DisputeID {
			if contract.Disputes[i].Status == "Active" {
				flag = true
				disputePos = i
				break
			} else {
				return false, fmt.Errorf("This contract is already closed.")
			}
		}
	}
	// Flag will be set to true if we find an active dispute with a matching id.
	if !flag {
		return false, fmt.Errorf("There is no matching disputes with the given ID: %s ", DisputeID)
	}

	// Because we don't want the user to enter response ID.
	response.ID = strconv.Itoa(len(contract.Disputes[disputePos].Responses))
	contract.Disputes[disputePos].Responses = append(contract.Disputes[disputePos].Responses, response)
	contractJson, err := json.Marshal(contract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(ID, contractJson)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SmartContract) HandleAddContract(ctx contractapi.TransactionContextInterface, jsonString string) (bool, error) {

	//Parsing jsonString
	jsonString = strings.ReplaceAll(jsonString, "'", "\"")
	var contract Contract
	err := json.Unmarshal([]byte(jsonString), &contract)
	if err != nil {
		return false, fmt.Errorf("Error Unmarshaling JSON: %s, \n %s", err, jsonString)
	}

	// Check if there is an existing contract with the same ID.
	exists, err := s.ContractExist(ctx, contract.ID)
	if err != nil {
		return false, err
	}
	if exists {
		return false, fmt.Errorf("the Contract %s already exists", contract.ID)
	}

	// To simplify the process we separate every struct.
	tempEmployer := contract.Employer
	tempEmployee := contract.Employee
	tempJob := contract.Job
	tempBenefits := contract.Benefits

	employer, flag1 := createEmployer(ctx, tempEmployer.ID, tempEmployer.Name, tempEmployer.EmployerAC, tempEmployer.Country)
	employee, flag2 := createEmployee(ctx, tempEmployee.ID, tempEmployee.Name, tempEmployee.EmployeeAC, tempEmployee.Country)
	job, flag3 := createJob(ctx, tempJob.Position, tempJob.Level, tempJob.Description)
	benefits, flag4 := createBenefits(ctx, tempBenefits.Currency, tempBenefits.Salary, tempBenefits.AnnualIncrease, tempBenefits.AnnualLeave,
		tempBenefits.Housing, tempBenefits.Allowances, tempBenefits.OtherBenefits)

	// in case one of the structs failed to initialize.
	if !(flag1 && flag2 && flag3 && flag4) {
		return false, err
	}

	flag5, err := checkDate(contract.StartDate, contract.EndDate)
	if !(flag5) {
		return false, fmt.Errorf("Time has failed Error: %s", err)
	}

	var disputes = []Dispute{} // A new contract should have no disputes.

	contract.Status = "Pending" // Every new contract should start with status as pending. Will ignore jsonString input.
	// In case one of the create method changes anything in the struct
	contract.Employer = employer
	contract.Employee = employee
	contract.Job = job
	contract.Benefits = benefits
	contract.Disputes = disputes

	contractJson, err := json.Marshal(contract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(contract.ID, contractJson)
	if err != nil {
		return false, err
	}

	return true, nil

}

/*
* This method will change the contract status from Active to Terminated.
* Will return true only if ID points to existing contract and said contract status is not Terminated.
 */
func (s *SmartContract) TerminateContract(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {
	// Retrieving the contract from world state and unmarshaling into contract.
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if contractJSON == nil {
		return false, fmt.Errorf("the contract %s does not exist", ID)
	}
	var contract Contract
	err = json.Unmarshal(contractJSON, &contract)
	if err != nil {
		return false, err
	}

	if contract.Status == "Terminated" {
		return false, fmt.Errorf("The contract is already %s.", contract.Status)
	}

	contract.Status = "Terminated"

	contractJson, err := json.Marshal(contract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(contract.ID, contractJson)
	if err != nil {
		return false, err
	}

	return true, nil

}

/*
* ID represents the contract ID.
* This method will change the contract status from Pending to Active.
* Will return true only if ID points to existing contract and said contract status is not Terminated/Active.
 */
func (s *SmartContract) ApproveContract(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {

	// Retrieving the contract from world state and unmarshaling into contract.
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if contractJSON == nil {
		return false, fmt.Errorf("the contract %s does not exist", ID)
	}
	var contract Contract
	err = json.Unmarshal(contractJSON, &contract)
	if err != nil {
		return false, err
	}

	if contract.Status == "Active" || contract.Status == "Terminated" {
		return false, fmt.Errorf("The contract is %s", contract.Status)
	}

	contract.Status = "Active"

	contractJson, err := json.Marshal(contract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(contract.ID, contractJson)
	if err != nil {
		return false, err
	}

	return true, nil
}

/*
* This method will update the contract duration.
* @Param ToDate represents the new date at which the contract will end.
* To return true ToDate must be further than the EndDate and currentDate should be no less than three months of currentDate.
 */
func (s *SmartContract) ExtendContract(ctx contractapi.TransactionContextInterface, ID string, ToDate string) (bool, error) {

	// Retrieving the contract from world state and unmarshaling into contract.
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if contractJSON == nil {
		return false, fmt.Errorf("the contract %s does not exist", ID)
	}
	var contract Contract
	err = json.Unmarshal(contractJSON, &contract)
	if err != nil {
		return false, err
	}

	if !(contract.Status == "Active") {
		return false, fmt.Errorf("The contract is not Active")
	}
	isCorrectDate, err := extensionCheck(ToDate, contract.EndDate)

	if !isCorrectDate {
		return false, err
	}

	contract.EndDate = ToDate
	contractJson, err := json.Marshal(contract)
	if err != nil {
		return false, err
	}
	err = ctx.GetStub().PutState(contract.ID, contractJson)
	if err != nil {
		return false, err
	}

	return true, nil

}

// Return true if inputs are not empty.
func createEmployer(ctx contractapi.TransactionContextInterface, ID string, Name string, EmployerAC string, Country string) (Employer, bool) {
	if ID == "" || Name == "" || EmployerAC == "" || Country == "" {
		return Employer{}, false
	}

	employer := Employer{
		ID:         ID,
		Name:       Name,
		EmployerAC: EmployerAC,
		Country:    Country,
	}
	return employer, true
}

// Return true if inputs are not empty.
func createEmployee(ctx contractapi.TransactionContextInterface, ID string, Name string, EmployeeAC string, Country string) (Employee, bool) {
	if ID == "" || Name == "" || EmployeeAC == "" || Country == "" {
		return Employee{}, false
	}

	employee := Employee{
		ID:         ID,
		Name:       Name,
		EmployeeAC: EmployeeAC,
		Country:    Country,
	}
	return employee, true
}

// Return true if inputs are not empty.
func createJob(ctx contractapi.TransactionContextInterface, Position string, Level string, Description string) (Job, bool) {
	if Position == "" || Level == "" || Description == "" {
		return Job{}, false
	}

	job := Job{
		Position:    Position,
		Level:       Level,
		Description: Description,
	}
	return job, true
}

// Return true if inputs are not empty, and the salary is not 0.
func createBenefits(ctx contractapi.TransactionContextInterface, Currency string, Salary int, AnnualIncrease string, AnnualLeave string, Housing int, Allowances int, OtherBenefits string) (Benefits, bool) {
	// Housing and allowances are missing because they are not as critical as the rest.
	if Currency == "" || AnnualIncrease == "" || AnnualLeave == "" || OtherBenefits == "" || Salary == 0 {
		return Benefits{}, false
	}

	benefits := Benefits{
		Currency:       Currency,
		Salary:         Salary,
		AnnualIncrease: AnnualIncrease,
		AnnualLeave:    AnnualLeave,
		OtherBenefits:  OtherBenefits,
	}
	return benefits, true
}

// Will return false if the given ID doesn't match any record in the blockchain.
func (s *SmartContract) ContractExist(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {
	contractJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	// It is false if we don't find anything
	if contractJSON == nil {
		return false, nil
	}

	return true, nil
}

// Will check if the given Dispute is up to standers. Returns true if the dispute is made correctly
func checkDispute(ctx contractapi.TransactionContextInterface, dispute Dispute) bool {
	if dispute.ID == "" || dispute.LastUpdatedDate == "" || dispute.Content == "" {
		return false
	}
	return true
}

// Will check if the given Dispute is up to standers. Returns true and Dispute if the dispute is made correctly
func createDispute(ctx contractapi.TransactionContextInterface, dispute Dispute) (Dispute, bool) {
	if dispute.ID == "" || dispute.LastUpdatedDate == "" || dispute.Content == "" {
		return Dispute{}, false
	}
	return dispute, true
}

// Will return true if the response structure is valid
func createResponse(ctx contractapi.TransactionContextInterface, response Response) (Response, bool) {
	if response.ID == "" || response.LastUpdatedDate == "" || response.Content == "" {
		return Response{}, false
	}
	return response, true
}

// Will check the contract start date and end date. If they are valid will return true otherwise will return false.
func checkDate(startDate string, endDate string) (bool, error) {
	currentDate := time.Now()
	_, err := time.Parse("01/02/2006", startDate) // We only need to know if the startDate is in correct format or not.
	EndDate, err1 := time.Parse("01/02/2006", endDate)

	if err != nil || err1 != nil {
		return false, err
	}

	if EndDate.Before(currentDate) {
		return false, fmt.Errorf("You can't create a new contract in the past. Please check End date")
	}

	return true, nil

}

// This method will verify if the conditions of the contract allow an extension. Also will check if the new dates are valid.
// ToDate represents the new end date we want to extend to. endDate represents the contract old date.
func extensionCheck(toDate string, endDate string) (bool, error) {
	currentDate := time.Now()
	ExtendedDate, err := time.Parse("01/02/2006", toDate) // The new end date.
	if err != nil {
		return false, err
	}
	EndDate, err := time.Parse("01/02/2006", endDate)
	if err != nil {
		return false, err
	}
	currentDatePlus := currentDate.AddDate(0, 3, 0) // To check if there is 3 months or less left on the contract

	if !currentDate.Before(ExtendedDate) {
		return false, fmt.Errorf("You can't extend the contract to a date in the past, CurrentDate: %s, You want to extend it to: %s", currentDate, ExtendedDate)
	}

	if ExtendedDate.Before(EndDate) {
		return false, fmt.Errorf("You can't shorten the length of the contract using this method.")
	}

	if !currentDatePlus.After(EndDate) {
		return false, fmt.Errorf("You can only extend the contract in the last three months.")
	}

	return true, nil
}

// Returns all assets found in the world state
func (s *SmartContract) GetAllContracts(ctx contractapi.TransactionContextInterface) ([]*Contract, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var contracts []*Contract
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var contract Contract
		err = json.Unmarshal(queryResponse.Value, &contract)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, &contract)
	}

	return contracts, nil
}

// getEmployeeHistory will sort through every contract in the blockchain. Will return an array of contracts containing all contracts of an employee.
func (s *SmartContract) getEmployeeHistory(ctx contractapi.TransactionContextInterface, EmployeeID string) ([]Contract, error) {

	Contracts, _ := s.GetAllContracts(ctx)

	var EmployeeContracts []Contract

	for i := 0; i < len(Contracts); i++ {
		if Contracts[i].Employee.ID == EmployeeID {
			EmployeeContracts = append(EmployeeContracts, *Contracts[i])

		}
	}

	return EmployeeContracts, nil
}

func (s *SmartContract) getEmployerHistory(ctx contractapi.TransactionContextInterface, EmployerID string) ([]Contract, error) {

	Contracts, _ := s.GetAllContracts(ctx)

	var EmployeeContracts []Contract

	for i := 0; i < len(Contracts); i++ {
		if Contracts[i].Employer.ID == EmployerID {
			EmployeeContracts = append(EmployeeContracts, *Contracts[i])

		}
	}

	return EmployeeContracts, nil
}

// Employee data struct let us reorganize data in a way that is easy to understand.
type EmployeeData struct { // We will reuse this struct for EmployerData.
	Contracts           string
	TerminatedContracts string
	ActiveContracts     string
	PendingContracts    string
	Disputes            string
	OpenDisputes        string
	ClosedDisputes      string
}

// This method will reorganize employee history.
func (s *SmartContract) ViewEmployeeHistory(ctx contractapi.TransactionContextInterface, EmployeeID string) EmployeeData {
	EmployeeContracts, _ := s.getEmployeeHistory(ctx, EmployeeID)

	var activeContracts, terminatedContracts, pendingContracts = 0, 0, 0
	var totalDisputes, openDisputes, closedDisputes = 0, 0, 0
	for i := 0; i < len(EmployeeContracts); i++ {
		totalDisputes += len(EmployeeContracts[i].Disputes)
		if EmployeeContracts[i].Status == "Active" {
			activeContracts += 1
		} else if EmployeeContracts[i].Status == "Pending" {
			pendingContracts += 1
		} else {
			terminatedContracts += 1
		}
		for j := 0; j < len(EmployeeContracts[i].Disputes); j++ {
			if EmployeeContracts[i].Disputes[j].Status == "Active" {
				openDisputes += 1
			} else {
				closedDisputes += 1
			}

		}

	}

	// We need to convert them into strings.
	EmployeeData := EmployeeData{
		Contracts:           strconv.Itoa(len(EmployeeContracts)),
		TerminatedContracts: strconv.Itoa(terminatedContracts),
		ActiveContracts:     strconv.Itoa(activeContracts),
		PendingContracts:    strconv.Itoa(pendingContracts),
		Disputes:            strconv.Itoa(totalDisputes),
		OpenDisputes:        strconv.Itoa(openDisputes),
		ClosedDisputes:      strconv.Itoa(closedDisputes),
	}

	return EmployeeData
}

// This method will reorganize employer history.
func (s *SmartContract) ViewEmployerHistory(ctx contractapi.TransactionContextInterface, EmployeeID string) EmployeeData {
	EmployeeContracts, _ := s.getEmployerHistory(ctx, EmployeeID)

	var activeContracts, terminatedContracts, pendingContracts = 0, 0, 0
	var totalDisputes, openDisputes, closedDisputes = 0, 0, 0
	for i := 0; i < len(EmployeeContracts); i++ {
		totalDisputes += len(EmployeeContracts[i].Disputes)
		if EmployeeContracts[i].Status == "Active" {
			activeContracts += 1
		} else if EmployeeContracts[i].Status == "Pending" {
			pendingContracts += 1
		} else {
			terminatedContracts += 1
		}
		for j := 0; j < len(EmployeeContracts[i].Disputes); j++ {
			if EmployeeContracts[i].Disputes[j].Status == "Active" {
				openDisputes += 1
			} else {
				closedDisputes += 1
			}

		}

	}

	// We need to convert them into strings.
	EmployeeData := EmployeeData{
		Contracts:           strconv.Itoa(len(EmployeeContracts)),
		TerminatedContracts: strconv.Itoa(terminatedContracts),
		ActiveContracts:     strconv.Itoa(activeContracts),
		PendingContracts:    strconv.Itoa(pendingContracts),
		Disputes:            strconv.Itoa(totalDisputes),
		OpenDisputes:        strconv.Itoa(openDisputes),
		ClosedDisputes:      strconv.Itoa(closedDisputes),
	}

	return EmployeeData
}
