# Conduit Connector Stripe

### General
The Stripe connector is one of [Conduit](https://github.com/ConduitIO/conduit) plugins. It provides a source Stripe connector.

### Prerequisites
- [Go](https://go.dev/) 1.18
- (optional) [golangci-lint](https://github.com/golangci/golangci-lint) 1.45.2
- (optional) [mock](https://github.com/golang/mock) 1.6.0

### Configuration
The config passed to `Configure` can contain the following fields:

| name            | description                                                                                                 | required | example                      |
|-----------------|-------------------------------------------------------------------------------------------------------------|----------|------------------------------|
| `secretKey`     | Stripe [secret key](https://dashboard.stripe.com/apikeys).                                                  | yes      | "sk_51Kr0QrJit566F2YtZAwMlh" |
| `resourceName`  | The name of Stripe resource. A list of supported resources can be found [here](models/resources/README.md). | yes      | "plan"                       |

### How to build it
Run `make build`.

### Testing
Run `make test`.

### Stripe Source
The `Configure` method parses the configuration and validates them.

The `Open` method parses the current position, initializes an [HTTP client](https://github.com/hashicorp/go-retryablehttp), and initializes Snapshot (only if in the position IteratorType equals Snapshot) and CDC iterators.

The `Read` method calls the method `Next` of the current iterator and returns the next record.

The `Teardown` method calls the method `Close` of the http client, which calls `CloseIdleConnections` method of the [net/http](https://pkg.go.dev/net/http) package.

### Snapshot
The system retrieves data from the list of objects of a defined Stripe resource (e.g. resource [plan](https://stripe.com/docs/api/plans/list)).

The system makes the first request to get a list of resource objects and appends the data to the resulting slice.

The data in the resulting slice is ready to be returned by the `Read` method record by record.

Each time the `Read` method is called, the system sets the current record ID to the `Cursor` to the position.

After all the data from the resulting slice has been read, the system makes the next request to Stripe with `starting_after` parameter, which is `Cursor` from the position.

The system stops making requests when the field `has_next` equals false. Then, the system updates `IteratorType` field in the position with `cdc` value and clears the `Cursor` field.

**Note:** The Snapshot iterator creates a copy of the data, which is sorted by date of creation in descending order.

### Change Data Capture
The system retrieves data from [Stripe events](https://api.stripe.com/v1/events) of the defined resource.

Because of this, receiving data is divided by requests to Stripe with different shift parameters:
- starting_after
- ending_before

Each time the `Read` method is called, the system updates `Index` field to the position. `Cursor` field updates when the last resulting slice element is read.

#### starting_after
The system makes requests in the loop to get all the data since the first start of the pipeline `CreatedAt`.

After each request, the system adds the received data to the final slice and makes the next request with the starting_after parameter, which is the event ID of the last element from the response.

The system stops making requests when the field `has_next` equals false.

Then the system reverses the whole resulting slice and sets the Cursor ID of the "freshest" event.

The data in the resulting slice is ready to be returned by the `Read` method record by record.

#### ending_before
After all the data from the resulting slice has been read, the system takes the next batch of data from Stripe.

The system makes requests in the loop to get all data starting from the position of the freshest event element (Cursor) using the ending_before parameter.

After each request, the system sets a new value for the next request before_ending (the ID of the first event element from the response), reverses the response, and appends the data to the resulting slice.

The system stops making requests when the field `has_next` equals false.

The data in the resulting slice is ready to be returned by the `Read` method record by record.

All the following data are taken by the [ending_before](#ending_before) script.

### Position
Position is a JSON object with the following fields:

| name            | type    | description                                                                               | 
|-----------------|---------|-------------------------------------------------------------------------------------------|
| `IteratorType`  | string  | type of iterator (`snapshot`, `cdc`)                                                      |
| `CreatedAt`     | int64   | Unix time from which the system should receive events of the resource in the CDC iterator |
| `Cursor`        | string  | resource or event identifier for receiving shifted data in the following requests         |
| `Index`         | int     | current index of the returning record from the bunch of previously received resources     |
Example:
```json
{
	"iterator_type":"cdc",
	"created_at":1652279623,
	"cursor":"evt_1KtXkmJit567F2YtZzGSIrsh",
	"index": 1
}
```

### Notes
- Data from Stripe is sorted by date of creation in descending order, with no manual sort option.
- If the user changes the resource name when the connector was created and was already in use, the connector will start working with the new resource from the beginning.