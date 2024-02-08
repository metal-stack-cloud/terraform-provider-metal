package kubeconfig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestKubeconfigExtraction(t *testing.T) {
	kubeconfig := `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRYWtuNjNlWThreitkZrcHdaazVIR3RPaWErRXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0VCQUlmblROZzNubU4yakFKWUkvYnVqWDhwMk11TWpZOStYVThwSjNIbi9xVjZId05MbXdXSwpIYk5NMDJ6WmxXeFBJZ2dDV3BFQW84VVFKU044RVBVQk1mNDZJYnpVbG9LN3JXMGFwK3crb0prUWZ6cTNZb3JtCjZEaDFyQzZMamg1Q3ZBZ0VtYzMyOUhlVFNDaUtVeTZ3V01kYy9oREQwNEtsQlhHaXpsY2RqYVN6ME1Ud2R2VSsKMURkWk1Lbmg1OVVCVi9TR2VqRGJsdWNhaU84RkZGbFp0SHYrQ1RlYVoxYzdqWFo2cTdVTFI0L2k3SzM0WE5qOQpxOXhsR3N1N045dWdoWGFoYUF4NlNVZ3B3bUZybmVaWkxZcDdLTDFybzZmKzNqZGo1K3JWV0dueVE2ZkR3ME1iCjRXcHFjYTcvcGFnT0Q4cVM4S1ZiWUJPek9ka1Z5N0ozMXFzPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://api.cluster.f8e67080ba.k8s.metalstackcloud.io
  name: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1akNDQWM2Z0F3SUJBZ0lRYWtuNjNlWThreitaRE1MZ09kd2sY2RqYVN6ME1Ud2R2VSsKMURkWk1Lbmg1OVVCVi9TR2VqRGJsdWNhaU84RkZGbFp0SHYrQ1RlYVoxYzdqWFo2cTdVTFI0L2k3SzM0WE5qOQpxOXhsR3N1N045dWdoWGFoYUF4NlNVZ3B3bUZybmVaWkxZcDdLTDFybzZmKzNqZGo1K3JWV0dueVE2ZkR3ME1iCjRXcHFjYTcvcGFnT0Q4cVM4S1ZiWUJPek9ka1Z5N0ozMXFzPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://api.cluster.f8e67080ba.internal.k8s.metalstackcloud.io
  name: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-internal
contexts:
- context:
    cluster: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
    user: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
  name: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
- context:
    cluster: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-internal
    user: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
  name: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-internal
current-context: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
kind: Config
preferences: {}
users:
- name: garden-f8e67080-ba68-41d2-ad44-59dc65a09a33--cluster-external
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURKekNDQWcrZ0F3SUJBZ0lFkWUl5YjFjbk16Tmx6OXBLa3lmM1V2OXNoTDAwQlJhVTRCVXNVVFdNekRKUTQ5dnRpYVg0cldlbCt2UwpLV3dqbGdLZEJiL25zczJlNmhrcG9qL3BvNXp0bmY4ditJSS9iY1MxM2Q5N3VOY2E3MFVyY2ZSTHFOQVY3VVc1Ckc0b2plUzk5YlkvRG1ZSm1LZldVWFN3am0rVXZ4bnBpRXpTY3FLeWNzRFBOVEFkdmZaNnExSGJiTCtkWEhOelQKMWpxZ2FHUGxIdGRWck5GajF3UmFQV0FFS3Ntai9oLzd6QmtlK0tzWUcvbTRnOXozOHFLVHd3Y0x2QT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBeml1SitIN2lZRFJieUZicWhnK2twMnpPQW11OVlLZ29HRzdTbkFDClRzUTNTYlFRNmlDU3gxbHgvT05WN08wQ2dZQXVMTkMwZmtnUjBaL0dRS3FiVDE5elZaWHJpNUhPMmZCNkJWankKMmtLSHNjZUFMQ05zNzhnK1I1dkt1TVFFMGRuQ05rTzlUdW1IMDVyeWhFTm5VYW9IR0ZFSUY5eWYxUWFQSFg2Sgo0YStwUnNwbWpVRlhiWXY3ZEdYUGRKSFpVV2xsb3ppckZaQ1BZS0pudm1sTXN6WEZmblcwWTJ3WjRGKy8yVFV1CmVnVVhXUUtCZ1FDYXBhUEk2WlEraEdsd21FaGU1L3FVNnZaNUVxdUFHN2tLU2gyZ1FjQkpzb0FvYTdHQXpZQlIKMVlhSmtrdytHRTNkZVU2Z0tFZTBXWVhGd2Z5VVhUQlVGUzI1WUp5NnYzOGhsSk0yb3JvR2hzaWpjZi9OUEVocApnYmQrbEg0bHRGMUtFOVpYMFQ2Q0pFTXpkaktNV3llMGU4TWMveldIdVNNblliTktJR0J1TWc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`
	diags := diag.Diagnostics{}
	external := parseKubeconfig(kubeconfig, diags)
	if len(diags) != 0 {
		t.Errorf("expected no diagnostics, got %s", diags)
	}
	if external.Host.IsNull() || external.Host.IsUnknown() {
		t.Errorf("host could not be parsed")
	}
	if external.ClientKey.IsNull() || external.ClientKey.IsUnknown() {
		t.Errorf("client key could not be parsed")
	}
	if external.ClientCertificate.IsNull() || external.ClientCertificate.IsUnknown() {
		t.Errorf("client certificate could not be parsed")
	}
	if external.ClusterCaCertificate.IsNull() || external.ClusterCaCertificate.IsUnknown() {
		t.Errorf("cluster ca certificate could not be parsed")
	}
}
