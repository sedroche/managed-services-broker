package fuse

import (
	"net/http"

	brokerapi "github.com/aerogear/managed-services-broker/pkg/broker"
	"github.com/pkg/errors"
	glog "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type FuseDeployer struct {
	id string
}

func NewDeployer(id string) *FuseDeployer {
	return &FuseDeployer{id: id}
}

func (fd *FuseDeployer) DoesDeploy(serviceID string) bool {
	return serviceID == "fuse-service-id"
}

func (fd *FuseDeployer) GetCatalogEntries() []*brokerapi.Service {
	glog.Infof("Getting fuse catalog entries")
	return []*brokerapi.Service{
		{
			Name:        "fuse",
			ID:          "fuse-service-id",
			Description: "fuse",
			Metadata:    map[string]string{"serviceName": "fuse", "serviceType": "fuse"},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					Name:        "default",
					ID:          "default",
					Description: "default fuse plan",
					Free:        true,
					Schemas: &brokerapi.Schemas{
						ServiceBinding: &brokerapi.ServiceBindingSchema{
							Create: &brokerapi.RequestResponseSchema{},
						},
						ServiceInstance: &brokerapi.ServiceInstanceSchema{
							Create: &brokerapi.InputParametersSchema{},
						},
					},
				},
			},
		},
	}
}

func (fd *FuseDeployer) Deploy(id string, k8sclient kubernetes.Interface, config *rest.Config) (*brokerapi.CreateServiceInstanceResponse, error) {
	ns, err := k8sclient.CoreV1().Namespaces().Create(getNamespace("fuse-" + id))
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create namespace for fuse service")
	}
	glog.Infof("created namespace: %s", ns.ObjectMeta.Name)

	glog.Infof("deploying fuse from deployer, id: %s", id)
	return &brokerapi.CreateServiceInstanceResponse{
		Code: http.StatusOK,
	}, nil
}

func (fd *FuseDeployer) GetID() string {
	glog.Infof("getting fuse id from deployer: " + fd.id)
	return fd.id
}
