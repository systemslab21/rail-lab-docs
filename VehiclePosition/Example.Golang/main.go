package main

import (
	"bytes"
	"context"
	"crypto/tls"
	fmt "fmt"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/examples/dtls/pki"
	"github.com/plgd-dev/go-coap/v3/message"
)

func main() {
	config, err := createClientConfig(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}

	co, err := dtls.Dial(fmt.Sprintf("%s:5684", os.Getenv("COAP_HOST")), config)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	va := &VehicleAttributes{
		NidEngine: "NID_ENGINE",
	}
	va_out, err := proto.Marshal(va)

	va_ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	va_req, err := co.NewPostRequest(va_ctx, "/api/v1/attributes", message.AppOctets, bytes.NewReader(va_out))

	_, err = co.Do(va_req)

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// Send vehicle position

	vp := &VehiclePosition{
		Timestamp: 1234,
		Position:  500,
		Speed:     2.4,
		ResetId:   "foo",
	}
	vp_out, err := proto.Marshal(vp)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	vp_req, err := co.NewPostRequest(ctx, "/api/v1/telemetry", message.AppOctets, bytes.NewReader(vp_out))
	vp_req.SetType(message.NonConfirmable)

	_, err = co.Do(vp_req)

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
}

func createClientConfig(ctx context.Context) (*piondtls.Config, error) {
	keyBytes, err := os.ReadFile("cert/keyplain.pem")
	if err != nil {
		return nil, err
	}

	certBytes, err := os.ReadFile("cert/cert.pem")
	if err != nil {
		return nil, err
	}

	certificate, err := pki.LoadKeyAndCertificate(keyBytes, certBytes)
	if err != nil {
		return nil, err
	}

	return &piondtls.Config{
		Certificates:         []tls.Certificate{*certificate},
		ExtendedMasterSecret: piondtls.RequireExtendedMasterSecret,
		InsecureSkipVerify:   true,
	}, nil
}
