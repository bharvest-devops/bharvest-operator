package fullnode

import (
	"fmt"
	cosmosv1 "github.com/bharvest-devops/cosmos-operator/api/v1"
	"github.com/bharvest-devops/cosmos-operator/internal/diff"
	"github.com/bharvest-devops/cosmos-operator/internal/kube"
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
	svcs := make([]diff.Resource[*corev1.Service], crd.Spec.Replicas)

	for i := int32(0); i < crd.Spec.Replicas; i++ {
		ordinal := i
		svc, componentName := new(corev1.Service), "p2p"

		p2pServiceSpec := findP2PServiceSpec(crd, i)

		if p2pServiceSpec != nil {
			svc = podService(crd, ordinal,
				componentName,
				p2pServiceSpec.Metadata.Labels,
				p2pServiceSpec.Metadata.Annotations,
				*valOrDefault(p2pServiceSpec.Type, ptr(corev1.ServiceTypeLoadBalancer)),
				valOrDefault(p2pServiceSpec.ExternalTrafficPolicy, ptr(corev1.ServiceExternalTrafficPolicyTypeLocal)))
		} else {
			svc = podService(crd, ordinal, componentName, map[string]string{}, map[string]string{}, corev1.ServiceTypeClusterIP, nil)
		}

		servicePortList := []corev1.ServicePort{
			{
				Name:       componentName,
				Protocol:   corev1.ProtocolTCP,
				Port:       p2pPort,
				TargetPort: intstr.FromString(componentName),
			},
		}

		for dp := 0; dp < len(servicePortList); dp++ {
			if p2pServiceSpec != nil {
				if p2pServiceSpec.Port != *new(int32) {
					servicePortList[dp].Port = p2pServiceSpec.Port
				}
				if p2pServiceSpec.Protocol != corev1.ProtocolTCP {
					servicePortList[dp].Protocol = p2pServiceSpec.Protocol
				}
				if p2pServiceSpec.NodePort != *new(int32) {
					servicePortList[dp].NodePort = p2pServiceSpec.NodePort
				}
			}
		}

		svc.Spec.Ports = servicePortList

		svcs[i] = diff.Adapt(svc, i)
	}

	rpc := rpcService(crd)

	if crd.Spec.Type == cosmosv1.Sentry {
		for i := int32(0); i < crd.Spec.Replicas; i++ {
			ordinal := i
			svc, componentName := new(corev1.Service), "cosmos-sentry"

			svc = podService(crd, ordinal,
				componentName,
				map[string]string{},
				map[string]string{},
				corev1.ServiceTypeClusterIP,
				nil)

			servicePortList := []corev1.ServicePort{
				{
					Name:       "sentry-privval",
					Protocol:   corev1.ProtocolTCP,
					Port:       privvalPort,
					TargetPort: intstr.FromString("privval"),
				},
			}

			svc.Spec.Ports = servicePortList
			svc.Spec.PublishNotReadyAddresses = true

			svcs[i] = diff.Adapt(svc, i)
		}
	}

	return append(svcs, diff.Adapt(rpc, len(svcs)))
}

func podService(crd *cosmosv1.CosmosFullNode, ordinal int32, componentName string, labels, annotations map[string]string, serviceType corev1.ServiceType, externalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType) *corev1.Service {
	svc := new(corev1.Service)
	svc.Name = podServiceName(crd, componentName, ordinal)
	svc.Namespace = crd.Namespace
	svc.Kind = "Service"
	svc.APIVersion = "v1"

	svc.Labels = defaultLabels(crd,
		kube.InstanceLabel, instanceName(crd, ordinal),
		kube.ComponentLabel, componentName,
	)
	svc.Annotations = map[string]string{}
	svc.Spec.Selector = map[string]string{kube.InstanceLabel: instanceName(crd, ordinal)}

	preserveMergeInto(svc.Labels, labels)
	preserveMergeInto(svc.Annotations, annotations)

	svc.Spec.Type = serviceType
	// To set svc.Spec.ExternalTrafficPolicy, you should set type as NodePort, or LoadBalancer
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer || svc.Spec.Type == corev1.ServiceTypeNodePort {
		svc.Spec.ExternalTrafficPolicy = *valOrDefault(externalTrafficPolicy, ptr(corev1.ServiceExternalTrafficPolicyTypeLocal))
	}

	return svc
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
		for _, p := range rpcSpec.Ports {

			// Prevents error occurrence from wrong request when not configured nodePort but entered nodePort
			if *rpcSpec.Type != corev1.ServiceTypeNodePort && p.NodePort != *new(int32) {
				p.NodePort = *new(int32)
			}
			if p.Name == n {
				servicePortList[i] = p
				servicePortList[i].Name = n
			}
			break
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

func podServiceName(crd *cosmosv1.CosmosFullNode, componentName string, ordinal int32) string {
	return fmt.Sprintf("%s-%s-%d", appName(crd), componentName, ordinal)
}

func rpcServiceName(crd *cosmosv1.CosmosFullNode) string {
	return fmt.Sprintf("%s-rpc", appName(crd))
}

func findP2PServiceSpec(crd *cosmosv1.CosmosFullNode, idx int32) *cosmosv1.P2PServiceSpec {

	for _, p := range crd.Spec.Service.P2PServiceSpecs {
		if *p.PodIdx == uint32(idx) {
			return &p
		}
	}

	return nil
}
