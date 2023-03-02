package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/olekukonko/tablewriter"
)

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

func printScreen() {
	options := []string{
		"1. Add Contract",
		"2. Approve Contract",
		"3. Update Contract",
		"4. Extend Contract",
		"5. Terminate Contract",
		"6. Issue Dispute",
		"7. Update Dispute",
		"8. Close Dispute",
		"9. Respond to Dispute",
		"10. Read Contract",
		"11. View Employee History",
		"12. View Employer History",
		"13. Get All Contracts",
	}

	numCols := 3
	numRows := 5

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetColumnSeparator("  ")
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderLine(false)
	table.SetRowLine(false)
	table.SetTablePadding("\t")
	table.SetHeader([]string{"Options", "", "Options", "", "Options"})
	table.SetRowSeparator("-")
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")

	for i := 0; i < numRows; i++ {
		row := []string{}
		for j := 0; j < numCols; j++ {
			idx := i + j*numRows
			if idx >= len(options) {
				row = append(row, "", "")
			} else {
				row = append(row, options[idx], "")
			}
		}
		table.Append(row)
	}

	table.Render()
	fmt.Println()
	fmt.Println("Please select a transaction to execute")
}

func choice() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your choice: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return 0
	}

	choice = strings.TrimSpace(choice)

	methodNumber, err := strconv.Atoi(choice)
	if err != nil {
		fmt.Printf("Invalid input. Please enter a number between 1 and 15. \n")
		return 0
	}

	if methodNumber < 1 || methodNumber > 15 {
		fmt.Printf("Invalid input. Please enter a number between 1 and 15. \n")
		return 0
	}
	return methodNumber

}

func main() {
	for j := 0; j < 20; j++ { // the program will loop for 20 times only
		printScreen()
		var input int = 0
		for i := 0; i < 100; i++ {
			input = choice()
			if input != 0 {
				break
			}
			if i == 99 {
				println("You have exceeded the limit.")
				return
			}
		}

		switch input {
		case 1:

			fmt.Println("You selected to execute create contract transaction ")
			createContract()
		case 2:
			fmt.Println("You selected to execute approve contract transaction ")
			approveContract()
		case 3:
			fmt.Println("You selected to execute update contract transaction ")
			updateContract()
		case 4:
			fmt.Println("You selected to execute extend contract transaction ")
			extendContract()
		case 5:
			fmt.Println("You selected to execute terminate contract transaction ")
			terminateContract()
		case 6:
			fmt.Println("You selected to execute issue dispute transaction ")
			issueDispute()
		case 7:
			fmt.Println("You selected to execute update dispute transaction ")
			updateDispute()
		case 8:
			fmt.Println("You selected to execute close dispute transaction ")
			closeDispute()
		case 9:
			fmt.Println("You selected to execute respond to dispute transaction ")
			respondToDispute()
		case 10:
			fmt.Println("You selected to execute read contract transaction ")
			prettifyContract(readContract())
		case 11:
			fmt.Println("You selected to execute view employee history transaction ")
			viewEmployeeHistory()
		case 12:
			fmt.Println("You selected to execute view employer history transaction ")
			viewEmployerHistory()
		case 13:
			fmt.Println("You selected to execute view all contracts transaction ")
			prettifyAllContracts(getAllContracts())
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Println()
		fmt.Print("Press Enter to continue.")
		reader.ReadString('\n')
	}

}

// will post your request to Fablo rest api and return the response. Without changing anything.
func postRequest(input string, methodName string) string {
	client := &http.Client{}
	var data = strings.NewReader(`{"method": "` + methodName + `",
"args": [` + input + `]}`)
	req, err := http.NewRequest("POST", "http://localhost:8801/invoke/my-channel/chaincode1", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	token := getToken()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(bodyText)
}

// This method will determine if the response from the smart contract is an error or not.
// if @Param err represents an error this method will return true.
func isError(err string) bool {

	flag1 := strings.Contains(err, "\"message\":")
	flag2 := strings.Contains(err, "message=")

	return (flag1 && flag2) // Both flags must be true to be sure that the returned values is an error.
}

// Will format and print the error.
func printError(err string) {
	_, err, _ = strings.Cut(err, "message=")
	_, err, _ = strings.Cut(err, "message=") // Because we got two peers.

	err = strings.ReplaceAll(err, "\"}", "")
	err = strings.ReplaceAll(err, "\"", "")
	fmt.Printf("Error: %s", err)
}

func updateDispute() {
	// Taking all the need inputs from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	fmt.Print("Enter Dispute ID: ")
	DisputeID, err2 := reader.ReadString('\n')

	fmt.Print("Enter the new content of your dispute: ")
	Content, err3 := reader.ReadString('\n')

	if err != nil || err2 != nil || err3 != nil {
		fmt.Printf("Could not read string \n")
	}

	combinedInputs := combineStrings(ID, DisputeID, Content)
	combinedInputs = strings.ReplaceAll(combinedInputs, "\n", "")

	// Will send and get a response from the blockchain
	bodyText := postRequest(combinedInputs, "UpdateDispute")
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}

	println(jsonString)
}

func approveContract() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string \n")
	}

	ID = combineStrings(ID)
	ID = strings.ReplaceAll(ID, "\n", "")
	bodyText := postRequest(ID, "ApproveContract")

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}

	println("Contract has been approved.")
	prettifyTopContract(choseContract(ID))

}

func extendContract() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string \n")
	}

	fmt.Print("Enter the extension date in the following format: 01/01/2023. You must add the 0 in 01:  ")
	ToDate, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string2 \n")
	}

	ToDate = strings.ReplaceAll(ToDate, "\n", "")
	ID = strings.ReplaceAll(ID, "\n", "")

	input := combineStrings(ID, ToDate)
	input = strings.ReplaceAll(input, "\n", "")

	bodyText := postRequest(input, "ExtendContract")

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}

	ID = combineStrings(ID) // To properly format ID
	ID = strings.ReplaceAll(ID, "\n", "")

	println("The contract has been extended successfully")
	prettifyTopContract(choseContract(ID))

}

func terminateContract() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string \n")
	}

	ID = combineStrings(ID)
	ID = strings.ReplaceAll(ID, "\n", "")
	bodyText := postRequest(ID, "TerminateContract")

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}

	prettifyTopContract(choseContract(ID))

}

func respondToDispute() {
	// Taking all the need inputs from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	fmt.Print("Enter Dispute ID: ")
	DisputeID, err2 := reader.ReadString('\n')

	// fmt.Print("Enter Response ID: ")
	// RID, err4 := reader.ReadString('\n')

	fmt.Print("Enter the content of your response: ")
	Content, err3 := reader.ReadString('\n')

	if err != nil || err2 != nil || err3 != nil {
		fmt.Printf("Could not read string \n")
	}

	combinedInputs := combineStrings(ID, DisputeID, Content)
	combinedInputs = strings.ReplaceAll(combinedInputs, "\n", "")

	// Will send and get a response from the blockchain
	bodyText := postRequest(combinedInputs, "RespondToDispute")
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}

	ID = combineStrings(ID)
	ID = strings.ReplaceAll(ID, "\n", "")
	prettifyDispute(choseContract(ID))
}

func closeDispute() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	fmt.Print("Enter Dispute ID: ")
	DisputeID, err2 := reader.ReadString('\n')

	if err != nil || err2 != nil {
		fmt.Printf("Could not read string \n")
	}

	combinedInputs := combineStrings(ID, DisputeID)
	combinedInputs = strings.ReplaceAll(combinedInputs, "\n", "")

	// Will send and get a response from the blockchain
	bodyText := postRequest(combinedInputs, "CloseDispute")
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}

	ID = combineStrings(ID)
	ID = strings.ReplaceAll(ID, "\n", "")
	prettifyDispute(choseContract(ID))
}

func issueDispute() {
	// Taking all the need inputs from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	// fmt.Print("Enter Dispute ID: ")
	// DisputeID, err2 := reader.ReadString('\n')

	fmt.Print("Enter the content of your dispute: ")
	Content1, err3 := reader.ReadString('\n')

	if err != nil || err3 != nil {
		fmt.Printf("Could not read string \n")
	}

	combinedInputs := combineStrings(ID, Content1)
	combinedInputs = strings.ReplaceAll(combinedInputs, "\n", "")

	// Will send and get a response from the blockchain
	bodyText := postRequest(combinedInputs, "IssueDispute")
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	if isError(jsonString) {
		printError(jsonString)
		return
	}
	ID = strings.ReplaceAll(ID, "\n", "")
	ID = combineStrings(ID)
	prettifyDispute(choseContract(ID))

}

// Get will return a token without any spaces.
func getToken() string {
	client := &http.Client{}
	var data = strings.NewReader(`{"id": "admin", "secret": "adminpw"}`)
	req, err := http.NewRequest("POST", "http://localhost:8801/user/enroll", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer ")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	str := string(bodyText)

	str = strings.TrimPrefix(str, "{\"token\":\"")
	str = strings.TrimSuffix(str, "\"}")

	return str
}

type EmployeeData struct {
	Contracts           string
	TerminatedContracts string
	ActiveContracts     string
	PendingContracts    string
	Disputes            string
	OpenDisputes        string
	ClosedDisputes      string
}

func viewEmployeeHistory() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter EmployeeID: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return
	}
	choice = strings.ReplaceAll(choice, "\n", "")
	choice = combineStrings(choice)

	// Will send and get a response from the blockchain
	bodyText := postRequest(choice, "ViewEmployeeHistory")
	if bodyText == "" {
		return
	}

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	EmpData := EmployeeData{}
	json.Unmarshal([]byte(jsonString), &EmpData)

	if EmpData.Contracts == "0" {
		fmt.Println("Invalid EmployeeID. Please try again.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Contracts", "Terminated Contracts", "Active Contracts", "Pending Contracts", "Disputes", "Open Disputes", "Closed Disputes"})
	table.SetRowLine(true)

	// Set the table style
	table.SetBorder(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Append employee data.
	row := []string{EmpData.Contracts, EmpData.TerminatedContracts, EmpData.ActiveContracts, EmpData.PendingContracts, EmpData.Disputes, EmpData.OpenDisputes, EmpData.ClosedDisputes}
	table.Append(row)

	table.Render()

}

func viewEmployerHistory() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter EmployerID: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return
	}
	choice = strings.ReplaceAll(choice, "\n", "")
	choice = combineStrings(choice)

	// Will send and get a response from the blockchain
	bodyText := postRequest(choice, "ViewEmployerHistory")
	if bodyText == "" {
		return
	}
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	EmpData := EmployeeData{}
	json.Unmarshal([]byte(jsonString), &EmpData)

	if EmpData.Contracts == "0" {
		fmt.Println("Invalid EmployerID. Please try again.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Contracts", "Terminated Contracts", "Active Contracts", "Pending Contracts", "Disputes", "Open Disputes", "Closed Disputes"})
	table.SetRowLine(true)
	// Set the table style
	table.SetBorder(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Append employee data.
	row := []string{EmpData.Contracts, EmpData.TerminatedContracts, EmpData.ActiveContracts, EmpData.PendingContracts, EmpData.Disputes, EmpData.OpenDisputes, EmpData.ClosedDisputes}
	table.Append(row)

	table.Render()

}

func updateContract() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the contract file name: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return
	}

	choice = strings.ReplaceAll(choice, "\n", "")
	// Read the contents of the file
	input, err := ioutil.ReadFile(choice)
	if err != nil {
		fmt.Println("Failed to located the file. Please don't forget to add .json at the end.")
		return

	}

	// Convert the byte slice to a string and remove newline characters
	jsonData := string(input)
	jsonData = strings.ReplaceAll(jsonData, "\n", "")
	jsonData = strings.ReplaceAll(jsonData, "\"", "'")
	jsonData = combineStrings(jsonData)

	bodyText := postRequest(jsonData, "UpdateContract")

	if isError(bodyText) {
		printError(bodyText)
		return
	}
	fmt.Println("The contract has been updated.")

	// We need to retrieve the id from the json string
	_, jsonData, _ = strings.Cut(jsonData, ": ")
	jsonData, _, _ = strings.Cut(jsonData, ",")
	jsonData = strings.ReplaceAll(jsonData, "'", "")
	jsonData = combineStrings(jsonData)
	prettifyTopContract(choseContract(jsonData))
}

func createContract() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the contract file name: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return
	}

	choice = strings.ReplaceAll(choice, "\n", "")
	// Read the contents of the file
	input, err := ioutil.ReadFile(choice)
	if err != nil {
		fmt.Println("Failed to located the file. Please don't forget to add .json at the end.")
		return

	}

	// Convert the byte slice to a string and remove newline characters
	jsonData := string(input)
	jsonData = strings.ReplaceAll(jsonData, "\n", "")
	jsonData = strings.ReplaceAll(jsonData, "\"", "'")
	jsonData = combineStrings(jsonData)

	bodyText := postRequest(jsonData, "HandleAddContract")

	if isError(bodyText) {
		printError(bodyText)
		return
	}
	fmt.Println("A contract with the status Pending has been created.")

	// We need to retrieve the id from the json string
	_, jsonData, _ = strings.Cut(jsonData, ": ")
	jsonData, _, _ = strings.Cut(jsonData, ",")
	jsonData = strings.ReplaceAll(jsonData, "'", "")
	jsonData = combineStrings(jsonData)
	prettifyTopContract(choseContract(jsonData))
}

// Will only show the contract top level info
func showContractInfo(jsonString string) {

}

func choseContract(ID string) Contract {

	bodyText := postRequest(ID, "ReadContract")
	contract := Contract{}

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")
	err := json.Unmarshal([]byte(jsonString), &contract)
	if err != nil {
		fmt.Errorf("%s", err)
		return contract
	}
	return contract
}

func readContract() Contract {
	contract := Contract{}

	// Taking input from user.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return contract
	}

	input = strings.ReplaceAll(input, "\n", "")
	input = combineStrings(input)

	// Will send and get a response from the blockchain
	bodyText := postRequest(input, "ReadContract")

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")
	err = json.Unmarshal([]byte(jsonString), &contract)
	if err != nil {
		fmt.Errorf("%s", err)
		return contract
	}
	return contract

}

func getAllContracts() []Contract {

	contracts := []Contract{}

	// Will send and get a response from the blockchain
	bodyText := postRequest("", "GetAllContracts")

	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")
	err := json.Unmarshal([]byte(jsonString), &contracts)

	if err != nil {
		fmt.Print("Got this string: %s", jsonString)
		fmt.Errorf("%s", err)
		return contracts
	}
	return contracts

}

// will format the strings into accepted format for Fablo rest api.
func combineStrings(strs ...string) string {
	var sb strings.Builder
	for i, s := range strs {
		sb.WriteString(`"` + s + `"`)
		if i < len(strs)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func prettifyAllContracts(contracts []Contract) {
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)

	// Set the table headers
	table.SetHeader([]string{"ID", "Status", "Notes", "Start Date", "End Date"})

	// Loop through the contracts and add each one to the table
	for _, contract := range contracts {
		table.Append([]string{
			contract.ID,
			contract.Status,
			contract.Notes,
			contract.StartDate,
			contract.EndDate,
		})
	}

	// Set the table style
	table.SetBorder(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Render the table
	table.Render()

}

func prettifyContract(contract Contract) {
	// In case the user enters a wrong ID
	if contract.ID == "" {
		println("No matching ID in the blockchain. Please try again.")
		return
	}
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)

	// Set the table headers
	table.SetHeader([]string{"Field", "Value"})

	// Append each field and value as a row in the table
	table.Append([]string{"ID", contract.ID})
	table.Append([]string{"Status", contract.Status})
	table.Append([]string{"Notes", contract.Notes})
	table.Append([]string{"Start Date", contract.StartDate})
	table.Append([]string{"End Date", contract.EndDate})
	table.Append([]string{"Extension Details", contract.ExtensionDetails})

	// Append employer details
	table.Append([]string{"Employer ID", contract.Employer.ID})
	table.Append([]string{"Employer Name", contract.Employer.Name})
	table.Append([]string{"Employer Address and Contact", contract.Employer.EmployerAC})
	table.Append([]string{"Employer Country", contract.Employer.Country})

	// Append employee details
	table.Append([]string{"Employee ID", contract.Employee.ID})
	table.Append([]string{"Employee Name", contract.Employee.Name})
	table.Append([]string{"Employee Address and Contact", contract.Employee.EmployeeAC})
	table.Append([]string{"Employee Country", contract.Employee.Country})

	// Append job details
	table.Append([]string{"Position", contract.Job.Position})
	table.Append([]string{"Level", contract.Job.Level})
	table.Append([]string{"Description", contract.Job.Description})

	// Append benefits details
	table.Append([]string{"Currency", contract.Benefits.Currency})
	table.Append([]string{"Salary", strconv.Itoa(contract.Benefits.Salary)})
	table.Append([]string{"Annual Increase", contract.Benefits.AnnualIncrease})
	table.Append([]string{"Annual Leave", contract.Benefits.AnnualLeave})
	table.Append([]string{"Housing", strconv.Itoa(contract.Benefits.Housing)})
	table.Append([]string{"Allowances", strconv.Itoa(contract.Benefits.Allowances)})
	table.Append([]string{"Other Benefits", contract.Benefits.OtherBenefits})

	// Append disputes details
	for _, dispute := range contract.Disputes {
		table.Append([]string{"Dispute ID", dispute.ID})
		table.Append([]string{"Dispute Status", dispute.Status})
		table.Append([]string{"Dispute Last Updated Date", dispute.LastUpdatedDate})
		table.Append([]string{"Dispute Content", dispute.Content})
		for _, response := range dispute.Responses {
			table.Append([]string{"Response ID", response.ID})
			table.Append([]string{"Response Last Updated Date", response.LastUpdatedDate})
			table.Append([]string{"Response Content", response.Content})
		}
	}

	// Setting the colors of the columns
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)
	// Set the table style
	table.SetBorder(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Render the table
	table.Render()
}

func prettifyTopContract(contract Contract) {
	// In case the user enters a wrong ID
	if contract.ID == "" {
		println("No matching ID in the blockchain. Please try again.")
		return
	}
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)

	// Set the table headers
	table.SetHeader([]string{"Field", "Value"})

	// Append each field and value as a row in the table
	table.Append([]string{"ID", contract.ID})
	table.Append([]string{"Status", contract.Status})
	table.Append([]string{"Notes", contract.Notes})
	table.Append([]string{"Start Date", contract.StartDate})
	table.Append([]string{"End Date", contract.EndDate})
	table.Append([]string{"Extension Details", contract.ExtensionDetails})
	table.Append([]string{"Salary", strconv.Itoa(contract.Benefits.Salary)})

	// Setting the colors of the columns
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)
	// Set the table style
	table.SetBorder(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Render the table
	table.Render()
}

func prettifyDispute(contract Contract) {
	// In case the user enters a wrong ID
	if contract.ID == "" {
		println("No matching ID in the blockchain. Please try again.")
		return
	}
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)

	// Set the table headers
	table.SetHeader([]string{"Field", "Value"})

	// Append each field and value as a row in the table

	for i := 0; i < len(contract.Disputes); i++ {
		dispute := contract.Disputes[i]

		table.Append([]string{"Dispute ID", dispute.ID})
		table.Append([]string{"Status", dispute.Status})
		table.Append([]string{"Last Updated", dispute.LastUpdatedDate})
		table.Append([]string{"Content", dispute.Content})
		for j := 0; j < len(dispute.Responses); j++ {
			response := contract.Disputes[i].Responses[j]
			table.Append([]string{"Response ID", response.ID})
			table.Append([]string{"Last Updated", response.LastUpdatedDate})
			table.Append([]string{"Content", response.Content})
		}
	}

	// Setting the colors of the columns
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)
	// Set the table style
	table.SetBorder(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Render the table
	table.Render()
}
