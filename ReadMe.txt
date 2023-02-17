curl --request POST   --url http://localhost:8801/user/enroll   --header 'Authorization: Bearer '   --data '{"id": "admin", "secret": "adminpw"}'

curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"InitLedger\",\"args\": []}"

curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"ReadContract\",
\"args\": [\"C-1234\"]}"


curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"DummyRead\",
\"args\": [\"Res123\"]}"


curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' \
-H 'Content-Type: application/json' --data "{\"args\": [ \"1234\",\"D-12345\",\"Did not pay me \"],
\"method\": \"IssueDispute\"}"

curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' \
-H 'Content-Type: application/json' --data "{\"args\": [ \"C-1234\",\"D-1234\",\"He paid me less\"],
\"method\": \"UpdateDispute\"}"


curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' \
-H 'Content-Type: application/json' --data "{\"args\": [ \"C-1234\",\"D-1234\",\"R-421\",\"Sorry we will pay him.\" ],
\"method\": \"RespondToDispute\"}"


curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"ExtendContract\",
\"args\": [\"C-1234\",\"02/06/2025\"]}"




curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"TerminateContract\",
\"args\": [\"C-1234\"]}"

curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"CloseDispute\",
\"args\": [\"1234\", \"D-12345\"]}"




curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"ReadContract\",
\"args\": [\"C-1234\"]}"

curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"ApproveContract\",
\"args\": [\"C-1234\"]}"



//32313131
curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin'  --data "{\"method\": \"AddContract\",
\"args\": [ \"1234\", \"C-no\", \"C-Sdat\",
\"C-Edate\", \"C-Ext\", \"Em-1423\", \"STC\", \"011-44321\", \"SaudiArabia\", \"441101772\", \"Fahd\", \"05079\", \"Saudi Arabia\", \"C-level Officer\", \"High\", \"Oversee Tech department\",
\"SAR\", \"30000\", \"3%\", \"60 Days\", \"1200\",\"444\", \"NA\"]}"







curl --request POST --url http://localhost:8801/invoke/my-channel/chaincode1 --header 'Authorization: Bearer b76207c0-ae5d-11ed-9067-4b24ea7f8dc0-admin' --data "{\"method\": \"ViewEmployeeHistory\",
\"args\": [\"441101772\"]}"
