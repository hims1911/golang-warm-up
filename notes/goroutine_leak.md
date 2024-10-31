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

Packaging cake example *
Channel buffering may also affect program performance. Imagine three cooks in a cake shop, one baking, one icing, and one inscribing each cake before passing it on to the next cook in the assembly line:

In a kitchen with little space, each cook that has finished a cake must wait for the next cook to become ready to accept it. This is analogous to communication over an unbuffered channel.
If there is space for one cake between each cook, a cook may place a finished cake there and immediately start work on the next. This is analogous to a buffered channel with capacity 1.
As long as the cooks work at about the same rate on average, most of these handovers proceed quickly, smoothing out transient differences in their respective rates.
More space between cooks (larger buffers) can smooth out bigger transient variations in their rates without stalling the assembly line, such as happens when one cook takes a short break, then later rushes to catch up.
On the other hand, a buffer provides no benefit in either of the following case:

When an earlier stage of the assembly line is consistently faster than the following stage, the buffer between them will spend most of its time full.
When a later stage is faster, the buffer will usually be empty.
If the second stage is more elaborate, a single cook may not be able to keep up with the supply from the first cook or meet the demand from the third. To solve the problem, we could hire another cook to help the second, performing the same task but working independently. This is analogous to creating another goroutine communicating over the same channels.

The gopl.io/ch8/cake package simulates this cake shop, with several parameters.
