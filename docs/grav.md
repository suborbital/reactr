# Hive ➕ Grav

Hive is designed to integrate with the other [Suborbital](https://suborbital.dev) projects such as [Grav](https://github.com/suborbital/grav). Grav is a decentralized message bus which allows for your application code to communicate in a scalable, resilient way.

## Handle Messages
Hive can respond to messages by connecting to a `grav.Pod` using `HandleMsg`:
```golang
hive := hive.New()
g := grav.New()

hive.HandleMsg(g.Connect(), msgTypeLogin, &loginEmailRunner{})
```
Whenever a message with the given type is received from the bus, a `Job` will be queued to be handled by the provided Runnable. The `Job` will contain the message, and `job.Msg()` makes it easy to retreive (with type conversions happening automatically).

The result returned by the Runnable's `Run` function may be a `grav.Message`. If so, it will be sent back out over the message bus. Anything else will be put into a mesage (by converting it to a byte array) and sent back over the bus. If `Run` returns an error a message with type `hive.joberr` will be sent. If `Run` returns `nil, nil`, then a message of type `hive.nil` will be sent. All messages sent will be a reply to the message that triggered the job.

Further integrations with `Grav` are in the works, along with improvements to Hive's [FaaS](./faas.md) capabilities, which is powered by Suborbital's [Vektor](https://github.com/suborbital/vektor) framework. 