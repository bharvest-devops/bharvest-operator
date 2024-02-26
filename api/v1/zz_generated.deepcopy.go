//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2024 B-Harvest Corporation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoDataSource) DeepCopyInto(out *AutoDataSource) {
	*out = *in
	if in.VolumeSnapshotSelector != nil {
		in, out := &in.VolumeSnapshotSelector, &out.VolumeSnapshotSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoDataSource.
func (in *AutoDataSource) DeepCopy() *AutoDataSource {
	if in == nil {
		return nil
	}
	out := new(AutoDataSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChainSpec) DeepCopyInto(out *ChainSpec) {
	*out = *in
	if in.Comet != nil {
		in, out := &in.Comet, &out.Comet
		*out = new(CometConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.CosmosSDK != nil {
		in, out := &in.CosmosSDK, &out.CosmosSDK
		*out = new(SDKAppConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Namada != nil {
		in, out := &in.Namada, &out.Namada
		*out = new(NamadaConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.LogLevel != nil {
		in, out := &in.LogLevel, &out.LogLevel
		*out = new(string)
		**out = **in
	}
	if in.LogFormat != nil {
		in, out := &in.LogFormat, &out.LogFormat
		*out = new(string)
		**out = **in
	}
	if in.AddrbookURL != nil {
		in, out := &in.AddrbookURL, &out.AddrbookURL
		*out = new(string)
		**out = **in
	}
	if in.AddrbookScript != nil {
		in, out := &in.AddrbookScript, &out.AddrbookScript
		*out = new(string)
		**out = **in
	}
	if in.GenesisURL != nil {
		in, out := &in.GenesisURL, &out.GenesisURL
		*out = new(string)
		**out = **in
	}
	if in.GenesisScript != nil {
		in, out := &in.GenesisScript, &out.GenesisScript
		*out = new(string)
		**out = **in
	}
	if in.PrivvalSleepSeconds != nil {
		in, out := &in.PrivvalSleepSeconds, &out.PrivvalSleepSeconds
		*out = new(int32)
		**out = **in
	}
	if in.DatabaseBackend != nil {
		in, out := &in.DatabaseBackend, &out.DatabaseBackend
		*out = new(string)
		**out = **in
	}
	if in.Versions != nil {
		in, out := &in.Versions, &out.Versions
		*out = make([]ChainVersion, len(*in))
		copy(*out, *in)
	}
	if in.AdditionalInitArgs != nil {
		in, out := &in.AdditionalInitArgs, &out.AdditionalInitArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AdditionalStartArgs != nil {
		in, out := &in.AdditionalStartArgs, &out.AdditionalStartArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChainSpec.
func (in *ChainSpec) DeepCopy() *ChainSpec {
	if in == nil {
		return nil
	}
	out := new(ChainSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChainVersion) DeepCopyInto(out *ChainVersion) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChainVersion.
func (in *ChainVersion) DeepCopy() *ChainVersion {
	if in == nil {
		return nil
	}
	out := new(ChainVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CometConfig) DeepCopyInto(out *CometConfig) {
	*out = *in
	if in.RPC != nil {
		in, out := &in.RPC, &out.RPC
		*out = new(RPC)
		(*in).DeepCopyInto(*out)
	}
	if in.P2P != nil {
		in, out := &in.P2P, &out.P2P
		*out = new(P2P)
		(*in).DeepCopyInto(*out)
	}
	if in.Consensus != nil {
		in, out := &in.Consensus, &out.Consensus
		*out = new(Consensus)
		(*in).DeepCopyInto(*out)
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(Storage)
		(*in).DeepCopyInto(*out)
	}
	if in.TxIndex != nil {
		in, out := &in.TxIndex, &out.TxIndex
		*out = new(TxIndex)
		(*in).DeepCopyInto(*out)
	}
	if in.Instrumentation != nil {
		in, out := &in.Instrumentation, &out.Instrumentation
		*out = new(Instrumentation)
		(*in).DeepCopyInto(*out)
	}
	if in.Statesync != nil {
		in, out := &in.Statesync, &out.Statesync
		*out = new(Statesync)
		(*in).DeepCopyInto(*out)
	}
	if in.TomlOverrides != nil {
		in, out := &in.TomlOverrides, &out.TomlOverrides
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CometConfig.
func (in *CometConfig) DeepCopy() *CometConfig {
	if in == nil {
		return nil
	}
	out := new(CometConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Consensus) DeepCopyInto(out *Consensus) {
	*out = *in
	if in.DoubleSignCheckHeight != nil {
		in, out := &in.DoubleSignCheckHeight, &out.DoubleSignCheckHeight
		*out = new(uint64)
		**out = **in
	}
	if in.SkipTimeoutCommit != nil {
		in, out := &in.SkipTimeoutCommit, &out.SkipTimeoutCommit
		*out = new(bool)
		**out = **in
	}
	if in.CreateEmptyBlocks != nil {
		in, out := &in.CreateEmptyBlocks, &out.CreateEmptyBlocks
		*out = new(bool)
		**out = **in
	}
	if in.CreateEmptyBlocksInterval != nil {
		in, out := &in.CreateEmptyBlocksInterval, &out.CreateEmptyBlocksInterval
		*out = new(string)
		**out = **in
	}
	if in.PeerGossipSleepDuration != nil {
		in, out := &in.PeerGossipSleepDuration, &out.PeerGossipSleepDuration
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Consensus.
func (in *Consensus) DeepCopy() *Consensus {
	if in == nil {
		return nil
	}
	out := new(Consensus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CosmosFullNode) DeepCopyInto(out *CosmosFullNode) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CosmosFullNode.
func (in *CosmosFullNode) DeepCopy() *CosmosFullNode {
	if in == nil {
		return nil
	}
	out := new(CosmosFullNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CosmosFullNode) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CosmosFullNodeList) DeepCopyInto(out *CosmosFullNodeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CosmosFullNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CosmosFullNodeList.
func (in *CosmosFullNodeList) DeepCopy() *CosmosFullNodeList {
	if in == nil {
		return nil
	}
	out := new(CosmosFullNodeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CosmosFullNodeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FullNodeProbesSpec) DeepCopyInto(out *FullNodeProbesSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FullNodeProbesSpec.
func (in *FullNodeProbesSpec) DeepCopy() *FullNodeProbesSpec {
	if in == nil {
		return nil
	}
	out := new(FullNodeProbesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FullNodeSnapshotStatus) DeepCopyInto(out *FullNodeSnapshotStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FullNodeSnapshotStatus.
func (in *FullNodeSnapshotStatus) DeepCopy() *FullNodeSnapshotStatus {
	if in == nil {
		return nil
	}
	out := new(FullNodeSnapshotStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FullNodeSpec) DeepCopyInto(out *FullNodeSpec) {
	*out = *in
	in.ChainSpec.DeepCopyInto(&out.ChainSpec)
	in.PodTemplate.DeepCopyInto(&out.PodTemplate)
	in.RolloutStrategy.DeepCopyInto(&out.RolloutStrategy)
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	if in.RetentionPolicy != nil {
		in, out := &in.RetentionPolicy, &out.RetentionPolicy
		*out = new(RetentionPolicy)
		**out = **in
	}
	in.Service.DeepCopyInto(&out.Service)
	if in.InstanceOverrides != nil {
		in, out := &in.InstanceOverrides, &out.InstanceOverrides
		*out = make(map[string]InstanceOverridesSpec, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.SelfHeal != nil {
		in, out := &in.SelfHeal, &out.SelfHeal
		*out = new(SelfHealSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FullNodeSpec.
func (in *FullNodeSpec) DeepCopy() *FullNodeSpec {
	if in == nil {
		return nil
	}
	out := new(FullNodeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FullNodeStatus) DeepCopyInto(out *FullNodeStatus) {
	*out = *in
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
	if in.ScheduledSnapshotStatus != nil {
		in, out := &in.ScheduledSnapshotStatus, &out.ScheduledSnapshotStatus
		*out = make(map[string]FullNodeSnapshotStatus, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.SelfHealing.DeepCopyInto(&out.SelfHealing)
	if in.Peers != nil {
		in, out := &in.Peers, &out.Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SyncInfo != nil {
		in, out := &in.SyncInfo, &out.SyncInfo
		*out = make(map[string]*SyncInfoPodStatus, len(*in))
		for key, val := range *in {
			var outVal *SyncInfoPodStatus
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SyncInfoPodStatus)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.Height != nil {
		in, out := &in.Height, &out.Height
		*out = make(map[string]uint64, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FullNodeStatus.
func (in *FullNodeStatus) DeepCopy() *FullNodeStatus {
	if in == nil {
		return nil
	}
	out := new(FullNodeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeightDriftMitigationSpec) DeepCopyInto(out *HeightDriftMitigationSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeightDriftMitigationSpec.
func (in *HeightDriftMitigationSpec) DeepCopy() *HeightDriftMitigationSpec {
	if in == nil {
		return nil
	}
	out := new(HeightDriftMitigationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceOverridesSpec) DeepCopyInto(out *InstanceOverridesSpec) {
	*out = *in
	if in.DisableStrategy != nil {
		in, out := &in.DisableStrategy, &out.DisableStrategy
		*out = new(DisableStrategy)
		**out = **in
	}
	if in.VolumeClaimTemplate != nil {
		in, out := &in.VolumeClaimTemplate, &out.VolumeClaimTemplate
		*out = new(PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ExternalAddress != nil {
		in, out := &in.ExternalAddress, &out.ExternalAddress
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceOverridesSpec.
func (in *InstanceOverridesSpec) DeepCopy() *InstanceOverridesSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceOverridesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Instrumentation) DeepCopyInto(out *Instrumentation) {
	*out = *in
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(bool)
		**out = **in
	}
	if in.PrometheusListenAddr != nil {
		in, out := &in.PrometheusListenAddr, &out.PrometheusListenAddr
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Instrumentation.
func (in *Instrumentation) DeepCopy() *Instrumentation {
	if in == nil {
		return nil
	}
	out := new(Instrumentation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metadata) DeepCopyInto(out *Metadata) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metadata.
func (in *Metadata) DeepCopy() *Metadata {
	if in == nil {
		return nil
	}
	out := new(Metadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamadaConfig) DeepCopyInto(out *NamadaConfig) {
	*out = *in
	if in.WasmDir != nil {
		in, out := &in.WasmDir, &out.WasmDir
		*out = new(string)
		**out = **in
	}
	if in.Ledger != nil {
		in, out := &in.Ledger, &out.Ledger
		*out = new(NamadaLedger)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamadaConfig.
func (in *NamadaConfig) DeepCopy() *NamadaConfig {
	if in == nil {
		return nil
	}
	out := new(NamadaConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamadaEthereumBridge) DeepCopyInto(out *NamadaEthereumBridge) {
	*out = *in
	if in.Mode != nil {
		in, out := &in.Mode, &out.Mode
		*out = new(string)
		**out = **in
	}
	if in.OracleRPCEndpoint != nil {
		in, out := &in.OracleRPCEndpoint, &out.OracleRPCEndpoint
		*out = new(string)
		**out = **in
	}
	if in.ChannelBufferSize != nil {
		in, out := &in.ChannelBufferSize, &out.ChannelBufferSize
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamadaEthereumBridge.
func (in *NamadaEthereumBridge) DeepCopy() *NamadaEthereumBridge {
	if in == nil {
		return nil
	}
	out := new(NamadaEthereumBridge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamadaLedger) DeepCopyInto(out *NamadaLedger) {
	*out = *in
	if in.Shell != nil {
		in, out := &in.Shell, &out.Shell
		*out = new(NamadaShell)
		(*in).DeepCopyInto(*out)
	}
	if in.EthereumBridge != nil {
		in, out := &in.EthereumBridge, &out.EthereumBridge
		*out = new(NamadaEthereumBridge)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamadaLedger.
func (in *NamadaLedger) DeepCopy() *NamadaLedger {
	if in == nil {
		return nil
	}
	out := new(NamadaLedger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamadaShell) DeepCopyInto(out *NamadaShell) {
	*out = *in
	if in.BaseDir != nil {
		in, out := &in.BaseDir, &out.BaseDir
		*out = new(string)
		**out = **in
	}
	if in.StorageReadPastHeightLimit != nil {
		in, out := &in.StorageReadPastHeightLimit, &out.StorageReadPastHeightLimit
		*out = new(uint64)
		**out = **in
	}
	if in.DbDir != nil {
		in, out := &in.DbDir, &out.DbDir
		*out = new(string)
		**out = **in
	}
	if in.TendermintMode != nil {
		in, out := &in.TendermintMode, &out.TendermintMode
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamadaShell.
func (in *NamadaShell) DeepCopy() *NamadaShell {
	if in == nil {
		return nil
	}
	out := new(NamadaShell)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *P2P) DeepCopyInto(out *P2P) {
	*out = *in
	if in.Laddr != nil {
		in, out := &in.Laddr, &out.Laddr
		*out = new(string)
		**out = **in
	}
	if in.ExternalAddress != nil {
		in, out := &in.ExternalAddress, &out.ExternalAddress
		*out = new(string)
		**out = **in
	}
	if in.Seeds != nil {
		in, out := &in.Seeds, &out.Seeds
		*out = new(string)
		**out = **in
	}
	if in.PersistentPeers != nil {
		in, out := &in.PersistentPeers, &out.PersistentPeers
		*out = new(string)
		**out = **in
	}
	if in.MaxNumInboundPeers != nil {
		in, out := &in.MaxNumInboundPeers, &out.MaxNumInboundPeers
		*out = new(int32)
		**out = **in
	}
	if in.MaxNumOutboundPeers != nil {
		in, out := &in.MaxNumOutboundPeers, &out.MaxNumOutboundPeers
		*out = new(int32)
		**out = **in
	}
	if in.Pex != nil {
		in, out := &in.Pex, &out.Pex
		*out = new(bool)
		**out = **in
	}
	if in.SeedMode != nil {
		in, out := &in.SeedMode, &out.SeedMode
		*out = new(bool)
		**out = **in
	}
	if in.PrivatePeerIds != nil {
		in, out := &in.PrivatePeerIds, &out.PrivatePeerIds
		*out = new(string)
		**out = **in
	}
	if in.UnconditionalPeerIDs != nil {
		in, out := &in.UnconditionalPeerIDs, &out.UnconditionalPeerIDs
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new P2P.
func (in *P2P) DeepCopy() *P2P {
	if in == nil {
		return nil
	}
	out := new(P2P)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PVCAutoScaleSpec) DeepCopyInto(out *PVCAutoScaleSpec) {
	*out = *in
	out.MaxSize = in.MaxSize.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PVCAutoScaleSpec.
func (in *PVCAutoScaleSpec) DeepCopy() *PVCAutoScaleSpec {
	if in == nil {
		return nil
	}
	out := new(PVCAutoScaleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PVCAutoScaleStatus) DeepCopyInto(out *PVCAutoScaleStatus) {
	*out = *in
	out.RequestedSize = in.RequestedSize.DeepCopy()
	in.RequestedAt.DeepCopyInto(&out.RequestedAt)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PVCAutoScaleStatus.
func (in *PVCAutoScaleStatus) DeepCopy() *PVCAutoScaleStatus {
	if in == nil {
		return nil
	}
	out := new(PVCAutoScaleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistentVolumeClaimSpec) DeepCopyInto(out *PersistentVolumeClaimSpec) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	in.Resources.DeepCopyInto(&out.Resources)
	if in.AccessModes != nil {
		in, out := &in.AccessModes, &out.AccessModes
		*out = make([]corev1.PersistentVolumeAccessMode, len(*in))
		copy(*out, *in)
	}
	if in.VolumeMode != nil {
		in, out := &in.VolumeMode, &out.VolumeMode
		*out = new(corev1.PersistentVolumeMode)
		**out = **in
	}
	if in.DataSource != nil {
		in, out := &in.DataSource, &out.DataSource
		*out = new(corev1.TypedLocalObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.AutoDataSource != nil {
		in, out := &in.AutoDataSource, &out.AutoDataSource
		*out = new(AutoDataSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistentVolumeClaimSpec.
func (in *PersistentVolumeClaimSpec) DeepCopy() *PersistentVolumeClaimSpec {
	if in == nil {
		return nil
	}
	out := new(PersistentVolumeClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSpec) DeepCopyInto(out *PodSpec) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]map[string]string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]corev1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(corev1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Priority != nil {
		in, out := &in.Priority, &out.Priority
		*out = new(int32)
		**out = **in
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.TerminationGracePeriodSeconds != nil {
		in, out := &in.TerminationGracePeriodSeconds, &out.TerminationGracePeriodSeconds
		*out = new(int64)
		**out = **in
	}
	out.Probes = in.Probes
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]corev1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]corev1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]corev1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSpec.
func (in *PodSpec) DeepCopy() *PodSpec {
	if in == nil {
		return nil
	}
	out := new(PodSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pruning) DeepCopyInto(out *Pruning) {
	*out = *in
	if in.Interval != nil {
		in, out := &in.Interval, &out.Interval
		*out = new(uint32)
		**out = **in
	}
	if in.KeepEvery != nil {
		in, out := &in.KeepEvery, &out.KeepEvery
		*out = new(uint32)
		**out = **in
	}
	if in.KeepRecent != nil {
		in, out := &in.KeepRecent, &out.KeepRecent
		*out = new(uint32)
		**out = **in
	}
	if in.MinRetainBlocks != nil {
		in, out := &in.MinRetainBlocks, &out.MinRetainBlocks
		*out = new(uint32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pruning.
func (in *Pruning) DeepCopy() *Pruning {
	if in == nil {
		return nil
	}
	out := new(Pruning)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RPC) DeepCopyInto(out *RPC) {
	*out = *in
	if in.Laddr != nil {
		in, out := &in.Laddr, &out.Laddr
		*out = new(string)
		**out = **in
	}
	if in.CorsAllowedOrigins != nil {
		in, out := &in.CorsAllowedOrigins, &out.CorsAllowedOrigins
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.CorsAllowedMethods != nil {
		in, out := &in.CorsAllowedMethods, &out.CorsAllowedMethods
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.TimeoutBroadcastTxCommit != nil {
		in, out := &in.TimeoutBroadcastTxCommit, &out.TimeoutBroadcastTxCommit
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RPC.
func (in *RPC) DeepCopy() *RPC {
	if in == nil {
		return nil
	}
	out := new(RPC)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RolloutStrategy) DeepCopyInto(out *RolloutStrategy) {
	*out = *in
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RolloutStrategy.
func (in *RolloutStrategy) DeepCopy() *RolloutStrategy {
	if in == nil {
		return nil
	}
	out := new(RolloutStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SDKAppConfig) DeepCopyInto(out *SDKAppConfig) {
	*out = *in
	if in.SnapshotURL != nil {
		in, out := &in.SnapshotURL, &out.SnapshotURL
		*out = new(string)
		**out = **in
	}
	if in.SnapshotScript != nil {
		in, out := &in.SnapshotScript, &out.SnapshotScript
		*out = new(string)
		**out = **in
	}
	if in.Pruning != nil {
		in, out := &in.Pruning, &out.Pruning
		*out = new(Pruning)
		(*in).DeepCopyInto(*out)
	}
	if in.HaltHeight != nil {
		in, out := &in.HaltHeight, &out.HaltHeight
		*out = new(uint64)
		**out = **in
	}
	if in.TomlOverrides != nil {
		in, out := &in.TomlOverrides, &out.TomlOverrides
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SDKAppConfig.
func (in *SDKAppConfig) DeepCopy() *SDKAppConfig {
	if in == nil {
		return nil
	}
	out := new(SDKAppConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SelfHealSpec) DeepCopyInto(out *SelfHealSpec) {
	*out = *in
	if in.PVCAutoScale != nil {
		in, out := &in.PVCAutoScale, &out.PVCAutoScale
		*out = new(PVCAutoScaleSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.HeightDriftMitigation != nil {
		in, out := &in.HeightDriftMitigation, &out.HeightDriftMitigation
		*out = new(HeightDriftMitigationSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SelfHealSpec.
func (in *SelfHealSpec) DeepCopy() *SelfHealSpec {
	if in == nil {
		return nil
	}
	out := new(SelfHealSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SelfHealingStatus) DeepCopyInto(out *SelfHealingStatus) {
	*out = *in
	if in.PVCAutoScale != nil {
		in, out := &in.PVCAutoScale, &out.PVCAutoScale
		*out = make(map[string]*PVCAutoScaleStatus, len(*in))
		for key, val := range *in {
			var outVal *PVCAutoScaleStatus
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(PVCAutoScaleStatus)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SelfHealingStatus.
func (in *SelfHealingStatus) DeepCopy() *SelfHealingStatus {
	if in == nil {
		return nil
	}
	out := new(SelfHealingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOverridesSpec) DeepCopyInto(out *ServiceOverridesSpec) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(corev1.ServiceType)
		**out = **in
	}
	if in.ExternalTrafficPolicy != nil {
		in, out := &in.ExternalTrafficPolicy, &out.ExternalTrafficPolicy
		*out = new(corev1.ServiceExternalTrafficPolicy)
		**out = **in
	}
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]corev1.ServicePort, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOverridesSpec.
func (in *ServiceOverridesSpec) DeepCopy() *ServiceOverridesSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceOverridesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceSpec) DeepCopyInto(out *ServiceSpec) {
	*out = *in
	if in.MaxP2PExternalAddresses != nil {
		in, out := &in.MaxP2PExternalAddresses, &out.MaxP2PExternalAddresses
		*out = new(int32)
		**out = **in
	}
	in.P2PTemplate.DeepCopyInto(&out.P2PTemplate)
	in.RPCTemplate.DeepCopyInto(&out.RPCTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSpec.
func (in *ServiceSpec) DeepCopy() *ServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Statesync) DeepCopyInto(out *Statesync) {
	*out = *in
	if in.Enable != nil {
		in, out := &in.Enable, &out.Enable
		*out = new(bool)
		**out = **in
	}
	if in.RPCServers != nil {
		in, out := &in.RPCServers, &out.RPCServers
		*out = new(string)
		**out = **in
	}
	if in.TrustHeight != nil {
		in, out := &in.TrustHeight, &out.TrustHeight
		*out = new(uint64)
		**out = **in
	}
	if in.TrustHash != nil {
		in, out := &in.TrustHash, &out.TrustHash
		*out = new(string)
		**out = **in
	}
	if in.TrustPeriod != nil {
		in, out := &in.TrustPeriod, &out.TrustPeriod
		*out = new(string)
		**out = **in
	}
	if in.DiscoveryTime != nil {
		in, out := &in.DiscoveryTime, &out.DiscoveryTime
		*out = new(string)
		**out = **in
	}
	if in.TempDir != nil {
		in, out := &in.TempDir, &out.TempDir
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Statesync.
func (in *Statesync) DeepCopy() *Statesync {
	if in == nil {
		return nil
	}
	out := new(Statesync)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Storage) DeepCopyInto(out *Storage) {
	*out = *in
	if in.DiscardAbciResponses != nil {
		in, out := &in.DiscardAbciResponses, &out.DiscardAbciResponses
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Storage.
func (in *Storage) DeepCopy() *Storage {
	if in == nil {
		return nil
	}
	out := new(Storage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncInfoPodStatus) DeepCopyInto(out *SyncInfoPodStatus) {
	*out = *in
	in.Timestamp.DeepCopyInto(&out.Timestamp)
	if in.Height != nil {
		in, out := &in.Height, &out.Height
		*out = new(uint64)
		**out = **in
	}
	if in.InSync != nil {
		in, out := &in.InSync, &out.InSync
		*out = new(bool)
		**out = **in
	}
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncInfoPodStatus.
func (in *SyncInfoPodStatus) DeepCopy() *SyncInfoPodStatus {
	if in == nil {
		return nil
	}
	out := new(SyncInfoPodStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TxIndex) DeepCopyInto(out *TxIndex) {
	*out = *in
	if in.Indexer != nil {
		in, out := &in.Indexer, &out.Indexer
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TxIndex.
func (in *TxIndex) DeepCopy() *TxIndex {
	if in == nil {
		return nil
	}
	out := new(TxIndex)
	in.DeepCopyInto(out)
	return out
}
