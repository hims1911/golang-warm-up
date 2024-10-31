The example below shows an application of a buffered channel. It makes parallel requests to three mirror sites. It sends their responses over a buffered channel, then receives and returns only the first response, which is the quickest one to arrive. Thus mirroredQuery returns a result even before the two slower servers have responded. (It's quite normal for several goroutines to send values to the same channel concurrently, as in this example, or to receive from the same channel.)

```
func mirroredQuery() string {
    responses := make(chan string, 3)
    go func() { responses <- request("asia.gopl.io") }()
    go func() { responses <- request("europe.gopl.io") }()
    go func() { responses <- request("americas.gopl.io") }()
    return <-responses // return the quickest response
}
```

func request(hostname string) (response string) { /* ... */ }
If we had used an unbuffered channel in the above example, the two slower goroutines would have gotten stuck trying to send their responses on a channel from which no goroutine will ever receive. This situation, called a goroutine leak, would be a bug. Unlike garbage variables, leaked goroutines are not automatically collected, so it is important to make sure that goroutines terminate themselves when no longer needed.

The choice between unbuffered and buffered channels, and the choice of a buffered channel's capacity, may both affect the correctness of a program.

Unbuffered channels give stronger synchronization guarantees because every send operation is synchronized with its corresponding receive
With buffered channels, the send and receive operations are decoupled.
When we know an upper bound on the number of values that will be sent on a channel, it's not unusual to create a buffered channel of that size and perform all the sends before the first value is received. Failure to allocate sufficient buffer capacity would cause the program to deadlock.
