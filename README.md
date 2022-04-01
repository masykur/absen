# keico
Time Attendance and access control machine TCP communication writen in Go

Implemented machines are below
## Keico SF3000

### Features
#### General
- [x] Obtain product code
- [x] Obtain product serial number
- [ ] Obtain device info
- [ ] Obtain detail device info 
- [ ] Obtain device status
- [ ] Enable device
- [ ] Power off device
- [ ] Upgrade firmware
#### Date and time
- [x] Retrieve current date and time from machine
- [x] Set current date and time to machine
#### User data
- [x] Get enrolled user information including card number, card facility code and fingerprint templates from machine
- [x] Get number enrolled users from machine
- [x] Get list of enrolled users from machine
- [x] Enroll user and it information including card number, card facility code and fingerprint templates to machine
- [ ] Delete enrolled user from machine
- [ ] Delete all enrolled users from machine
- [ ] Modify user privilage
#### Log data
- [ ] Read general log data
- [ ] Read super log data
- [ ] Pull general log data
- [ ] Pull super log data 
- [ ] Clear keeper data
- [ ] Delete general log data
- [ ] Delete super log data
- [ ] Delete all general log data
- [ ] Delete all super log data


## RECO RAC2000, AC2200PC
### Features
#### Date and time
- [x] Retrieve current date and time from machine
- [ ] Set current date and time to machine
#### Card
- [ ] Add / register new card to machine 
- [ ] Delete / unregister card from machine
- [ ] Add / register visitor card for certain periods of time
- [ ] Delete/ unregister visitor card
#### Log data
- [ ] Read log data
- [ ] Clear log data 


> **Note:** _The software is reverse engineered by sniffing network package between official software and machine. It is developed without any official references or documentations from hardware maker._ 
