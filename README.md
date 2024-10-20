# Event store experiment

## What it is
This Go-based Eventstore is an experimental implementation of a minimalist event store service. Its primary use case is appending events to an Azure Append Blob. The service exposes two main REST API handlers:

- An idempotent PUT /events endpoint for adding new events
- A GET /events endpoint for retrieving either all stored events or a specified range
- Additionally, a GET /eventstream endpoint is planned for future implementation.

The idempotency of the PUT operation is not directly addressed: Each call to the PUT events method results in an event record being written. To solve this issue, we have several options. One approach involves implementing additional infrastructure elements, such as a distributed Bloom filter, to check whether an event has already been stored. Alternatively, we could build a distributed/synchronized Cuckoo filter using the Raft protocol. However, the complexity of such solutions is somewhat excessive for a playground project.

A simpler solution is to have each Eventstore instance maintain a Bloom filter when reading events. If an event is not present in the Bloom filter, it is passed to the client, and the message ID is added to the filter. Since we use Azure Append-Only Blob Storage, duplicates cannot be removed later (theoretically, we could create a copy, filter out duplicate events, and write a new record). On the other hand, each duplicated record tells a story in time...

## Run the service
To effectively experiment with this system, you'll need an Azure Storage Account with appropriate permissions and a Storage Container. Regarding access rights:

For local development:
Log in to your Azure tenant using 'az login' with a user principal. The user principal must have "Storage Blob Data Contributor" permissions on the storage account.

For a cloud native infrastructure:
A managed identity should be employed.
The code is designed to work seamlessly with both authentication methods out of the box.

- Don't use this in a production environment
- Be sure you have go installed
- Be sure you have the latest azure cli installed (when using the azure path)
- Clone the repo
- Use go >= go 1.23
- run the server:
```bash
$ cd src
$ go mod tidy
$ az login
$ go run cmd/server/main.go \
  -storage-type=azure \
  -azure-account-name=<your storage account name> \
  -azure-container-name=<your storage container name> \
  -azure-blob-name=<event-log-name, defaults to event.log>
```

## Help: 
```
$ go run cmd/server/main.go -h
  -azure-account-name string
    	Azure storage account name
  -azure-blob-name string
    	Azure Blob name
  -azure-container-name string
    	Azure Blob container name
  -estimated-event-count uint
    	Server address (default 1000000)
  -file-path string
    	File path for local file storage (default "events.log")
  -h	Display help
  -help
    	Display help
  -server-address string
    	Server address (default ":8080")
  -storage-type string
    	Type of storage to use (memory, file, aws, azure) (default "memory")```
Instead of the azure way, you can also configure it a step simpler with local storage or even local memory. The ephemeral way to hell.
```


## Test the server with curl

Write two events (you can also write one single event with /event endpoint)
```
$ curl \
  --location \
  --request PUT 'http://localhost:8080/events' \
  --header 'Content-Type: application/json' \
  --data '[{"message_id":"4e976abd-f891-456e-9ad6-4e0a5d66dc5a","timestamp": "1730040745283000000","event_type":"USER_CREATED","payload": "{\"userId\":4729,\"username\":\"leonwvo\",\"firstName\":\"Leon\",\"lastName\":\"Fischer\",\"email\":\"leonwvo@example.com\"}"},{"message_id": "044313bd-5687-4509-a2a2-f85a49f8b00f","timestamp": "1730040765965000000","event_type": "USER_CREATED","payload": "{\"userId\":6112,\"username\":\"annawdw\",\"firstName\":\"Anna\",\"lastName\":\"Weber\"\"email\":\"annawdw@example.com\"}"}]'
```

Read all events since first appended event:
```
$ curl --location 'http://localhost:8080/events'
[
    {
        "message_id": "4e976abd-f891-456e-9ad6-4e0a5d66dc5a",
        "timestamp": "1730040745283000000",
        "event_type": "USER_CREATED",
        "payload": "{\"userId\":4729,\"username\":\"leonwvo\",\"firstName\":\"Leon\",\"lastName\":\"Fischer\",\"email\":\"leonwvo@example.com\"}"
    },
    {
        "message_id": "044313bd-5687-4509-a2a2-f85a49f8b00f",
        "timestamp": "1730040765965000000",
        "event_type": "USER_CREATED",
        "payload": "{\"userId\":6112,\"username\":\"annawdw\",\"firstName\":\"Anna\",\"lastName\":\"Weber\"\"email\":\"annawdw@example.com\"}"
    }
]
```

You also can do simple paged requests with query start=0&limit=10.