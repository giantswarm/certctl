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

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/certctl/integration/env"
	certsigner "github.com/giantswarm/certctl/service/cert-signer"
	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certctl/service/token"
	vaultfactory "github.com/giantswarm/certctl/service/vault-factory"
)

func TestIssuance(t *testing.T) {
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

	c.Logger.Log("level", "debug", "message", fmt.Sprintf("setup Vault PKI successful, returned token %q", token))

	client.SetToken(token)

	err = issueCerts(client)
	if err != nil {
		t.Fatalf("could not issue signed certificates, %#v", err)
	}
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

func issueCerts(client *vaultclient.Client) error {
	certSigner, err := getCertSigner(client)
	if err != nil {
		return microerror.Mask(err)
	}

	// where do organizations come from?

	newIssueConfig := spec.IssueConfig{
		ClusterID:  defaultClusterID,
		CommonName: defaultCertCommonName,
		TTL:        defaultCertTTL,
		RoleTTL:    defaultCertTokenTTL,
	}
	_, err = certSigner.Issue(newIssueConfig)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func getVaultClient(vaultAddr string) (*vaultclient.Client, error) {
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.Address = vaultAddr
	newVaultFactoryConfig.AdminToken = env.VaultToken()
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

func getCertSigner(client *vaultclient.Client) (spec.CertSigner, error) {
	newCertSignerConfig := certsigner.DefaultConfig()
	newCertSignerConfig.VaultClient = client
	newCertSigner, err := certsigner.New(newCertSignerConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return newCertSigner, nil
}

func createPKIBackend(svc pki.Service) error {
	createConfig := pki.CreateConfig{
		AllowedDomains: defaultCommonName,
		ClusterID:      defaultClusterID,
		CommonName:     defaultCommonName,
		TTL:            defaultCATTL,
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
	vaultSvc, err := c.Clients.K8sClient().CoreV1().Services("default").Get("vault", meta_v1.GetOptions{})
	if err != nil {
		return "", microerror.Mask(err)
	}

	// We will access Vault service from the test container using the k8s API
	// server address and the service NodePort.
	hostURL, err := url.Parse(c.Clients.RESTConfig().Host)
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
	b := backoff.NewExponential(backoff.MediumMaxWait, backoff.LongMaxInterval)
	n := func(err error, delay time.Duration) {
		log.Printf("failed connection to vault %#v", err)
	}

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
