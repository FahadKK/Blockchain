#  Managing Overseas Employment Contracts with Blockchain
Using Blockchain to Manage Overseas Employment Contracts  <br>
<br>


(Insert Name) is a proof of concept program that will demonstrate the solution provided in the following paper: <br>
ww.google.com



## Prerequisites 
--------------


<br>
This program has been tested mainly on the Linux operating system. <br>
 So, I recommend using Ubuntu to try this application <br>
 <br>

In addition, To Ubuntu, you will need to install the following: 

- Golang
- curl
- Docker
-  [Fablo](https://github.com/hyperledger-labs/fablo)

<br>

## Installation
---------------
<br>



Inside BlockChain folder run this command :
``` 
    sudo fablo recreate
```

When the network finishes setting up run this command: 
``` 
    go run Main/Main.go
```

Now interact with the system as much as you want.
* sudo Fablo prune will shut done the network, including all stored information.
* sudo Fablo recreate will reset the network.

<br>

## Fablo Installation | Setting Up The Environment
-----
<br>

If you need help setting up Fablo, this guide will help you. <br>
Be aware that we will assume this is a fresh installation of Ubuntu OS.
```
sudo apt-get update
sudo apt --yes install software-properties-common
sudo apt install --yes golang
sudo apt-get install -y jq
sudo apt-get install -y curl
sudo apt install docker --yes
sudo apt install docker.io --yes
sudo apt install docker-compose --yes
```

After setting up the environment we will install Fablo
```
sudo curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o /usr/local/bin/fablo && sudo chmod +x /usr/local/bin/fablo
```
Credits to @(Insert Email)


<br>


## Examples
-----------

We will start with an empty blockchain and populate it.
<br>



### **Example 1 - Starting Up**
```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 
```
<br>


### **Example 2 - Create Contract (Simple)**


Creating a contract involves many variables, so I made a method that will streamline the process. <br>
This method will only require a contract id to run. The rest of the variables will be filled by the system. <br>

```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 1
Enter Contract ID: C-1234
```
**Output**
```
{"response":true}
```
<br>







### **Example 3 - Read Contract**
```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 11
Enter Contract ID: C-1234
```

**Output**
```
+------------------------------+--------------------------------+
| FIELD                        | VALUE                          |
+------------------------------+--------------------------------+
| ID                           | C-1234                         |
| Status                       | Pending                        |
| Notes                        | N/A                            |
| Start Date                   | 01/01/2022                     |
| End Date                     | 08/08/2023                     |
| Extension Details            | N/A                            |
| Employer ID                  | Comp-1                         |
| Employer Name                | CompanyA                       |
| Employer Address and Contact | First st,Riyadh1234            |
| Employer Country             | Saudi Arabia                   |
| Employee ID                  | 441101772                      |
| Employee Name                | JohnDoe                        |
| Employee Address and Contact | Second st,New Delhi 3342       |
| Employee Country             | India                          |
| Position                     | Developer                      |
| Level                        | Senior                         |
| Description                  | Manage teams of junior         |
|                              | developers                     |
| Currency                     | SAR                            |
| Salary                       | 10000                          |
| Annual Increase              | 3-7%                           |
| Annual Leave                 | 30 days                        |
| Housing                      | 0                              |
| Allowances                   | 0                              |
| Other Benefits               | Schooling for children and     |
|                              | yearly tickets                 |
+------------------------------+--------------------------------+
```
<br>


### **Example 4 - Approve Contract**

```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 2
Enter Contract ID: C-1234
```
**Output** <br>
<font size = "2" > We will read the contract to see changes. The contract status will become Active.  </font>

```
+------------------------------+--------------------------------+
| FIELD                        | VALUE                          |
+------------------------------+--------------------------------+
| ID                           | C-1234                         |
| Status                       | Active                         |
| Notes                        | N/A                            |
| Start Date                   | 01/01/2022                     |
| End Date                     | 08/08/2023                     |
| Extension Details            | N/A                            |
| Employer ID                  | Comp-1                         |
| Employer Name                | CompanyA                       |
| Employer Address and Contact | First st,Riyadh1234            |
| Employer Country             | Saudi Arabia                   |
| Employee ID                  | 441101772                      |
| Employee Name                | JohnDoe                        |
| Employee Address and Contact | Second st,New Delhi 3342       |
| Employee Country             | India                          |
| Position                     | Developer                      |
| Level                        | Senior                         |
| Description                  | Manage teams of junior         |
|                              | developers                     |
| Currency                     | SAR                            |
| Salary                       | 10000                          |
| Annual Increase              | 3-7%                           |
| Annual Leave                 | 30 days                        |
| Housing                      | 0                              |
| Allowances                   | 0                              |
| Other Benefits               | Schooling for children and     |
|                              | yearly tickets                 |
+------------------------------+--------------------------------+
```




### **Example 5 - issue Dispute**
```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 6
Enter Contract ID: C-1234
Enter Dispute ID: D-1234
Enter the content of your dispute: Never Received the promised schooling for children. 
```

**Output** <br>

<font size = "2" > We will read the contract to see changes. </font>
```
+------------------------------+--------------------------------+
| FIELD                        | VALUE                          |
+------------------------------+--------------------------------+
| ID                           | C-1234                         |
| Status                       | Active                         |
| Notes                        | N/A                            |
| Start Date                   | 01/01/2022                     |
| End Date                     | 08/08/2023                     |
| Extension Details            | N/A                            |
| Employer ID                  | Comp-1                         |
| Employer Name                | CompanyA                       |
| Employer Address and Contact | First st,Riyadh1234            |
| Employer Country             | Saudi Arabia                   |
| Employee ID                  | 441101772                      |
| Employee Name                | JohnDoe                        |
| Employee Address and Contact | Second st,New Delhi 3342       |
| Employee Country             | India                          |
| Position                     | Developer                      |
| Level                        | Senior                         |
| Description                  | Manage teams of junior         |
|                              | developers                     |
| Currency                     | SAR                            |
| Salary                       | 10000                          |
| Annual Increase              | 3-7%                           |
| Annual Leave                 | 30 days                        |
| Housing                      | 0                              |
| Allowances                   | 0                              |
| Other Benefits               | Schooling for children and     |
|                              | yearly tickets                 |
| Dispute ID                   | D-1234                         |
| Dispute Status               | Active                         |
| Dispute Last Updated Date    | 02/18/2023                     |
| Dispute Content              | Never Received the promised    |
|                              | schooling for children.        |
+------------------------------+--------------------------------+
```
As you can see, the dispute has been submitted correctly.

<br>

### **Example 6 - Responding to a dispute**

```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 9
Enter Contract ID: C-1234
Enter Dispute ID: D-1234
Enter Response ID: R-1234
Enter the content of your response: Compensation for tuition will be paid immediately.   
```

**Output** <br>
<font size = "2" > We will read the contract to see changes </font>

```
+------------------------------+--------------------------------+
| FIELD                        | VALUE                          |
+------------------------------+--------------------------------+
| ID                           | C-1234                         |
| Status                       | Active                         |
| Notes                        | N/A                            |
| Start Date                   | 01/01/2022                     |
| End Date                     | 08/08/2023                     |
| Extension Details            | N/A                            |
| Employer ID                  | Comp-1                         |
| Employer Name                | CompanyA                       |
| Employer Address and Contact | First st,Riyadh1234            |
| Employer Country             | Saudi Arabia                   |
| Employee ID                  | 441101772                      |
| Employee Name                | JohnDoe                        |
| Employee Address and Contact | Second st,New Delhi 3342       |
| Employee Country             | India                          |
| Position                     | Developer                      |
| Level                        | Senior                         |
| Description                  | Manage teams of junior         |
|                              | developers                     |
| Currency                     | SAR                            |
| Salary                       | 10000                          |
| Annual Increase              | 3-7%                           |
| Annual Leave                 | 30 days                        |
| Housing                      | 0                              |
| Allowances                   | 0                              |
| Other Benefits               | Schooling for children and     |
|                              | yearly tickets                 |
| Dispute ID                   | D-1234                         |
| Dispute Status               | Active                         |
| Dispute Last Updated Date    | 02/18/2023                     |
| Dispute Content              | Never Received the promised    |
|                              | schooling for children.        |
| Response ID                  | R-1234                         |
| Response Last Updated Date   | 02/18/2023                     |
| Response Content             | Compensation for tuition will  |
|                              | be paid immediately.           |
+------------------------------+--------------------------------+
```
<br>

### **Example 7 - Create Contract**

```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 10
Creating a contract from scratch involes alot of vairables.
This is why I need you to go to the file Main/contract.json. 
In that file you will find a valid contract. Please change whatever value you desire.
Kindly don't mess with the structure of the contract.
After changing the values of the contract save the json file and press ENTER.
```
<font size = "4" > [contract.json](Main/contract.json) </font>


**Output** <br>
<font size = "2" > We will read the contract to see changes id = 10 </font>

```
+------------------------------+--------------------------------+
| FIELD                        | VALUE                          |
+------------------------------+--------------------------------+
| ID                           | 10                             |
| Status                       | Pending                        |
| Notes                        | N/A                            |
| Start Date                   | 01/02/2022                     |
| End Date                     | 05/05/2023                     |
| Extension Details            | N/A                            |
| Employer ID                  | Comp-1                         |
| Employer Name                | Company A                      |
| Employer Address and Contact | First st, Riyadh 12345         |
| Employer Country             | Saudi Arabia                   |
| Employee ID                  | 441101772                      |
| Employee Name                | John Doe                       |
| Employee Address and Contact | Second st, New Delhi 3342      |
| Employee Country             | India                          |
| Position                     | Developer                      |
| Level                        | Senior                         |
| Description                  | Manage teams of junior         |
|                              | developers                     |
| Currency                     | SAR                            |
| Salary                       | 10000                          |
| Annual Increase              | 3-7%                           |
| Annual Leave                 | 30 days                        |
| Housing                      | 0                              |
| Allowances                   | 0                              |
| Other Benefits               | Schooling for children and     |
|                              | yearly tickets                 |
+------------------------------+--------------------------------+
```

<br>

### **Example 8 - View Employer Statistic**

```
|-----------------------------|--|-----------------------|--|---------------------------|--|
|           Options           |  |        Options        |  |          Options          |  |
| 1. Create Contract (Simple) |  | 6. Issue Dispute      |  | 11. Read Contract         |  |
| 2. Approve Contract         |  | 7. Update Dispute     |  | 12. View Employee History |  |
| 3. Update Contract          |  | 8. Close Dispute      |  | 13. View Employer History |  |
| 4. Extend Contract          |  | 9. Respond to Dispute |  | 14. Get All Contracts     |  |
| 5. Terminate Contract       |  | 10. Create Contract   |  |                           |  |
|-----------------------------|--|-----------------------|--|---------------------------|--|
Please choose one of the above options by entering its number:
Enter your choice: 13
Enter EmployerID: Comp-1
```

**Output** <br>

```
+-----------+----------------------+------------------+-------------------+----------+---------------+-----------------+
| CONTRACTS | TERMINATED CONTRACTS | ACTIVE CONTRACTS | PENDING CONTRACTS | DISPUTES | OPEN DISPUTES | CLOSED DISPUTES |
+-----------+----------------------+------------------+-------------------+----------+---------------+-----------------+
| 2         | 0                    | 1                | 1                 | 1        | 1             | 0               |
+-----------+----------------------+------------------+-------------------+----------+---------------+-----------------+
```


## License
------
<br>

(Insert License)
