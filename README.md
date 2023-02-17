#  Managing Overseas Employment Contracts with Blockchain 
Using Blockchain to Manage Overseas Employment Contracts \n
(Insert Name) is a proof of concept program that will demonstrate the solution provided in the following paper: \n
ww.google.com


## Prerequisites 
--------------
This program has been tested mainly on the Linux operating system.
 So, I recommend using ubuntu to try this application
- Golang
- curl
- Docker
-  [Fablo](https://github.com/hyperledger-labs/fablo)


## Installation
---------------
This installation guide will assume that you have installed Fablo globally. 

Inside BlockChain folder run this command :
```
    sudo fablo recreate
```

When the network finshes setting up run this command: 
```
go run Main/Main.go
```

Now interact with the system as much as you want.
* sudo Fablo prune will shut done the network, including all stored information.
* sudo Fablo recreate will reset the network.

## Fablo Installation | Setting Up The Environment
-----
If you need help setting up Fablo, this guide will help you.
Do note that we will assume that this is a fresh installation of Ubuntu OS. 
```
sudo apt-get update
sudo apt --yes install software-properties-common
sudo apt install --yes golang
su
do apt-get install -y jq
sudo apt-get install -y curl
sudo apt install docker --yes
sudo apt install docker.io --yes
sudo apt install docker-compose --yes
```

After setting up the enviorment we will install Fablo
```
sudo curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o /usr/local/bin/fablo && sudo chmod +x /usr/local/bin/fablo
```
Credits to @(Insert Email)


## License
------
(Insert Licesne)
