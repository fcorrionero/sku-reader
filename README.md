# Sku Reader
<hr/>

### Prerequisites

You need docker-compose and go installed in your machine

### Run the project

```
make run
```

### Connect client to host

After run the project in another terminal

```
nc localhost 4000
```

Then you can start sending messages to the host.

### Tear down the project

```
make teardown-env
```

### Run Tests

```
make integration-test
make unit-test
```

