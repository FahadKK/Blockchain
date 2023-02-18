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

This installation guide will assume that you have installed Fablo globally. <br>


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

## Fablo Installation | Setting Up The Environment
-----
<br>

If you need help setting up Fablo, this guide will help you. <br>
Be aware that we will assume that this is a fresh installation of Ubuntu OS. 
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











## License
------
<br>

(Insert License)
