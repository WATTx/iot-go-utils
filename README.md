# go-aws-iot-utils

## mqttclient

`mqttclient` helps to construct an aws iot mqtt client.

Example usage:

```
import "github.com/wattx/go-aws-iot-utils/mqttclient"

clientConfig := mqttclient.ClientConfig{
		ClientID:    "client-id",
		Endpoint:    "aws-mqtt-endpoint",
		CertsFolder: "./certs",
}


mqttc, onConnect, err := mqttclient.NewClient(clientConfig)
if err != nil {
    log.Fatal(err)
}
```

`Endpoint` usually has a format of `tls://FOOBAR.iot.eu-west-1.amazonaws.com:8883`

`onConnect` is used to notify when reconnect to the mqtt broker is happening. Since aws iot doesn't support persistent sessions you can use as a work around for refreshing your subscriptions.

`CertsFolder` should contain:
- `certificate.pem.crt`
- `private.pem.key`
- `root-CA.crt`

### How to get iot certificates

1. Log to: https://<your-account>.signin.aws.amazon.com/console
2. AWS-iot -> Create-a-resource -> certificate
3. Download the private key and certificate files to your machine (you don't need the private key file).
4. Get the last certificate from:
```
wget https://www.symantec.com/content/en/us/enterprise/verisign/roots/VeriSign-Class%203-Public-Primary-Certification-Authority-G5.pem
```
5. Rename the certificates:
```
mv VeriSign-Class\ 3-Public-Primary-Certification-Authority-G5.pem root-CA.crt
mv <your-num>-certificate.pem.crt certificate.pem.crt
mv <your-num>-private.pem.key private.pem.key
```
6. Move the renamed certificates to `CertsFolder` (for example `./certs`)
