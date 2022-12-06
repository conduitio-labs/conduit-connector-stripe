# Conduit Connector Stripe

### General
The Stripe connector is one of [Conduit](https://github.com/ConduitIO/conduit) plugins. It provides a source Stripe connector.

### Prerequisites
- [Go](https://go.dev/) 1.18
- (optional) [golangci-lint](https://github.com/golangci/golangci-lint) 1.50.1
- (optional) [mock](https://github.com/golang/mock) 1.6.0

### Configuration
The config passed to `Configure` can contain the following fields:

| name           | description                                                                                                          | required | example                    |
|----------------|----------------------------------------------------------------------------------------------------------------------|----------|----------------------------|
| `secretKey`    | Stripe [secret key](https://dashboard.stripe.com/apikeys).                                                           | yes      | sk_51Kr0QrJit566F2YtZAwMlh |
| `resourceName` | The name of Stripe resource. A list of supported resources can be found [here](models/resources/README.md).          | yes      | plan                       |
| `snapshot`     | The field determines whether the connector will take a snapshot of the entire resource before starting cdc mode.     | no       | false                      |
| `batchSize`    | A batch size is the number of objects to be returned. Batch size can range between 1 and 100, and the default is 10. | no       | 20                         |

### How to build it
Run `make build`.

### Testing
Run `make test`.

### Stripe Source
The `Configure` method parses the configuration and validates them.

The `Open` method parses the current position, initializes an [http client](#http-client), and initializes Snapshot (only if in the position IteratorType equals Snapshot) and CDC iterators.

The `Read` method calls the method `Next` of the current iterator and returns the next record.

The `Teardown` method calls the method `Close` of the [http client](#http-client), which calls `CloseIdleConnections` method of the [net/http](https://pkg.go.dev/net/http) package.

#### Snapshot

`Snapshot` iterator makes a copy of the data of the selected resource, sorted by date of creation in descending order.

`Snapshot` iterator algorithm:
1. iterator makes a request to the list of resource objects in Stripe without a "shift" parameter;
2. the system stores the result of the request in memory, which is a slice of objects;
3. the `Read` method creates a record from each element of the slice and updates the `Cursor` position with the `id` of the current slice element;
4. if all elements of the slice have been returned, the iterator makes the next request with the `starting_after` parameter whose value is `Cursor`;
5. if the answer is empty, the system proceeds to the `CDC` iterator, if not, it repeats from step 2.

#### CDC

The `CDC` iterator runs after Snapshot, takes data from events, and, based on those events, adds, updates, and deletes data.

`CDC` iterator algorithm:
1. the first request reads all events over the time of the connector, reverses them, and stores them in a slice;
2. the `Read` method creates a record from each event in the slice, and updates `Index` with the index of the slice element itself; 
3. if the last element of the slice is returned, it updates the `Cursor` with the index of the latest event and clears the `Index` value;
4. if all slice elements have been returned, the iterator makes the next request with the `ending_before` parameter, whose value is the `Cursor`, reverses the results, and stores them in the slice;
5. then it repeats from step 2.

**Note:** All queries in Stripe contain a `limit` parameter, the value of which is `batchSize` from the configuration, which specifies the number of returned objects.

### Position
Position is a JSON object with the following fields:

| name            | type     | description                                                                                                                                                         | 
|-----------------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `IteratorType`  | `string` | type of iterator (`snapshot`, `cdc`)                                                                                                                                |
| `CreatedAt`     | `int64`  | unix time from which the system should receive events of the resource in the CDC iterator (the parameter is set with the present time when the Position is created) |
| `Cursor`        | `string` | resource or event identifier for receiving shifted data in the following requests                                                                                   |
| `Index`         | `int`    | current index of the returning record from the batch of previously received resources                                                                               |
Example:
```json
{
	"iterator_type":"cdc",
	"created_at":1652279623,
	"cursor":"evt_1KtXkmJit567F2YtZzGSIrsh",
	"index": 1
}
```

### HTTP Client
To receive data from Stripe the connector uses [retryable HTTP client](https://github.com/hashicorp/go-retryablehttp).
Data are taken in batches (batch size parameter from [configuration](#configuration)).
In the case of an unsuccessful request, the client makes a new one.
Maximum number of retries is 4.

### Stripe
Stripe allows up to 100 read operations per second in live mode, and 25 operations per second in test mode.

Data from Stripe is sorted by date of creation in descending order, with no manual sort option.

Stripe stores [events](https://api.stripe.com/v1/events) for the last 30 days.