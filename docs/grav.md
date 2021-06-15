# Reactr ➕ Grav

Reactr is designed to integrate with the other [Suborbital](https://suborbital.dev) projects such as [Grav](https://github.com/suborbital/grav). Grav is a decentralized message bus which allows for your application code to communicate in a scalable, resilient way.

## Handle Messages
Reactr can respond to messages by connecting to a `grav.Pod` using `HandleMsg`:
```golang
reactr := rt.New()
g := grav.New()

reactr.HandleMsg(g.Connect(), msgTypeLogin, &loginEmailRunner{})
```
Whenever a message with the given type is received from the bus, a `Job` will be queued to be handled by the provided Runnable. The `Job` will contain the message data.

The result returned by the Runnable's `Run` function may be a `grav.Message`. If so, it will be sent back out over the message bus. Anything else will be put into a mesage (by converting it into bytes) and sent back over the bus. If `Run` returns an error, a message with type `reactr.runerr` will be sent. If `Run` returns `nil, nil`, then a message of type `reactr.nil` will be sent. All messages sent will be a reply to the message that triggered the job.

Further integrations with `Grav` are in the works, along with improvements to Reactr's [FaaS](./faas.md) capabilities, which is powered by Suborbital's [Vektor](https://github.com/suborbital/vektor) framework. 