package token

import (
	"bytes"
	"text/template"
)

// pkiIssuePolicyContext is the template context provided to the rendering of
// the pkiIssuePolicyTemplate.
type pkiIssuePolicyContext struct {
	ClusterID string
}

// pkiIssueOrgPolicyContext is the template context provided to the rendering of
// the pkiIssueOrgPolicyTemplate.
type pkiIssueOrgPolicyContext struct {
	ClusterID             string
	OrganizationsRoleHash string
}

// pkiIssuePolicyTemplate provides a template of Vault policies used to
// restrict access to only being able to issue signed certificates specific to
// a Vault PKI backend of a cluster ID.
var pkiIssuePolicyTemplate = `
	path "pki-{{.ClusterID}}/issue/role-{{.ClusterID}}" {
		capabilities = ["create", "update", "delete"]
	}
	path "pki-{{.ClusterID}}/roles/" {
		capabilities = ["list"]
	}
`

// pkiIssueOrgPolicyTemplate provides a template of Vault policy used to
// restrict access to only being able to issue signed certificates specific to
// a Vault PKI backend of a organization.
var pkiIssueOrgPolicyTemplate = `
	path "pki-{{.ClusterID}}/issue/role-org-{{.OrganizationsRoleHash}}" {
		capabilities = ["create", "update", "delete"]
	}
	path "pki-{{.ClusterID}}/roles/role-org-{{.OrganizationsRoleHash}}" {
		capabilities = ["create", "update", "delete"]
	}
`

func execTemplate(t string, v interface{}) (string, error) {
	var result bytes.Buffer

	tmpl, err := template.New("policy-template").Parse(t)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&result, v)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
