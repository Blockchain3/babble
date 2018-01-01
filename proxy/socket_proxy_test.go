package proxy

import (
	"reflect"
	"testing"
	"time"

	"github.com/babbleio/babble/common"
	"github.com/babbleio/babble/crypto"
	"github.com/babbleio/babble/hashgraph"
	aproxy "github.com/babbleio/babble/proxy/app"
)

func createDummyEvent(t *testing.T) hashgraph.Event {

	privateKey, _ := crypto.GenerateECDSAKey()
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)

	body := hashgraph.EventBody{}
	body.Transactions = [][]byte{[]byte("abc"), []byte("def")}
	body.Parents = []string{"self", "other"}
	body.Creator = publicKeyBytes
	body.Timestamp = time.Now().UTC()

	event := hashgraph.Event{Body: body}
	if err := event.Sign(privateKey); err != nil {
		t.Fatalf("Error signing Event: %s", err)
	}
	res, err := event.Verify()
	if err != nil {
		t.Fatalf("Error verifying signature: %s", err)
	}
	if !res {
		t.Fatalf("Verify returned false")
	}

	return event
}

func TestSokcetProxyServer(t *testing.T) {
	clientAddr := "127.0.0.1:9990"
	proxyAddr := "127.0.0.1:9991"
	proxy := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))
	submitCh := proxy.SubmitCh()

	tx := []byte("the test transaction")

	// Listen for a request
	go func() {
		select {
		case st := <-submitCh:
			// Verify the command
			if !reflect.DeepEqual(st, tx) {
				t.Fatalf("tx mismatch: %#v %#v", tx, st)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// now client part connecting to RPC service
	// and calling methods
	dummyClient, err := NewDummySocketClient(clientAddr, proxyAddr, common.NewTestLogger(t))
	if err != nil {
		t.Fatal(err)
	}
	err = dummyClient.SubmitTx(tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSocketProxyClient(t *testing.T) {
	clientAddr := "127.0.0.1:9992"
	proxyAddr := "127.0.0.1:9993"
	proxy := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))

	dummyClient, err := NewDummySocketClient(clientAddr, proxyAddr, common.NewTestLogger(t))
	if err != nil {
		t.Fatal(err)
	}
	clientCh := dummyClient.babbleProxy.CommitCh()

	event := createDummyEvent(t)

	// Listen for a request
	go func() {
		select {
		case se := <-clientCh:
			if !reflect.DeepEqual(se, event) {
				t.Fatalf("event mismatch: %#v %#v", se, event)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	err = proxy.CommitEvent(event)
	if err != nil {
		t.Fatal(err)
	}
}
