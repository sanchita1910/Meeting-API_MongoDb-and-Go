# meetingsapiappointy
A Meeting Scheduling API  
- An API made with Go and MongoDB
## Requirements
* MongoDB
* Go
## Features
Sticks to the given constraint and uses only standard lib
* It uses only the package/libraries from https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.4.0 and https://golang.org/pkg/ and nothing else (apart from mongodb) is needed !
## Runnning
You can access the API server at http://localhost:8083
## API

### /meetings
* POST : Schedule a new meeting

### /meeting/\<id here>
* GET : Get meeting using the id
 
### /meetings?starttime=\<start time>&endtime=\<endtime>  
**Please note this function uses "starttime" and "endtime" instead of "start" and "end" in the url.**
* GET: Returns an array of meetings in JSON format that are within the given certain time range

