package fullnode

import (
	"fmt"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/diff"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const maxP2PServiceDefault = int32(1)

// BuildServices returns a list of services given the crd.
//
// Creates a single RPC service, likely for use with an Ingress.
//
// Creates 1 p2p service per pod. P2P diverges from traditional web and kubernetes architecture which calls for a single
// p2p service backed by multiple pods.
// Pods may be in various states even with proper readiness probes.
// Therefore, we do not want to confuse or disrupt peer exchange (PEX) within CometBFT.
// If using a single p2p service, an outside peer discovering a pod out of sync it could be
// interpreted as byzantine behavior if the peer previously connected to a pod that was in sync through the same
// external address.
func BuildServices(crd *cosmosv1.CosmosFullNode) []diff.Resource[*corev1.Service] {
	max := maxP2PServiceDefault
	if v := crd.Spec.Service.MaxP2PExternalAddresses; v != nil {
		max = *v
	}
	maxExternal := lo.Clamp(max, 0, crd.Spec.Replicas)
	p2ps := make([]diff.Resource[*corev1.Service], crd.Spec.Replicas)

	for i := int32(0); i < crd.Spec.Replicas; i++ {
		ordinal := i
		var svc corev1.Service
		svc.Name = p2pServiceName(crd, ordinal)
		svc.Namespace = crd.Namespace
		svc.Kind = "Service"
		svc.APIVersion = "v1"

		svc.Labels = defaultLabels(crd,
			kube.InstanceLabel, instanceName(crd, ordinal),
			kube.ComponentLabel, "p2p",
		)
		svc.Annotations = map[string]string{}
		svc.Spec.Selector = map[string]string{kube.InstanceLabel: instanceName(crd, ordinal)}

		if i < maxExternal {
			preserveMergeInto(svc.Labels, crd.Spec.Service.P2PTemplate.Metadata.Labels)
			preserveMergeInto(svc.Annotations, crd.Spec.Service.P2PTemplate.Metadata.Annotations)
			svc.Spec.Type = *valOrDefault(crd.Spec.Service.P2PTemplate.Type, ptr(corev1.ServiceTypeLoadBalancer))
			svc.Spec.ExternalTrafficPolicy = *valOrDefault(crd.Spec.Service.P2PTemplate.ExternalTrafficPolicy, ptr(corev1.ServiceExternalTrafficPolicyTypeLocal))
		} else {
			svc.Spec.Type = corev1.ServiceTypeClusterIP
		}

		var (
			portNameP2P = "p2p"
		)

		servicePortList := []corev1.ServicePort{
			{
				Name:       portNameP2P,
				Protocol:   corev1.ProtocolTCP,
				Port:       p2pPort,
				TargetPort: intstr.FromString(portNameP2P),
			},
		}

		for dp := 0; dp < len(servicePortList); dp++ {
			n := servicePortList[dp].Name
			for _, p := range crd.Spec.Service.P2PTemplate.Ports {
				if p.Name == fmt.Sprintf("%s.%s", p2pServiceName(crd, ordinal), portNameP2P) {
					servicePortList[dp] = p
					servicePortList[dp].Name = n
				}
			}
		}

		svc.Spec.Ports = servicePortList

		p2ps[i] = diff.Adapt(&svc, i)
	}

	rpc := rpcService(crd)

	return append(p2ps, diff.Adapt(rpc, len(p2ps)))
}

func rpcService(crd *cosmosv1.CosmosFullNode) *corev1.Service {
	var svc corev1.Service
	svc.Name = rpcServiceName(crd)
	svc.Namespace = crd.Namespace
	svc.Kind = "Service"
	svc.APIVersion = "v1"
	svc.Labels = defaultLabels(crd,
		kube.ComponentLabel, "rpc",
	)
	svc.Annotations = map[string]string{}

	svc.Spec.Selector = map[string]string{kube.NameLabel: appName(crd)}
	svc.Spec.Type = corev1.ServiceTypeClusterIP

	rpcSpec := crd.Spec.Service.RPCTemplate
	preserveMergeInto(svc.Labels, rpcSpec.Metadata.Labels)
	preserveMergeInto(svc.Annotations, rpcSpec.Metadata.Annotations)
	kube.NormalizeMetadata(&svc.ObjectMeta)

	var (
		portNameAPI     = "api"
		portNameRosetta = "rosetta"
		portNameGrpc    = "grpc"
		portNameRPC     = "rpc"
		portNameGrpcWeb = "grpc-web"
	)

	servicePortList := []corev1.ServicePort{
		{
			Name:       portNameAPI,
			Protocol:   corev1.ProtocolTCP,
			Port:       apiPort,
			TargetPort: intstr.FromString(portNameAPI),
		},
		{
			Name:       portNameRosetta,
			Protocol:   corev1.ProtocolTCP,
			Port:       rosettaPort,
			TargetPort: intstr.FromString(portNameRosetta),
		},
		{
			Name:       portNameGrpc,
			Protocol:   corev1.ProtocolTCP,
			Port:       grpcPort,
			TargetPort: intstr.FromString(portNameGrpc),
		},
		{
			Name:       portNameRPC,
			Protocol:   corev1.ProtocolTCP,
			Port:       rpcPort,
			TargetPort: intstr.FromString(portNameRPC),
		},
		{
			Name:       portNameGrpcWeb,
			Protocol:   corev1.ProtocolTCP,
			Port:       grpcWebPort,
			TargetPort: intstr.FromString(portNameGrpcWeb),
		},
	}

	for i := 0; i < len(servicePortList); i++ {
		n := servicePortList[i].Name
		for _, p := range crd.Spec.Service.RPCTemplate.Ports {
			if p.Name == fmt.Sprintf("%s.%s", rpcServiceName(crd), n) {
				servicePortList[i] = p
				servicePortList[i].Name = n
			}
		}
	}

	if v := rpcSpec.ExternalTrafficPolicy; v != nil {
		svc.Spec.ExternalTrafficPolicy = *v
	}
	if v := rpcSpec.Type; v != nil {
		svc.Spec.Type = *v
	}

	svc.Spec.Ports = servicePortList

	return &svc
}

func p2pServiceName(crd *cosmosv1.CosmosFullNode, ordinal int32) string {
	return fmt.Sprintf("%s-p2p-%d", appName(crd), ordinal)
}

func rpcServiceName(crd *cosmosv1.CosmosFullNode) string {
	return fmt.Sprintf("%s-rpc", appName(crd))
}
