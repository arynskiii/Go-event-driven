# Publishing Messages

### Picking the library

When you serve a website, you don't handle the HTTP protocol manually.
You use a library that does the heavy lifting for you and gives you a nice API to define all the endpoints.

Most popular Pub/Subs have their own SDK in Go.
But they are usually low-level libraries, and working with each Pub/Sub is quite different.

Back in 2018, we thought working with messages should be as easy as working with HTTP requests.
There was no library in Go that allowed it, so we decided to create it.
It's called [Watermill](https://github.com/ThreeDotsLabs/watermill), and we've been using it in production for many different projects since then.

**We're going to use Watermill for the rest of the training.
We believe it will help you focus on the high-level concepts instead of the low-level details.**
We want to teach you the ideas behind event-driven architecture, not the details of a specific Pub/Sub.

We designed Watermill as a lightweight library, not a framework, so there's no vendor lock-in.
If you prefer to use anything else, it should be straightforward to translate the examples.

Remember, **this training is about timeless event-driven patterns**, not Watermill-specific approach.
You can apply the same ideas in any other programming language or library.

### The Publisher

Watermill hides all the complexity of Pub/Subs behind just two interfaces: the `Publisher` and the `Subscriber`.

For now, let's consider the first one.

```go
type Publisher interface {
	Publish(topic string, messages ...*Message) error
	Close() error
}
```

To publish a message, you need to pass a `topic` and a slice of messages.

To *publish* a message means to append it to the given topic.
Anyone who subscribes to the same topic will receive the messages on it in a first-in, first-out (FIFO) fashion.

{{tip}}

FIFO is a common way to deliver messages, but it can vary depending on the Pub/Sub and how it's configured.

This is true for many behaviors we describe in this training.
Always check the documentation of the Pub/Sub you use to confirm it works as you expect.

{{endtip}}

How you create the publisher instance depends on the Pub/Sub you choose.
Each library provides its own constructors.

Here's how to create one for Redis Streams:

```go
logger := watermill.NewStdLogger(false, false)

rdb := redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDR"),
})

publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
	Client: rdb,
}, logger)
```

{{tip}}

The `NewStdLogger`'s arguments are for `debug` and `trace`, respectively.
You don't need to use this logger: You can adapt any other logger you use to the [`watermill.LoggerAdapter`](https://github.com/ThreeDotsLabs/watermill/blob/559222086a70e83f930fd904c2a53991749f3877/log.go#L43) interface.

{{endtip}}

You need the following imports to make the code above work:

```go
import (
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)
```

To create a message, use the `NewMessage` constructor.
It takes just two arguments.

```go
msg := message.NewMessage(watermill.NewUUID(), []byte(orderID))
```

The first argument is the message's UUID, which is used mainly for debugging.
Most of the time, any kind of UUID is fine.
Watermill provides helper functions to generate them.

The second argument is the payload.
It's a slice of bytes, so it can be anything you want, as long as you can marshal it.
You can send a string, a JSON object, or even a binary file.

To publish the message, call the `Publish` method on the publisher:

```go
err := publisher.Publish("orders", msg)
```

Remember to handle the error, as publishing works over the network and can fail for many reasons.

{{.Exercise}}

Create a Redis Streams publisher, and **publish two messages on the `progress` topic.**
The first one's payload should be `50`, and the second one's should be `100`.

To get the necessary dependencies, either run `go get` for them individually or run `go mod tidy` after adding the imports.
