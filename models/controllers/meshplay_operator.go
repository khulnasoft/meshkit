package controllers

import (
	"context"
	"fmt"
	"sync"

	meshplaykube "github.com/khulnasoft/meshkit/utils/kubernetes"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/polymorphichelpers"
)

type meshplayOperator struct {
	name           string
	status         MeshplayControllerStatus
	client         *meshplaykube.Client
	deploymentConf OperatorDeploymentConfig
	mx             sync.Mutex
}

type OperatorDeploymentConfig struct {
	GetHelmOverrides      func(delete bool) map[string]interface{}
	HelmChartRepo         string
	MeshplayReleaseVersion string
}

func NewMeshplayOperatorHandler(client *meshplaykube.Client, deploymentConf OperatorDeploymentConfig) IMeshplayController {
	return &meshplayOperator{
		name:           "MeshplayOperator",
		status:         Unknown,
		client:         client,
		deploymentConf: deploymentConf,
	}
}

func (mo *meshplayOperator) GetName() string {
	return mo.name
}

func (mo *meshplayOperator) GetStatus() MeshplayControllerStatus {
	if mo.status == Undeployed {
		return Undeployed
	}
	// check if the deployment exists
	deployment, err := mo.client.DynamicKubeClient.Resource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}).Namespace("meshplay").Get(context.TODO(), "meshplay-operator", metav1.GetOptions{})
	if err != nil {
		if kubeerror.IsNotFound(err) {
			mo.setStatus(NotDeployed)
			return mo.status
		}
		return Unknown
	}

	sv, err := polymorphichelpers.StatusViewerFor(deployment.GroupVersionKind().GroupKind())
	if err != nil {
		mo.setStatus(Unknown)
		return mo.status
	}
	_, done, err := sv.Status(deployment, 0)
	if err != nil {
		mo.setStatus(Unknown)
		return mo.status
	}
	if done {
		mo.setStatus(Deployed)
	} else {
		mo.setStatus(Deploying)
	}
	return mo.status
}

func (mo *meshplayOperator) Deploy(force bool) error {
	status := mo.GetStatus()
	if status == Undeployed && !force {
		return nil
	}
	if status == Deploying {
		return ErrDeployController(fmt.Errorf("Already a Meshplay Operator is being deployed."))
	}
	err := applyOperatorHelmChart(mo.deploymentConf.HelmChartRepo, *mo.client, mo.deploymentConf.MeshplayReleaseVersion, false, mo.deploymentConf.GetHelmOverrides(false))
	if err != nil {
		return ErrDeployController(err)
	}
	mo.setStatus(Deployed)
	return nil
}
func (mo *meshplayOperator) Undeploy() error {
	err := applyOperatorHelmChart(mo.deploymentConf.HelmChartRepo, *mo.client, mo.deploymentConf.MeshplayReleaseVersion, true, mo.deploymentConf.GetHelmOverrides(false))
	if err != nil {
		return ErrDeployController(err)
	}
	mo.setStatus(Undeployed)
	return nil
}

func (mo *meshplayOperator) GetPublicEndpoint() (string, error) {
	return "", nil
}

func (mo *meshplayOperator) GetVersion() (string, error) {
	return "", nil
}

func (mo *meshplayOperator) setStatus(st MeshplayControllerStatus) {
	mo.mx.Lock()
	defer mo.mx.Unlock()
	mo.status = st
}
