package mqttclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/eclipse/paho.mqtt.golang"
)

// ClientConfig contains the configuration for the remote mqtt client
type ClientConfig struct {
	CertsFolder string
	CloudMQTT   string
	ClientID    string
}

// NewClient returns a new mqtt client connected to AWS.
func NewClient(config ClientConfig) (mqtt.Client, chan bool, error) {
	opts := mqtt.NewClientOptions().AddBroker(config.CloudMQTT)
	opts.SetClientID(config.ClientID)

	lastWillMessage, err := lastWillMessage(config.ClientID)
	if err != nil {
		return nil, nil, err
	}

	opts.SetBinaryWill(lastWillTopic(config.ClientID), lastWillMessage, 1, false)

	if config.CertsFolder != "" {
		tlsConfig, err := newTLSConfig(config.CertsFolder)
		if err != nil {
			return nil, nil, fmt.Errorf("can't create a tls config: %s", err)
		}
		opts.SetTLSConfig(tlsConfig)
	}

	// Since aws IoT broker doesn't support persistency
	// we have to resubscribe on our own.
	onConnect := make(chan bool)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		onConnect <- true
	})

	// Create and start a client using the above ClientOptions
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, nil, fmt.Errorf("can't connect: %s", token.Error())
	}

	return c, onConnect, nil
}

func newTLSConfig(folder string) (*tls.Config, error) {
	certFile := path.Join(folder, "certificate.pem.crt")
	keyFile := path.Join(folder, "private.pem.key")
	caFile := path.Join(folder, "root-CA.crt")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("can't load client cert: %s", err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("can't load ca cert: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup tls
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return tlsConfig, nil
}

func lastWillTopic(clientID string) string {
	return fmt.Sprintf("last_will/%s", clientID)
}

func lastWillMessage(clientID string) ([]byte, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"clientID": clientID,
	})
	if err != nil {
		return nil, fmt.Errorf("Error while marshaling last will: %s", err)
	}

	return payload, nil
}
