package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"net/http"
	"net/url"

	"github.com/khulnasoft/meshplay-operator/api/v1alpha1"
	"github.com/khulnasoft/meshkit/utils"
	meshplaykube "github.com/khulnasoft/meshkit/utils/kubernetes"
)

const BrokerPingEndpoint = "8222/connz"

type Connections struct {
	Connections []connection `json:"connections"`
}

type connection struct {
	Name string `json:"name"`
}

func GetBrokerEndpoint(kclient *meshplaykube.Client, broker *v1alpha1.Broker) string {
	endpoint := broker.Status.Endpoint.Internal
	if len(strings.Split(broker.Status.Endpoint.Internal, ":")) > 1 {
		port, _ := strconv.Atoi(strings.Split(broker.Status.Endpoint.Internal, ":")[1])
		if !utils.TcpCheck(&utils.HostPort{
			Address: strings.Split(broker.Status.Endpoint.Internal, ":")[0],
			Port:    int32(port),
		}, nil) {
			endpoint = broker.Status.Endpoint.External
			port, _ = strconv.Atoi(strings.Split(broker.Status.Endpoint.External, ":")[1])
			if !utils.TcpCheck(&utils.HostPort{
				Address: strings.Split(broker.Status.Endpoint.External, ":")[0],
				Port:    int32(port),
			}, nil) {
				if !utils.TcpCheck(&utils.HostPort{
					Address: "host.docker.internal",
					Port:    int32(port),
				}, nil) {
					u, _ := url.Parse(kclient.RestConfig.Host)
					if utils.TcpCheck(&utils.HostPort{
						Address: u.Hostname(),
						Port:    int32(port),
					}, nil) {
						endpoint = fmt.Sprintf("%s:%d", u.Hostname(), int32(port))
					}
				} else {
					endpoint = fmt.Sprintf("host.docker.internal:%d", int32(port))
				}
			}
		}
	}

	return endpoint
}

func applyOperatorHelmChart(chartRepo string, client meshplaykube.Client, meshplayReleaseVersion string, delete bool, overrides map[string]interface{}) error {
	var (
		act   = meshplaykube.INSTALL
		chart = "meshplay-operator"
	)
	if delete {
		act = meshplaykube.UNINSTALL
	}
	err := client.ApplyHelmChart(meshplaykube.ApplyHelmChartConfig{
		Namespace:   "meshplay",
		ReleaseName: "meshplay-operator",
		ChartLocation: meshplaykube.HelmChartLocation{
			Repository: chartRepo,
			Chart:      chart,
			Version:    meshplayReleaseVersion,
		},
		// CreateNamespace doesn't have any effect when the action is UNINSTALL
		CreateNamespace: true,
		Action:          act,
		// Setting override values
		OverrideValues: overrides,
	})
	if err != nil {
		return err
	}
	return nil
}

func ConnectivityTest(clientName, externalIP string) bool {
	endpoint, err := url.Parse("http://" + externalIP + ":" + BrokerPingEndpoint)
	if err != nil {
		return false
	}

	resp, err := http.Get(endpoint.String())
	if err != nil {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var natsResponse Connections
	err = json.Unmarshal(body, &natsResponse)
	if err != nil {
		return false
	}

	for _, client := range natsResponse.Connections {
		if client.Name == clientName {
			return true
		}
	}
	return false
}
