package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/pflag"
)

type Config struct {
	// HMACSecretKey is the secret key to authenticate HMAC endpoint
	// connection JWTs.
	HMACSecretKey string `json:"hmac_secret_key" yaml:"hmac_secret_key"`

	// RSAPublicKey is the public key to authenticate RSA endpoint
	// connection JWTs.
	RSAPublicKey string `json:"rsa_public_key" yaml:"rsa_public_key"`

	// ECDSAPublicKey is the public key to authenticate ECDSA endpoint
	// connection JWTs.
	ECDSAPublicKey string `json:"ecdsa_public_key" yaml:"ecdsa_public_key"`

	// Audience is the required 'aud' claim of the authenticated JWTs.
	//
	// If not given the 'aud' claim will be ignored.
	Audience string `json:"audience" yaml:"audience"`

	// Issuer is the required 'iss' claim of the authenticated JWTs.
	//
	// If not given the 'iss' claim will be ignored.
	Issuer string `json:"issuer" yaml:"issuer"`
}

func (c *Config) AuthEnabled() bool {
	return c.HMACSecretKey != "" || c.RSAPublicKey != "" || c.ECDSAPublicKey != ""
}

func (c *Config) Load() (JWTVerifierConfig, error) {
	verifierConf := JWTVerifierConfig{
		HMACSecretKey: []byte(c.HMACSecretKey),
		Audience:      c.Audience,
		Issuer:        c.Issuer,
	}

	if c.RSAPublicKey != "" {
		rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(
			[]byte(c.RSAPublicKey),
		)
		if err != nil {
			return JWTVerifierConfig{}, fmt.Errorf("parse rsa public key: %w", err)
		}
		verifierConf.RSAPublicKey = rsaPublicKey
	}
	if c.ECDSAPublicKey != "" {
		ecdsaPublicKey, err := jwt.ParseECPublicKeyFromPEM(
			[]byte(c.ECDSAPublicKey),
		)
		if err != nil {
			return JWTVerifierConfig{}, fmt.Errorf("parse ecdsa public key: %w", err)
		}
		verifierConf.ECDSAPublicKey = ecdsaPublicKey
	}

	return verifierConf, nil
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	prefix += ".auth."

	fs.StringVar(
		&c.HMACSecretKey,
		prefix+"hmac-secret-key",
		c.HMACSecretKey,
		`
Secret key to authenticate HMAC endpoint connection JWTs.`,
	)
	fs.StringVar(
		&c.RSAPublicKey,
		prefix+"rsa-public-key",
		c.RSAPublicKey,
		`
Public key to authenticate RSA endpoint connection JWTs.`,
	)
	fs.StringVar(
		&c.ECDSAPublicKey,
		prefix+"ecdsa-public-key",
		c.ECDSAPublicKey,
		`
Public key to authenticate ECDSA endpoint connection JWTs.`,
	)
	fs.StringVar(
		&c.Audience,
		prefix+"audience",
		c.Audience,
		`
Audience of endpoint connection JWT  to verify.

If given the JWT 'aud' claim must match the given audience. Otherwise it
is ignored.`,
	)
	fs.StringVar(
		&c.Issuer,
		prefix+"issuer",
		c.Issuer,
		`
Issuer of endpoint connection JWT  to verify.

If given the JWT 'iss' claim must match the given issuer. Otherwise it
is ignored.`,
	)
}
