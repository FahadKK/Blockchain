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
	fmt.Println("Please choose one of the following methods by entering its number:")
	fmt.Println(" 1. Create Contract (Simple)")
	fmt.Println(" 2. Approve Contract")
	fmt.Println(" 3. Update Contract ")
	fmt.Println(" 4. Extend Contract")
	fmt.Println(" 5. Terminate Contract")
	fmt.Println(" 6. Issue Dispute")
	fmt.Println(" 7. Update Dispute")
	fmt.Println(" 8. Close Dispute")
	fmt.Println(" 9. Respond to Dispute")
	fmt.Println("10. Create Contract")
	fmt.Println("11. Read Contract")
	fmt.Println("12. View Employee History")
	fmt.Println("13. View Employer History")
	fmt.Println("14. Get All Contracts")

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
		fmt.Printf("Invalid input. Please enter a number between 1 and 15.   %s \n", err)
		return 0
	}

	if methodNumber < 1 || methodNumber > 15 {
		fmt.Printf("Invalid input. Please enter a number between 1 and 115. \n")
		return 0
	}
	return methodNumber

}

func main() {
	for j := 0; j < 14; j++ { // the program will loop for 14 times only
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

			createContractSimple()
		case 2:
			approveContract()
		case 3:
			updateContract()
		case 4:
			extendContract()
		case 5:
			terminateContract()
		case 6:
			issueDispute()
		case 7:
			updateDispute() // Not implemented yet.
		case 8:
			closeDispute()
		case 9:
			respondToDispute()
		case 10:
			createContract()
		case 11:
			prettifyContract(readContract())
		case 12:
			viewEmployeeHistory()
		case 13:
			viewEmployerHistory()
		case 14:
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
	println(jsonString)

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
	println(jsonString)

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
	println(jsonString)

}

func respondToDispute() {
	// Taking all the need inputs from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	fmt.Print("Enter Dispute ID: ")
	DisputeID, err2 := reader.ReadString('\n')

	fmt.Print("Enter Response ID: ")
	RID, err4 := reader.ReadString('\n')

	fmt.Print("Enter the content of your response: ")
	Content, err3 := reader.ReadString('\n')

	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		fmt.Printf("Could not read string \n")
	}

	combinedInputs := combineStrings(ID, DisputeID, RID, Content)
	combinedInputs = strings.ReplaceAll(combinedInputs, "\n", "")

	// Will send and get a response from the blockchain
	bodyText := postRequest(combinedInputs, "RespondToDispute")
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	println(jsonString)
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

	println(jsonString)
}

func issueDispute() {
	// Taking all the need inputs from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	ID, err := reader.ReadString('\n')
	fmt.Print("Enter Dispute ID: ")
	DisputeID, err2 := reader.ReadString('\n')

	fmt.Print("Enter the content of your dispute: ")
	Content, err3 := reader.ReadString('\n')

	if err != nil || err2 != nil || err3 != nil {
		fmt.Printf("Could not read string \n")
	}

	combinedInputs := combineStrings(ID, DisputeID, Content)
	combinedInputs = strings.ReplaceAll(combinedInputs, "\n", "")

	// Will send and get a response from the blockchain
	bodyText := postRequest(combinedInputs, "IssueDispute")
	jsonString := string(bodyText)
	jsonString = strings.TrimPrefix(string(bodyText), "{\"response\":")
	jsonString = strings.TrimSuffix(jsonString, "}")

	println(jsonString)
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

// creating contract the simple way.
func createContractSimple() {
	client := &http.Client{}

	// Taking input from user.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Contract ID: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read string %s \n", err)
		return
	}

	input = strings.ReplaceAll(input, "\n", "")

	// We don't use the post method because here we are sending json string.
	// TODO: Create a method to easily embede inputs here.
	var data = strings.NewReader(`{"method": "HandleAddContract",
"args": ["{'ID': '` + input + `','Status': 'Active','Notes': 'N/A','Start date': '01/01/2022','End date': '08/08/2023','Extension details': 'N/A','Employer': {        'ID': 'Comp-1',        'Name': 'CompanyA',       'Employer address and contact details': 'First st,Riyadh1234','Country': 'Saudi Arabia'      },'Employee': {        'ID': '441101772',        'Name': 'JohnDoe',        'Employee address and contact details': 'Second st,New Delhi 3342',        'Country': 'India'      },'Job': {        'Position': 'Developer',        'Level': 'Senior',        'Description': 'Manage teams of junior developers'      },'Benefits': {        'Currency': 'SAR',        'Salary': 10000,        'Annual increase': '3-7',        'Annual leave': '30 days','Housing': 2000,        'Allowances': 1500,'Other benefits': 'Schooling for children and yearly tickets'      }, 'Disputes': [        {          'ID': 'D123',          'Status': 'Closed',          'Last updated date': '01/10/2023',          'Content': 'Employer did not provide the travel tickets for my annual leave in 2022','Response': [            {              'ID': 'Res123',              'Last updated date': '01/12/2023',              'Content': 'The employee was compensated'            }          ]}      ]    }"  ]}`) //

	req, err := http.NewRequest("POST", "http://localhost:8801/invoke/my-channel/chaincode1", data)
	if err != nil {
		log.Fatal(err)
	}
	token := getToken()
	req.Header.Set("Authorization", "Bearer "+token)
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

	fmt.Println(string(bodyText))
}

func updateContract() {
	fmt.Println("Updating a contract from scratch involes alot of vairables.")
	fmt.Println("This is why I need you to go to the file Main/contract.json. ")
	fmt.Println("In that file you will find a valid contract. Please change whatever value you desire.")
	fmt.Print("Don't forget to change cotnract ID.")
	fmt.Println("Kindly don't mess with the structure of the contract.")
	fmt.Println("After changing the values of the contract save the json file and press ENTER.")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("")
	_, err := reader.ReadString('\n')

	// Read the contents of the file
	input, err := ioutil.ReadFile("Main/contract.json")
	if err != nil {
		fmt.Println("Error")
		panic(err)

	}

	// Convert the byte slice to a string and remove newline characters
	jsonData := string(input)
	jsonData = strings.ReplaceAll(jsonData, "\n", "")
	jsonData = strings.ReplaceAll(jsonData, "\"", "'")
	jsonData = combineStrings(jsonData)

	bodyText := postRequest(jsonData, "UpdateContract")
	fmt.Println(bodyText)
}

func createContract() {

	fmt.Println("Creating a contract from scratch involes alot of vairables.")
	fmt.Println("This is why I need you to go to the file Main/contract.json. ")
	fmt.Println("In that file you will find a valid contract. Please change whatever value you desire.")
	fmt.Println("Kindly don't mess with the structure of the contract.")
	fmt.Println("After changing the values of the contract save the json file and press ENTER.")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("")
	_, err := reader.ReadString('\n')

	// Read the contents of the file
	input, err := ioutil.ReadFile("Main/contract.json")
	if err != nil {
		fmt.Println("Error")
		panic(err)

	}

	// Convert the byte slice to a string and remove newline characters
	jsonData := string(input)
	jsonData = strings.ReplaceAll(jsonData, "\n", "")
	jsonData = strings.ReplaceAll(jsonData, "\"", "'")
	jsonData = combineStrings(jsonData)

	bodyText := postRequest(jsonData, "HandleAddContract")
	fmt.Println(bodyText)
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
