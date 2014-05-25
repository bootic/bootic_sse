## Go Client for Bootic's HTTP events stream

### Example

```go
import(
  "github.com/bootic/bootic_sse"
  data "github.com/bootic/bootic_go_data"
)

func main() {
	client, _ := bootic_sse.NewClient("https://some.stream.com", "token")

	events := make(data.EventsChannel)
	client.Subscribe(events)
	
	for {
		log.Println(<-events)
	}
	
}
```

### ToDo

Not production ready yet. Useful to subscribe to remote HTTP stream locally.

* Reconnect on failure