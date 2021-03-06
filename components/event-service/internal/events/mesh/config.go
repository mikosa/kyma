package mesh

import (
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go"
)

type Configuration struct {
	Source           string
	CloudEventClient cloudevents.Client
}

// InitConfig returns an Event mesh configuration instance which is required by the Event publish and delivery flows.
func InitConfig(source string, eventMeshURL string) (*Configuration, error) {
	ceClient, err := newCloudEventClient(eventMeshURL)
	if err != nil {
		return nil, err
	}
	config := &Configuration{
		Source:           source,
		CloudEventClient: ceClient,
	}
	return config, nil
}

// newCloudEventClient returns a new cloudevent client which points to the HTTP adapter created via the "Event Source",
// this is the internal entrypoint to our Event mesh.
func newCloudEventClient(eventMeshURL string) (cloudevents.Client, error) {
	if _, err := url.Parse(eventMeshURL); err != nil {
		return nil, err
	}

	transport, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(eventMeshURL),
		cloudevents.WithBinaryEncoding(),
	)
	if err != nil {
		return nil, err
	}
	client, err := cloudevents.NewClient(transport, cloudevents.WithTimeNow())
	if err != nil {
		return nil, err
	}
	return client, nil
}
