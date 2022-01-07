package rancher

import (
	"context"
	"github.com/ebauman/prometheus-rancher-exporter/semver"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	settingGVR                 = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "settings"}
	settingGVRNumberOfClusters = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "clusters"}
)

type Client struct {
	Client dynamic.Interface
	Config *rest.Config
}

func (r Client) GetRancherVersion() (map[string]int64, error) {

	res, err := r.Client.Resource(settingGVR).Get(context.Background(), "server-version", v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	version, _, _ := unstructured.NestedString(res.UnstructuredContent(), "value")

	result, err := semver.Parse(TrimVersionChar(version))

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Client) GetNumberOfManagedClusters() (int, error) {

	res, err := r.Client.Resource(settingGVRNumberOfClusters).Get(context.Background(), "", v1.GetOptions{})
	if err != nil {
		return 0, err
	}

	number, _, _ := unstructured.NestedStringMap(res.UnstructuredContent())

	log.Info("Number of Clusters", len(number))

	return 1, nil
}

// Version returned from CRD is in the format of "v.N.N.N", trim the leading "v"
func TrimVersionChar(version string) string {
	for i := range version {
		if i > 0 {
			return version[i:]
		}
	}
	return version[:0]
}
