# common parts

Common interfaces and data structures, mainly for testing.
Implementing them is usually trivial and needs tiny pain to do.
But doing so hundreds, thousands times is not.

This package mitigate those pain.

## time-related interface

Timers are really hard to handle in tests. see [kubernetes internally mocks timer interface](https://github.com/kubernetes/kubernetes/blob/0d6dc14051612afe6a9b45708edef60f30a8cd4a/pkg/util/async/bounded_frequency_runner.go#L66-L87). The interfaces are defined to do similar things.

### GetNower

As Go convention is to name interfaces `<Verb>+er` if it only have a single method, GetNower is named after name of its only method, GetNow.

It only gets the current time, namely returns time.Now(). `GetNowImpl` does exactly that. You can implement any your desired mocked behavior in your own structure.

### Timer

The Timer is an interface almost identical to time.Timer struct. `TimerReal`, which implements Timer, is slightly improved version of time.Timer. `TimerFake` is fake implementation of Timer that is observable, synchronizable, manipulatable implementation, mainly for testing.
