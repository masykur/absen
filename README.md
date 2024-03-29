# Absen
Time Attendance and access control machine TCP communication writen in Go

> **Note:** _The software is reverse engineered by sniffing network package between official software and machine. It is developed without any official references or documentations from hardware maker._ 

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
- [x] Read general log data
- [ ] Read super log data
- [ ] Pull general log data
- [ ] Pull super log data 
- [ ] Clear keeper data
- [ ] Delete general log data
- [ ] Delete super log data
- [ ] Delete all general log data
- [ ] Delete all super log data

## Keico SF3500 (Face ID)

### Features
#### Client
- [x] Obtain device info
- [x] Obtain user info
- [x] Get list of registered users
- [x] Get log data
#### Server
- [x] Listen incoming log data
- [x] Listen incoming enrolled data
- [x] Save log data into database

## RECO RAC2000, AC2200PC
### Features
#### Date and time
- [x] Retrieve current date and time from machine
- [x] Set current date and time to machine
#### Card
- [x] Get list of registered cards from machine
- [x] Add / register new card to machine 
- [x] Delete / unregister card from machine
- [ ] Add / register visitor card for certain periods of time
- [ ] Delete/ unregister visitor card
#### Log data
- [x] Fetch log data