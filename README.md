# GoToDo - A multiplatform task scheduler

This project intends to use a desktop, rasperry pi and a CYD-Esp32. Though the pi could be ommitted to run the server on the desktop too.

I undertook this mini-project as a learning tool for starting out with GO and learning the language syntax and fundamentals.

Setup the Firebase as follows:
  - DocumentID (item) = e.g Task 1
  - Name = string
  - Due = timestamp
  - InProgress = boolean

How to use:
CLI: 
  - go run . add --name "{Task Name}" --due "{Due Date FORMAT: YYYY-MM-DD OR YYYY-MM-DD HH:MM}"
  - go run . list (optional) --all (NOTE: all flag shows completed tasks)
  - go run . done --docID "{TASK XYZ}"
Server:
  - I set up a sytemd service on the pi as follows:
[Unit]
Description=Todo API Firestore Server
After=network-online.target

[Service]
User=todoapi
Group=todoapi
WorkingDirectory=/opt/todoapi
ExecStart=/opt/todoapi/todoapi
Environment=GOOGLE_APPLICATION_CREDENTIALS=/opt/todoapi/(FIREBASE_CREDENTIALS).json // REPLACE THESE WITH YOUR CREDENTIALS
Environment=PROJECT_ID={FIREBASE_PROJECTID} // REPLACE THESE WITH YOUR CREDENTIALS
Environment=PORT=8080
Restart=always
RestartSec=2

NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true

CYD:
  - Insert your network SSID, Password and target IP
  - Flash the device and it will now automatically update with tasks


## Example of CYD display
[Example of ToDo on CYD](CYD-Example.png)
