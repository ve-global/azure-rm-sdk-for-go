package management

import(
	"testing"
)

func TestClient(t *testing.T) {
	clientConfig := clientConfig{

	}
	client, err := makeClient("subscription","applicationId","applicationSecret","resourceGroup", clientConfig)

	if err != nil {
		t.Error(err)
	}
	if(client.azureConfig.SubscriptionID != "subscription"){
		t.Fail()
	}
}