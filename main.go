// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"drone-env-signed/plugin"
	"github.com/drone/drone-go/plugin/environ"

	"github.com/gbrlsnchs/jwt/v3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type spec struct {
	Bind       string `envconfig:"DRONE_BIND"`
	Debug      bool   `envconfig:"DRONE_DEBUG"`
	Secret     string `envconfig:"DRONE_SECRET"`
	PrivateKey string `envconfig:"DRONE_PRIVATE_KEY"`
}

func main() {
	var alg jwt.Algorithm

	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":80"
	}

	block, _ := pem.Decode([]byte(spec.PrivateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		logrus.Fatalln("invalid private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logrus.Fatal(err)
	}

	switch k := privKey.(type) {
	case ed25519.PrivateKey:
		alg = jwt.NewEd25519(jwt.Ed25519PrivateKey(k))
	case *rsa.PrivateKey:
		alg = jwt.NewRS256(jwt.RSAPrivateKey(k))
	case *ecdsa.PrivateKey:
		alg = jwt.NewES256(jwt.ECDSAPrivateKey(k))
	default:
		logrus.Fatalln("unsupported private key type")
	}

	handler := environ.Handler(
		spec.Secret,
		plugin.New(
			spec.Secret,
			alg,
		),
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
