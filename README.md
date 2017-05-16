# IISLogParser Project
Exercise in go to parse IIS logs, send to central server, and display the results 
in a web application.

# Dependencies
- Go -> Install go link https://golang.org/doc/install
- gRPC -> http://www.grpc.io/docs/quickstart/go.html
- 1.36Install Protocol Buffers v3 ->  MBprotoc-3.2.0-osx-x86_64.zip
- BigQuery locally?
    - https://cloud.google.com/bigquery/docs/reference/libraries#client-libraries-install-go
    - go get -u cloud.google.com/go/bigquery
- sample log files

## Testing gRPC
https://github.com/grpc/grpc-go/tree/master/examples


## Parser Name - DRD (diagnostic repair drone - dang that does not work... need something that pukes data...)

### Steps
- Runs once or more times a day
- Parses out data in new files, skips the files it has parsed
- Puts data into an object
- Posts or sends data to central server

### Configuration file

    * Directory of files
    * Customer Name
    * Server Name
    * Server IP
    
### History File
Since files follow a number pattern we will just save the last file we
imported

### Steps

    1. Define Schemas
    2. Import File
    3. Write summary data to .json file


### IIS Fields (Default)
https://stackoverflow.com/questions/11296698/understanding-iis-7-log-files
https://msdn.microsoft.com/en-us/library/windows/desktop/aa814385(v=vs.85).aspx

01. date (The date on which the activity occurred.)
02. time (The time, in coordinated universal time (UTC), at which the activity occurred.)
03. s-ip      (The IP address of the server on which the log file entry was generated.)   
04. cs-method (The requested verb, for example, a GET method.)
05. cs-uri-stem (The target of the verb, for example, Default.htm.)

06. cs-uri-query (The query, if any, that the client was trying to perform. A Universal Resource Identifier (URI) query is necessary only for dynamic pages.)

07. s-port (The server port number that is configured for the service.)
08. cs-username  (The name of the authenticated user that accessed the server. Anonymous users are indicated by a hyphen.)
09. c-ip (The IP address of the client that made the request.)
10. cs(User-Agent)  (The browser type that the client used.)
11. cs(Referer)  (The site that the user last visited. This site provided a link to the current site.)

12. sc-status    (The HTTP status code.)
13. sc-substatus  (The substatus error code.)
14. sc-win32-status  (The Windows status code.)
15. time-taken (The length of time that the action took, in milliseconds.)

Notes:
"::1" == 127.0.0.1
"-"  == Anonymous users are indicated by a hyphen.

#### Test 1

```bash

2015-11-04 22:22:31 ::1 POST /Home/ApplicationCategory_ByCategoryId - 80 - ::1 Mozilla/5.0+(Windows+NT+6.3;+WOW64;+Trident/7.0;+rv:11.0)+like+Gecko http://localhost/ 200 0 0 46

```


```json

{
    "date" : "2015-11-04", 
    "time" : "22:22:31" ,
    "s-ip" : "::1" ,                
    "cs-method" :"POST", 
    "cs-uri-stem" : "/Home/ApplicationCategory_ByCategoryId",           
    "cs-uri-query" : "-" ,
    "s-port" : "80" ,
    "cs-username" : "-",         
    "c-ip" : "::1" ,
    "cs(User-Agent)" : "Mozilla/5.0+(Windows+NT+6.3;+WOW64;+Trident/7.0;+rv:11.0)+like+Gecko", 
    "cs(Referer)" : "http://localhost/" ,
    "sc-status" : 200,
    "sc-substatus" : 0,
    "sc-win32-status" : 0, 
    "time-taken" : 46
}

```

#### Test 2

```bash

2015-11-04 22:04:49 ::1 GET /MDMVEE/JobAnalysis.aspx - 80 ARUBA\mdmadmin ::1 Mozilla/5.0+(Windows+NT+6.3;+WOW64;+Trident/7.0;+rv:11.0)+like+Gecko http://localhost/ 200 0 0 234

```

```json

{
    "date" : "2015-11-04", 
    "time" : "22:04:4931" ,
    "s-ip" : "::1" ,                
    "cs-method" :"GET", 
    "cs-uri-stem" : "/MDMVEE/JobAnalysis.aspx",           
    "cs-uri-query" : "-" ,
    "s-port" : "80" ,
    "cs-username" : "ARUBA\mdmadmin",         
    "c-ip" : "::1" ,
    "cs(User-Agent)" : "Mozilla/5.0+(Windows+NT+6.3;+WOW64;+Trident/7.0;+rv:11.0)+like+Gecko", 
    "cs(Referer)" : "http://localhost/" ,
    "sc-status" : 200,
    "sc-substatus" : 0,
    "sc-win32-status" : 0, 
    "time-taken" : 234
}

```

## Server - Moya
Server accpets connections from authenticated DRD client which listens
for incoming data.

Types of data streams
- logs
- licencse checks or validations

## Web UI - Aurora
The Aurora Application most effective means of gaining information from unwiling data.

# Database Options

AWS Dynamo Db (nosql)

Download local version for development and testing
https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html

Create a directory for dynamoDB, cd into the directory, download the db, extract the db, run the database
cn $ cd 
cn $ mkdir dynamoDB
cn $ cd dynamoDB
cn dynamoDB $ wget https://s3-us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.tar.gz
cn dynamoDB $ tar xopf dynamodb_local_latest.tar.gz 


--testing this one...
cn dynamoDB $ aws dynamodb list-tables --endpoint-url http://localhost:9090 

-port


cn dynamoDB $ java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
Initializing DynamoDB Local with the following configuration:
Port:	8000
InMemory:	false
DbPath:	null
SharedDb:	true
shouldDelayTransientStatuses:	false
CorsParams:	*


## Database Schema(s)

We need to be multi-tenat and have counts by Customer and 
total counts across all instances.

    * Customer Users
        * Individual
        * Last Login Date
    * Customer Users
        * Hourly
        * Day
        * Week
        * Month
        * Date Range
        * Count
    * Customer Pages Viewed
        * Page Name
        * Hourly
        * Day
        * Week
        * Month
        * Count

# TODO

    * go web app to display the results
    * all customers in one app
    * AWS or GCP?