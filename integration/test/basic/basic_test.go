// +build k8srequired

package basic

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
	"github.com/giantswarm/certctl/service/vault-factory"
)

func TestSetup(t *testing.T) {
	vaultAddr, err := getVaultAddr()
	if err != nil {
		t.Fatalf("could not create Vault address, %#v", err)
	}

	client, err := getVaultClient(vaultAddr)
	if err != nil {
		t.Fatalf("could not create Vault client, %#v", err)
	}

	err = waitForVault(client)
	if err != nil {
		t.Fatalf("timeout waiting for Vault, %#v", err)
	}

	token, err := setUp(client)
	if err != nil {
		t.Fatalf("could not setup Vault PKI, %#v", err)
	}

	l.Log("level", "debug", "message", fmt.Sprintf("setup Vault PKI successful, returned token %q", token))
}

func setUp(client *vaultclient.Client) (string, error) {
	pkiService, err := getPKIService(client)
	if err != nil {
		return "", microerror.Mask(err)
	}

	tokenService, err := getTokenService(client)
	if err != nil {
		return "", microerror.Mask(err)
	}

	err = createPKIBackend(pkiService)
	if err != nil {
		return "", microerror.Mask(err)
	}

	token, err := createToken(tokenService)
	if err != nil {
		return "", microerror.Mask(err)
	}
	return token, nil
}

func getVaultClient(vaultAddr string) (*vaultclient.Client, error) {
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.Address = vaultAddr
	newVaultFactoryConfig.AdminToken = vaultToken
	newVaultFactoryConfig.TLS = &vaultclient.TLSConfig{}
	newVaultFactory, err := vaultfactory.New(newVaultFactoryConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	vc, err := newVaultFactory.NewClient()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return vc, nil
}

func getPKIService(client *vaultclient.Client) (pki.Service, error) {
	pkiConfig := pki.DefaultServiceConfig()
	pkiConfig.VaultClient = client
	pkiService, err := pki.NewService(pkiConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return pkiService, nil
}

func getTokenService(client *vaultclient.Client) (token.Service, error) {
	tokenConfig := token.DefaultServiceConfig()
	tokenConfig.VaultClient = client
	tokenService, err := token.NewService(tokenConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return tokenService, nil
}

func createPKIBackend(svc pki.Service) error {
	createConfig := pki.CreateConfig{
		ClusterID: defaultClusterID,
		TTL:       defaultCATTL,
	}
	err := svc.Create(createConfig)
	if err != nil {
		microerror.Mask(err)
	}
	return nil
}

func createToken(svc token.Service) (string, error) {
	createConfig := token.CreateConfig{
		ClusterID: defaultClusterID,
		Num:       1,
		TTL:       defaultTokenTTL,
	}
	tokens, err := svc.Create(createConfig)
	if err != nil {
		return "", microerror.Mask(err)
	}
	return tokens[0], nil
}

func getVaultAddr() (string, error) {

	vaultSvc, err := f.K8sClient().
		CoreV1().
		Services("default").
		Get("vault", meta_v1.GetOptions{})
	if err != nil {
		return "", microerror.Mask(err)
	}

	// we will access Vault service from the test container using the k8s API
	// server address and the service NodePort.
	restCfg := f.RestConfig()
	hostURL, err := url.Parse(restCfg.Host)
	if err != nil {
		return "", microerror.Mask(err)
	}
	serverAddr := strings.Split(hostURL.Host, ":")[0]
	port := vaultSvc.Spec.Ports[0].NodePort
	addr := fmt.Sprintf("http://%s:%d", serverAddr, port)

	return addr, nil
}

func waitForVault(client *vaultclient.Client) error {
	o := func() error {
		req := client.NewRequest("HEAD", "/sys/health")
		resp, err := client.RawRequest(req)
		if err != nil {
			return microerror.Mask(err)
		}
		if resp.StatusCode != http.StatusOK {
			return microerror.Mask(fmt.Errorf("unexpected status code from Vault, want %d, got %d", http.StatusOK, resp.StatusCode))
		}
		return nil
	}
	b := backoff.NewExponentialBackOff()
	n := func(err error, delay time.Duration) {
		log.Printf("failed connection to vault %#v", err)
	}

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
