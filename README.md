# appointy-project

A set of rest API's build using GoLang for scheduling and querying meetings.

## Installation

Clone into the repo
```bash 
git clone https://github.com/yashjjw/appointy-project
```
Install the go packages used in the project
```bash
go get -d ./...
```
Run the main.go file
```bash
go run main/main.go
```

## Usage

1. The API for scheduling meetings , POST `/meetings`  The format of JSON is as follows
```
{
    "title": "meeting",
    "participants": [
        {
            "name": "John Doe",
            "email": "John@Doe.com",
            "rsvp": "Yes"
        }
    ],
    "startTime": INT, //time input in UNIX format 
    "endTime": INT
}

```

This API will check for overlapping meeting time of participants with RSPV = YES and wont schedule a meeting if there is an overlap. 
If there is no overlap then the meeting will be scheduled and the details will be returned as JSON

2. The API for getting meetings by ID, GET `/meeting/<id> `, This route will query the database for the meeting id and return all the details of the meeting in JSON format.

3. The API for getting meetings in a time frame , GET `/meetings?start=<startTimestamp>&end=<endTimestamp>` where the time will be given in UNIX format.
This will return an array of all the meetings in the given time frame.

4. The API for getting all the meetings of a participant , GET `/meetings?participant=<email>` , This will return an array of all the meetings with the given participant. 
