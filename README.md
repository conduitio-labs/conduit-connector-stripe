# Conduit Connector Stripe

### General

The Stripe connector is one of [Conduit](https://github.com/ConduitIO/conduit) builtin plugins. It provides a source 
Stripe connector.

### Prerequisites

- [Go](https://go.dev/) 1.18
- (optional) [golangci-lint](https://github.com/golangci/golangci-lint) 1.45.2
- (optional) [mock](https://github.com/golang/mock) 1.6.0

### Configuration

The config passed to `Configure` can contain the following fields.

| name        | description                                                                                                      | required | example                      |
|-------------|------------------------------------------------------------------------------------------------------------------|----------|------------------------------|
| `key`       | Stripe [secret key](https://dashboard.stripe.com/apikeys).                                                       | yes      | "sk_51Kr0QrJit566F2YtZAwMlh" |
| `resource`  | The name of Stripe resource. A list of supported resources can be found [below](#a-list-of-supported-resources). | yes      | "plan"                       |
| `retry_max` | The maximum number of requests to Stripe in case of failure. By default is 3. The maximum is 10.                 | no       | "5"                          |                                                                                           | yes      | "id"                                            |
| `limit`     | Count of records in one butch. By default is 50. The maximum is 100.                                             | no       | "70"                         |

### How to build it

Run `make build`.

### Testing

Run `make test`.

### Stripe Source

The Stripe Source Connector calls `Configure` method to parse the configurations.
The `Open` method parses the current position, initializes both iterators (Snapshot and CDC), and then initializes an 
[HTTP client](https://github.com/hashicorp/go-retryablehttp) using `key` and `resource`.
The `Read` method returns the next record.
The `Ack` method checks if the record with the position was recorded (under development).

### Position

Position has fields: `iterator type`, `cursor`, `created_at`, `index`.

### A list of supported resources

- [subscription](https://stripe.com/docs/api/subscriptions)
- [plan](https://stripe.com/docs/api/plans)