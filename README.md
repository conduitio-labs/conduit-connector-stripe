# Conduit Connector Stripe

### General

The Stripe connector is one of [Conduit](https://github.com/ConduitIO/conduit) builtin plugins. It provides a source 
Stripe connector.

### How to build it

Run `make`.

### Testing

Run `make test` to run all the tests. You must set the environment variables (`STRIPE_SECRET_KEY`, `STRIPE_SOURCE_NAME`)
before you run all the tests. If not set, the tests that use these variables will be ignored.