# API Reference

## Packages
- [cosmos.bharvest/v1](#cosmosbharvestv1)
- [cosmos.bharvest/v1alpha1](#cosmosbharvestv1alpha1)


## cosmos.bharvest/v1

Package v1 contains API Schema definitions for the cosmos v1 API group

### Resource Types
- [CosmosFullNode](#cosmosfullnode)
- [CosmosFullNodeList](#cosmosfullnodelist)



#### AutoDataSource





_Appears in:_
- [PersistentVolumeClaimSpec](#persistentvolumeclaimspec)

| Field | Description |
| --- | --- |
| `volumeSnapshotSelector` _object (keys:string, values:string)_ | If set, chooses the most recent VolumeSnapshot matching the selector to use as the PVC dataSource.<br /><br />See ScheduledVolumeSnapshot for a means of creating periodic VolumeSnapshots.<br /><br />The VolumeSnapshots must be in the same namespace as the CosmosFullNode.<br /><br />If no VolumeSnapshots found, controller logs error and still creates PVC. |
| `matchInstance` _boolean_ | If true, the volume snapshot selector will make sure the PVC<br /><br />is restored from a VolumeSnapshot on the same node.<br /><br />This is useful if the VolumeSnapshots are local to the node, e.g. for topolvm. |


#### ChainSpec





_Appears in:_
- [FullNodeSpec](#fullnodespec)

| Field | Description |
| --- | --- |
| `chainID` _string_ | Genesis file chain-id. |
| `chainType` _string_ | Describes chain type to operate<br /><br />If not set, defaults to "cosmos". |
| `network` _string_ | The network environment. Typically, mainnet, testnet, devnet, etc. |
| `binary` _string_ | Binary name which runs commands. E.g. gaiad, junod, osmosisd |
| `homeDir` _string_ | The chain's home directory is where the chain's data and config is stored.<br /><br />This should be a single folder. E.g. .gaia, .dydxprotocol, .osmosisd, etc.<br /><br />Set via --home flag when running the binary.<br /><br />If empty, defaults to "cosmos" which translates to `chain start --home /home/operator/cosmos`.<br /><br />Historically, several chains do not respect the --home and save data outside --home which crashes the pods.<br /><br />Therefore, this option was introduced to mitigate those edge cases, so that you can specify the home directory<br /><br />to match the chain's default home dir. |
| `config` _[CometBFTConfig](#cometconfig)_ | CometBFT (formerly Tendermint) configuration applied to config.toml.<br /><br />Although optional, it's highly recommended you configure this field. |
| `cosmos` _[SDKAppConfig](#sdkappconfig)_ | CosmosSDK configuration applied to app.toml. |
| `namada` _[NamadaConfig](#namadaconfig)_ | Namada configuration applied to $CHAIN_ID/config.toml. |
| `logLevel` _string_ | One of trace\|debug\|info\|warn\|error\|fatal\|panic.<br /><br />If not set, defaults to info. |
| `logFormat` _string_ | One of plain or json.<br /><br />If not set, defaults to plain. |
| `addrbookURL` _string_ | URL to address book file to download from the internet.<br /><br />The operator detects and properly handles the following file extensions:<br /><br />.json, .json.gz, .tar, .tar.gz, .tar.gzip, .zip<br /><br />Use AddrbookScript if the chain has an unconventional file format or address book location. |
| `addrbookScript` _string_ | Specify shell (sh) script commands to properly download and save the address book file.<br /><br />Prefer AddrbookURL if the file is in a conventional format.<br /><br />The available shell commands are from docker image ghcr.io/strangelove-ventures/infra-toolkit, including wget and curl.<br /><br />Save the file to env var $ADDRBOOK_FILE.<br /><br />E.g. curl https://url-to-addrbook.com > $ADDRBOOK_FILE<br /><br />Takes precedence over AddrbookURL.<br /><br />Hint: Use "set -eux" in your script.<br /><br />Available env vars:<br /><br />$HOME: The home directory.<br /><br />$ADDRBOOK_FILE: The location of the final address book file.<br /><br />$CONFIG_DIR: The location of the config dir that houses the address book file. Used for extracting from archives. The archive must have a single file called "addrbook.json". |
| `genesisURL` _string_ | URL to genesis file to download from the internet.<br /><br />Although this field is optional, you will almost always want to set it.<br /><br />If not set, uses the genesis file created from the init subcommand. (This behavior may be desirable for new chains or testing.)<br /><br />The operator detects and properly handles the following file extensions:<br /><br />.json, .json.gz, .tar, .tar.gz, .tar.gzip, .zip<br /><br />Use GenesisScript if the chain has an unconventional file format or genesis location. |
| `genesisScript` _string_ | Specify shell (sh) script commands to properly download and save the genesis file.<br /><br />Prefer GenesisURL if the file is in a conventional format.<br /><br />The available shell commands are from docker image ghcr.io/strangelove-ventures/infra-toolkit, including wget and curl.<br /><br />Save the file to env var $GENESIS_FILE.<br /><br />E.g. curl https://url-to-genesis.com \| jq '.genesis' > $GENESIS_FILE<br /><br />Takes precedence over GenesisURL.<br /><br />Hint: Use "set -eux" in your script.<br /><br />Available env vars:<br /><br />$HOME: The home directory.<br /><br />$GENESIS_FILE: The location of the final genesis file.<br /><br />$CONFIG_DIR: The location of the config dir that houses the genesis file. Used for extracting from archives. The archive must have a single file called "genesis.json". |
| `privvalSleepSeconds` _integer_ | If configured as a Sentry, invokes sleep command with this value before running chain start command.<br /><br />Currently, requires the privval laddr to be available immediately without any retry.<br /><br />This workaround gives time for the connection to be made to a remote signer.<br /><br />If a Sentry and not set, defaults to 10.<br /><br />If set to 0, omits injecting sleep command.<br /><br />Assumes chain image has `sleep` in $PATH. |
| `databaseBackend` _string_ | DatabaseBackend must match in order to detect the block height<br /><br />of the chain prior to starting in order to pick the correct image version.<br /><br />options: goleveldb, rocksdb, pebbledb<br /><br />Defaults to goleveldb. |
| `versions` _[ChainVersion](#chainversion) array_ | Versions of the chain and which height they should be applied.<br /><br />When provided, the operator will automatically upgrade the chain as it reaches the specified heights.<br /><br />If not provided, the operator will not upgrade the chain, and will use the image specified in the pod spec. |
| `additionalInitArgs` _string array_ | Additional arguments to pass to the chain init command. |
| `additionalStartArgs` _string array_ | Additional arguments to pass to the chain start command. |


#### ChainVersion





_Appears in:_
- [ChainSpec](#chainspec)

| Field | Description |
| --- | --- |
| `height` _integer_ | The block height when this version should be applied. |
| `image` _string_ | The docker image for this version in "repository:tag" format. E.g. busybox:latest. |
| `setHaltHeight` _boolean_ | Determines if the node should forcefully halt at the upgrade height. |


#### CometBFTConfig



CometBFTConfig configures the config.toml.

_Appears in:_
- [ChainSpec](#chainspec)

| Field | Description |
| --- | --- |
| `rpc` _[RPC](#rpc)_ | RPC configuration for your config.toml |
| `p2p` _[P2P](#p2p)_ | P2P configuration for your config.toml |
| `consensus` _[Consensus](#consensus)_ | Consensus configuration for your config.toml |
| `storage` _[Storage](#storage)_ | Storage configuration for your config.toml |
| `txIndex` _[TxIndex](#txindex)_ | TxIndex configuration for your config.toml |
| `instrumentation` _[Instrumentation](#instrumentation)_ | Instrumentation configuration for your config.toml |
| `statesync` _[Statesync](#statesync)_ | Statesync configuration for your config.toml |
| `tomlOverrides` _string_ |  |


#### Consensus





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `doubleSignCheckHeight` _[uint64](#uint64)_ | If not set, defaults to 0 |
| `skipTimeoutCommit` _[bool](#bool)_ | If not set, defaults to false |
| `createEmptyBlocks` _[bool](#bool)_ | If not set, defaults to true |
| `createEmptyBlocksInterval` _string_ | If not set, defaults to 0s |
| `peerGossipSleepDuration` _string_ | If not set, defaults to 100ms |


#### CosmosFullNode



CosmosFullNode is the Schema for the cosmosfullnodes API

_Appears in:_
- [CosmosFullNodeList](#cosmosfullnodelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cosmos.bharvest/v1`
| `kind` _string_ | `CosmosFullNode`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br /><br />Servers may infer this from the endpoint the client submits requests to.<br /><br />Cannot be updated.<br /><br />In CamelCase.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br /><br />Servers should convert recognized schemas to the latest internal value, and<br /><br />may reject unrecognized values.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[FullNodeSpec](#fullnodespec)_ |  |
| `status` _[FullNodeStatus](#fullnodestatus)_ |  |


#### CosmosFullNodeList



CosmosFullNodeList contains a list of CosmosFullNode



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cosmos.bharvest/v1`
| `kind` _string_ | `CosmosFullNodeList`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br /><br />Servers may infer this from the endpoint the client submits requests to.<br /><br />Cannot be updated.<br /><br />In CamelCase.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br /><br />Servers should convert recognized schemas to the latest internal value, and<br /><br />may reject unrecognized values.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[CosmosFullNode](#cosmosfullnode) array_ |  |


#### DisableStrategy

_Underlying type:_ _string_



_Appears in:_
- [InstanceOverridesSpec](#instanceoverridesspec)



#### FullNodePhase

_Underlying type:_ _string_



_Appears in:_
- [FullNodeStatus](#fullnodestatus)



#### FullNodeProbeStrategy

_Underlying type:_ _string_



_Appears in:_
- [FullNodeProbesSpec](#fullnodeprobesspec)



#### FullNodeProbesSpec



FullNodeProbesSpec configures probes for created pods

_Appears in:_
- [PodSpec](#podspec)

| Field | Description |
| --- | --- |
| `strategy` _[FullNodeProbeStrategy](#fullnodeprobestrategy)_ | Strategy controls the default probes added by the controller.<br /><br />None = Do not add any probes. May be necessary for Sentries using a remote signer. |


#### FullNodeSnapshotStatus





_Appears in:_
- [FullNodeStatus](#fullnodestatus)

| Field | Description |
| --- | --- |
| `podCandidate` _string_ | Which pod name to temporarily delete. Indicates a ScheduledVolumeSnapshot is taking place. For optimal data<br /><br />integrity, pod is temporarily removed so PVC does not have any processes writing to it. |


#### FullNodeSpec



FullNodeSpec defines the desired state of CosmosFullNode

_Appears in:_
- [CosmosFullNode](#cosmosfullnode)

| Field | Description |
| --- | --- |
| `replicas` _integer_ | Number of replicas to create.<br /><br />Individual replicas have a consistent identity. |
| `type` _[FullNodeType](#fullnodetype)_ | Different flavors of the fullnode's configuration.<br /><br />'Sentry' configures the fullnode as a validator sentry, requiring a remote signer such as Horcrux or TMKMS.<br /><br />The remote signer is out of scope for the operator and must be deployed separately. Each pod exposes a privval port<br /><br />for use with the remote signer.<br /><br />If not set, configures node for RPC. |
| `chain` _[ChainSpec](#chainspec)_ | Blockchain-specific configuration. |
| `podTemplate` _[PodSpec](#podspec)_ | Template applied to all pods.<br /><br />Creates 1 pod per replica. |
| `strategy` _[RolloutStrategy](#rolloutstrategy)_ | How to scale pods when performing an update. |
| `volumeClaimTemplate` _[PersistentVolumeClaimSpec](#persistentvolumeclaimspec)_ | Will be used to create a stand-alone PVC to provision the volume.<br /><br />One PVC per replica mapped and mounted to a corresponding pod. |
| `volumeRetentionPolicy` _[RetentionPolicy](#retentionpolicy)_ | Determines how to handle PVCs when pods are scaled down.<br /><br />One of 'Retain' or 'Delete'.<br /><br />If 'Delete', PVCs are deleted if pods are scaled down.<br /><br />If 'Retain', PVCs are not deleted. The admin must delete manually or are deleted if the CRD is deleted.<br /><br />If not set, defaults to 'Delete'. |
| `service` _[ServiceSpec](#servicespec)_ | Configure Operator created services. A singe rpc service is created for load balancing api, grpc, rpc, etc. requests.<br /><br />This allows a k8s admin to use the service in an Ingress, for example.<br /><br />Additionally, multiple p2p services are created for CometBFT peer exchange. |
| `instanceOverrides` _object (keys:string, values:[InstanceOverridesSpec](#instanceoverridesspec))_ | Allows overriding an instance on a case-by-case basis. An instance is a pod/pvc combo with an ordinal.<br /><br />Key must be the name of the pod including the ordinal suffix.<br /><br />Example: cosmos-1<br /><br />Used for debugging. |
| `selfHeal` _[SelfHealSpec](#selfhealspec)_ | Strategies for automatic recovery of faults and errors.<br /><br />Managed by a separate controller, SelfHealingController, in an effort to reduce<br /><br />complexity of the CosmosFullNodeController. |


#### FullNodeStatus



FullNodeStatus defines the observed state of CosmosFullNode

_Appears in:_
- [CosmosFullNode](#cosmosfullnode)

| Field | Description |
| --- | --- |
| `observedGeneration` _integer_ | The most recent generation observed by the controller. |
| `phase` _[FullNodePhase](#fullnodephase)_ | The current phase of the fullnode deployment.<br /><br />"Progressing" means the deployment is under way.<br /><br />"Complete" means the deployment is complete and reconciliation is finished.<br /><br />"WaitingForP2PServices" means the deployment is complete but the p2p services are not yet ready.<br /><br />"Error" means an unrecoverable error occurred, which needs human intervention. |
| `status` _string_ | A generic message for the user. May contain errors. |
| `scheduledSnapshotStatus` _object (keys:string, values:[FullNodeSnapshotStatus](#fullnodesnapshotstatus))_ | Set by the ScheduledVolumeSnapshotController. Used to signal the CosmosFullNode to modify its<br /><br />resources during VolumeSnapshot creation.<br /><br />Map key is the source ScheduledVolumeSnapshot CRD that created the status. |
| `selfHealing` _[SelfHealingStatus](#selfhealingstatus)_ | Status set by the SelfHealing controller. |
| `peers` _string array_ | Persistent peer addresses. |
| `sync` _object (keys:string, values:[SyncInfoPodStatus](#syncinfopodstatus))_ | Current sync information. Collected every 60s. |
| `height` _object (keys:string, values:integer)_ | Latest Height information. collected when node starts up and when RPC is successfully queried. |


#### FullNodeType

_Underlying type:_ _string_



_Appears in:_
- [FullNodeSpec](#fullnodespec)



#### HeightDriftMitigationSpec





_Appears in:_
- [SelfHealSpec](#selfhealspec)

| Field | Description |
| --- | --- |
| `threshold` _integer_ | If pod's height falls behind the max height of all pods by this value or more AND the pod's RPC /status endpoint<br /><br />reports itself as in-sync, the pod is deleted. The CosmosFullNodeController creates a new pod to replace it.<br /><br />Pod deletion respects the CosmosFullNode.Spec.RolloutStrategy and will not delete more pods than set<br /><br />by the strategy to prevent downtime.<br /><br />This workaround is necessary to mitigate a bug in the Cosmos SDK and/or CometBFT where pods report themselves as<br /><br />in-sync even though they can lag thousands of blocks behind the chain tip and cannot catch up.<br /><br />A "rebooted" pod /status reports itself correctly and allows it to catch up to chain tip. |


#### InstanceOverridesSpec



InstanceOverridesSpec allows overriding an instance which is pod/pvc combo with an ordinal

_Appears in:_
- [FullNodeSpec](#fullnodespec)

| Field | Description |
| --- | --- |
| `disable` _[DisableStrategy](#disablestrategy)_ | Disables whole or part of the instance.<br /><br />Used for scenarios like debugging or deleting the PVC and restoring from a dataSource.<br /><br />Set to "Pod" to prevent controller from creating a pod for this instance, leaving the PVC.<br /><br />Set to "All" to prevent the controller from managing a pod and pvc. Note, the PVC may not be deleted if<br /><br />the RetainStrategy is set to "Retain". If you need to remove the PVC, delete manually. |
| `volumeClaimTemplate` _[PersistentVolumeClaimSpec](#persistentvolumeclaimspec)_ | Overrides an individual instance's PVC. |
| `image` _string_ | Overrides an individual instance's Image. |
| `externalAddress` _string_ | Sets an individual instance's external address. |


#### Instrumentation





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `prometheus` _[bool](#bool)_ | Whether you open prometheus service.<br /><br />If not set, defaults to true |
| `prometheusListenAddr` _string_ | Where you want to open prometheus.<br /><br />If not set, defaults to "26660" |


#### Metadata



Metadata is a subset of k8s object metadata.

_Appears in:_
- [PersistentVolumeClaimSpec](#persistentvolumeclaimspec)
- [PodSpec](#podspec)
- [ServiceOverridesSpec](#serviceoverridesspec)

| Field | Description |
| --- | --- |
| `labels` _object (keys:string, values:string)_ | Labels are added to a resource. If there is a collision between labels the Operator creates, the Operator<br /><br />labels take precedence. |
| `annotations` _object (keys:string, values:string)_ | Annotations are added to a resource. If there is a collision between annotations the Operator creates, the Operator<br /><br />annotations take precedence. |


#### NamadaConfig





_Appears in:_
- [ChainSpec](#chainspec)

| Field | Description |
| --- | --- |
| `wasmDir` _string_ | namada use wasm. you can specify dir for wasm.<br /><br />If not set, defaults to "wasm" |
| `ledger` _[NamadaLedger](#namadaledger)_ |  |




#### NamadaLedger





_Appears in:_
- [NamadaConfig](#namadaconfig)

| Field | Description |
| --- | --- |
| `shell` _[NamadaShell](#namadashell)_ |  |
| `ethereumBridge` _[NamadaEthereumBridge](#namadaethereumbridge)_ |  |




#### P2P





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `laddr` _string_ | Listening address for P2P cononection.<br /><br />If not set, defaults to "tcp://127.0.0.1:26656" |
| `externalAddress` _string_ | ExternalAddress using P2P connection.<br /><br />If not set, defaults to "tcp://0.0.0.0:26656" also other peer cannot find to you using PEX. |
| `seeds` _string_ | Seeds for P2P.<br /><br />Comma delimited list of p2p seed nodes in <ID>@<IP>:<PORT> format. |
| `persistentPeers` _string_ | PersistentPeer address list for your P2P connection.<br /><br />Comma delimited list of p2p nodes in <ID>@<IP>:<PORT> format to keep persistent p2p connections. |
| `maxNumInboundPeers` _integer_ | It could be different depending on what chain you run.<br /><br />Cosmos - 20, Namada - 40 |
| `maxNumOutboundPeers` _integer_ | It could be different depending on what chain you run.<br /><br />Cosmos - 20, Namada - 10 |
| `pex` _[bool](#bool)_ | Whether peers can be exchanged.<br /><br />If not set, defaults to true |
| `seedMode` _[bool](#bool)_ | Whether you'll run seed node.<br /><br />WARNING: If you run seed node, the node will disconnect with other peers after transfer your peers.<br /><br />If not set, defaults to false |
| `privatePeerIds` _string_ | For sentry node.<br /><br />Comma delimited list of node/peer IDs to keep private (will not be gossiped to other peers)<br /><br />If not set, defaults to "" |
| `unconditionalPeerIDs` _string_ | Comma delimited list of node/peer IDs, to which a connection will be (re)established ignoring any existing limits. |


#### PVCAutoScaleSpec





_Appears in:_
- [SelfHealSpec](#selfhealspec)

| Field | Description |
| --- | --- |
| `usedSpacePercentage` _integer_ | The percentage of used disk space required to trigger scaling.<br /><br />Example, if set to 80, autoscaling will not trigger until used space reaches >=80% of capacity. |
| `increaseQuantity` _string_ | How much to increase the PVC's capacity.<br /><br />Either a percentage (e.g. 20%) or a resource storage quantity (e.g. 100Gi).<br /><br /><br /><br /><br /><br />If a percentage, the existing capacity increases by the percentage.<br /><br />E.g. PVC of 100Gi capacity + IncreaseQuantity of 20% increases disk to 120Gi.<br /><br /><br /><br /><br /><br />If a storage quantity (e.g. 100Gi), increases by that amount. |
| `maxSize` _[Quantity](#quantity)_ | A resource storage quantity (e.g. 2000Gi).<br /><br />When increasing PVC capacity reaches >= MaxSize, autoscaling ceases.<br /><br />Safeguards against storage quotas and costs. |


#### PVCAutoScaleStatus





_Appears in:_
- [SelfHealingStatus](#selfhealingstatus)

| Field | Description |
| --- | --- |
| `requestedSize` _[Quantity](#quantity)_ | The PVC size requested by the SelfHealing controller. |
| `requestedAt` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#time-v1-meta)_ | The timestamp the SelfHealing controller requested a PVC increase. |


#### PersistentVolumeClaimSpec



PersistentVolumeClaimSpec describes the common attributes of storage devices
and allows a Source for provider-specific attributes

_Appears in:_
- [FullNodeSpec](#fullnodespec)
- [InstanceOverridesSpec](#instanceoverridesspec)

| Field | Description |
| --- | --- |
| `metadata` _[Metadata](#metadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `storageClassName` _string_ | storageClassName is the name of the StorageClass required by the claim.<br /><br />For proper pod scheduling, it's highly recommended to set "volumeBindingMode: WaitForFirstConsumer" in the StorageClass.<br /><br />More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1<br /><br />For GKE, recommended storage class is "premium-rwo".<br /><br />This field is immutable. Updating this field requires manually deleting the PVC.<br /><br />This field is required. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | resources represents the minimum resources the volume should have.<br /><br />If RecoverVolumeExpansionFailure feature is enabled users are allowed to specify resource requirements<br /><br />that are lower than previous value but must still be higher than capacity recorded in the<br /><br />status field of the claim.<br /><br />More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources<br /><br />Updating the storage size is allowed but the StorageClass must support file system resizing.<br /><br />Only increasing storage is permitted.<br /><br />This field is required. |
| `accessModes` _[PersistentVolumeAccessMode](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#persistentvolumeaccessmode-v1-core) array_ | accessModes contain the desired access modes the volume should have.<br /><br />More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1<br /><br />If not specified, defaults to ReadWriteOnce.<br /><br />This field is immutable. Updating this field requires manually deleting the PVC. |
| `volumeMode` _[PersistentVolumeMode](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#persistentvolumemode-v1-core)_ | volumeMode defines what type of volume is required by the claim.<br /><br />Value of Filesystem is implied when not included in claim spec.<br /><br />This field is immutable. Updating this field requires manually deleting the PVC. |
| `dataSource` _[TypedLocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#typedlocalobjectreference-v1-core)_ | Can be used to specify either:<br /><br />* An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot)<br /><br />* An existing PVC (PersistentVolumeClaim)<br /><br />If the provisioner or an external controller can support the specified data source,<br /><br />it will create a new volume based on the contents of the specified data source.<br /><br />If the AnyVolumeDataSource feature gate is enabled, this field will always have<br /><br />the same contents as the DataSourceRef field.<br /><br />If you choose an existing PVC, the PVC must be in the same availability zone. |
| `autoDataSource` _[AutoDataSource](#autodatasource)_ | If set, discovers and dynamically sets dataSource for the PVC on creation.<br /><br />No effect if dataSource field set; that field takes precedence.<br /><br />Configuring autoDataSource may help boostrap new replicas more quickly. |


#### PodSpec





_Appears in:_
- [FullNodeSpec](#fullnodespec)

| Field | Description |
| --- | --- |
| `metadata` _[Metadata](#metadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `envs` _object array_ |  |
| `image` _string_ | Image is the docker reference in "repository:tag" format. E.g. busybox:latest.<br /><br />This is for the main container running the chain process.<br /><br />Note: for granular control over which images are applied at certain block heights,<br /><br />use spec.chain.versions instead. |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#pullpolicy-v1-core)_ | Image pull policy.<br /><br />One of Always, Never, IfNotPresent.<br /><br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br /><br />Cannot be updated.<br /><br />More info: https://kubernetes.io/docs/concepts/containers/images#updating-images<br /><br />This is for the main container running the chain process. |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#localobjectreference-v1-core) array_ | ImagePullSecrets is a list of references to secrets in the same namespace to use for pulling any images<br /><br />in pods that reference this ServiceAccount. ImagePullSecrets are distinct from Secrets because Secrets<br /><br />can be mounted in the pod, but ImagePullSecrets are only accessed by the kubelet.<br /><br />More info: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod<br /><br />This is for the main container running the chain process. |
| `nodeSelector` _object (keys:string, values:string)_ | NodeSelector is a selector which must be true for the pod to fit on a node.<br /><br />Selector which must match a node's labels for the pod to be scheduled on that node.<br /><br />More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/<br /><br />This is an advanced configuration option. |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#affinity-v1-core)_ | If specified, the pod's scheduling constraints<br /><br />This is an advanced configuration option. |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#toleration-v1-core) array_ | If specified, the pod's tolerations.<br /><br />This is an advanced configuration option. |
| `priorityClassName` _string_ | If specified, indicates the pod's priority. "system-node-critical" and<br /><br />"system-cluster-critical" are two special keywords which indicate the<br /><br />highest priorities with the former being the highest priority. Any other<br /><br />name must be defined by creating a PriorityClass object with that name.<br /><br />If not specified, the pod priority will be default or zero if there is no<br /><br />default.<br /><br />This is an advanced configuration option. |
| `priority` _integer_ | The priority value. Various system components use this field to find the<br /><br />priority of the pod. When Priority Admission Controller is enabled, it<br /><br />prevents users from setting this field. The admission controller populates<br /><br />this field from PriorityClassName.<br /><br />The higher the value, the higher the priority.<br /><br />This is an advanced configuration option. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#resourcerequirements-v1-core)_ | Resources describes the compute resource requirements. |
| `terminationGracePeriodSeconds` _integer_ | Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request.<br /><br />Value must be non-negative integer. The value zero indicates stop immediately via<br /><br />the kill signal (no opportunity to shut down).<br /><br />If this value is nil, the default grace period will be used instead.<br /><br />The grace period is the duration in seconds after the processes running in the pod are sent<br /><br />a termination signal and the time when the processes are forcibly halted with a kill signal.<br /><br />Set this value longer than the expected cleanup time for your process.<br /><br />This is an advanced configuration option.<br /><br />Defaults to 30 seconds. |
| `probes` _[FullNodeProbesSpec](#fullnodeprobesspec)_ | Configure probes for the pods managed by the controller. |
| `volumes` _[Volume](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#volume-v1-core) array_ | List of volumes that can be mounted by containers belonging to the pod.<br /><br />More info: https://kubernetes.io/docs/concepts/storage/volumes<br /><br />A strategic merge patch is applied to the default volumes created by the controller.<br /><br />Take extreme caution when using this feature. Use only for critical bugs.<br /><br />Some chains do not follow conventions or best practices, so this serves as an "escape hatch" for the user<br /><br />at the cost of maintainability. |
| `initContainers` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#container-v1-core) array_ | List of initialization containers belonging to the pod.<br /><br />More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/<br /><br />A strategic merge patch is applied to the default init containers created by the controller.<br /><br />Take extreme caution when using this feature. Use only for critical bugs.<br /><br />Some chains do not follow conventions or best practices, so this serves as an "escape hatch" for the user<br /><br />at the cost of maintainability. |
| `containers` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#container-v1-core) array_ | List of containers belonging to the pod.<br /><br />A strategic merge patch is applied to the default containers created by the controller.<br /><br />Take extreme caution when using this feature. Use only for critical bugs.<br /><br />Some chains do not follow conventions or best practices, so this serves as an "escape hatch" for the user<br /><br />at the cost of maintainability. |


#### Pruning



Pruning controls the pruning settings.

_Appears in:_
- [SDKAppConfig](#sdkappconfig)

| Field | Description |
| --- | --- |
| `strategy` _[PruningStrategy](#pruningstrategy)_ | One of default\|nothing\|everything\|custom.<br /><br />default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals.<br /><br />nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node).<br /><br />everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals.<br /><br />custom: allow pruning options to be manually specified through Interval, KeepEvery, KeepRecent. |
| `interval` _[uint32](#uint32)_ | Bock height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom').<br /><br />If not set, defaults to 0. |
| `keepEvery` _[uint32](#uint32)_ | Offset heights to keep on disk after 'keep-every' (ignored if pruning is not 'custom')<br /><br />Often, setting this to 0 is appropriate.<br /><br />If not set, defaults to 0. |
| `keepRecent` _[uint32](#uint32)_ | Number of recent block heights to keep on disk (ignored if pruning is not 'custom')<br /><br />If not set, defaults to 0. |
| `minRetainBlocks` _[uint32](#uint32)_ | Defines the minimum block height offset from the current<br /><br />block being committed, such that all blocks past this offset are pruned<br /><br />from CometBFT. It is used as part of the process of determining the<br /><br />ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates<br /><br />that no blocks should be pruned.<br /><br /><br /><br /><br /><br />This configuration value is only responsible for pruning Comet blocks.<br /><br />It has no bearing on application state pruning which is determined by the<br /><br />"pruning-*" configurations.<br /><br /><br /><br /><br /><br />Note: CometBFT block pruning is dependent on this parameter in conjunction<br /><br />with the unbonding (safety threshold) period, state pruning and state sync<br /><br />snapshot parameters to determine the correct minimum value of<br /><br />ResponseCommit.RetainHeight.<br /><br /><br /><br /><br /><br />If not set, defaults to 0. |


#### PruningStrategy

_Underlying type:_ _string_

PruningStrategy control pruning.

_Appears in:_
- [Pruning](#pruning)



#### RPC





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `laddr` _string_ | Listening address for RPC.<br /><br />If not set, defaults to "tcp://0.0.0.0:26657" |
| `corsAllowedOrigins` _[string](#string)_ | rpc list of origins a cross-domain request can be executed from.<br /><br />Default value '[]' disables cors support.<br /><br />Use '["*"]' to allow any origin. |
| `corsAllowedMethods` _[string](#string)_ | If not set, defaults to "["HEAD", "GET", "POST"]" |
| `timeoutBroadcastTxCommit` _string_ | timeout for broadcast_tx_commit<br /><br />If not set, defaults to "10000ms"(also "10s") |


#### RetentionPolicy

_Underlying type:_ _string_



_Appears in:_
- [FullNodeSpec](#fullnodespec)



#### RolloutStrategy



RolloutStrategy is an update strategy that can be shared between several Cosmos CRDs.

_Appears in:_
- [FullNodeSpec](#fullnodespec)

| Field | Description |
| --- | --- |
| `maxUnavailable` _[IntOrString](#intorstring)_ | The maximum number of pods that can be unavailable during an update.<br /><br />Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).<br /><br />Absolute number is calculated from percentage by rounding down. The minimum max unavailable is 1.<br /><br />Defaults to 25%.<br /><br />Example: when this is set to 30%, pods are scaled down to 70% of desired pods<br /><br />immediately when the rolling update starts. Once new pods are ready, pods<br /><br />can be scaled down further, ensuring that the total number of pods available<br /><br />at all times during the update is at least 70% of desired pods. |


#### SDKAppConfig



SDKAppConfig configures the cosmos sdk application app.toml.

_Appears in:_
- [ChainSpec](#chainspec)

| Field | Description |
| --- | --- |
| `skipInvariants` _boolean_ | Skip x/crisis invariants check on startup. |
| `snapshotURL` _string_ | URL for a snapshot archive to download from the internet.<br /><br />Unarchiving the snapshot populates the data directory.<br /><br />Although this field is optional, you will almost always want to set it.<br /><br />The operator detects and properly handles the following file extensions:<br /><br />.tar, .tar.gz, .tar.gzip, .tar.lz4<br /><br />Use SnapshotScript if the snapshot archive is unconventional or requires special handling. |
| `snapshotScript` _string_ | Specify shell (sh) script commands to properly download and process a snapshot archive.<br /><br />Prefer SnapshotURL if possible.<br /><br />The available shell commands are from docker image ghcr.io/strangelove-ventures/infra-toolkit, including wget and curl.<br /><br />Save the file to env var $GENESIS_FILE.<br /><br />Takes precedence over SnapshotURL.<br /><br />Hint: Use "set -eux" in your script.<br /><br />Available env vars:<br /><br />$HOME: The user's home directory.<br /><br />$CHAIN_HOME: The home directory for the chain, aka: --home flag<br /><br />$DATA_DIR: The directory for the database files. |
| `minGasPrice` _string_ | The minimum gas prices a validator is willing to accept for processing a<br /><br />transaction. A transaction's fees must meet the minimum of any denomination<br /><br />specified in this config (e.g. 0.25token1;0.0001token2). |
| `apiEnableUnsafeCORS` _boolean_ | Defines if CORS should be enabled for the API (unsafe - use it at your own risk). |
| `grpcWebEnableUnsafeCORS` _boolean_ | Defines if CORS should be enabled for grpc-web (unsafe - use it at your own risk). |
| `pruning` _[Pruning](#pruning)_ | Controls pruning settings. i.e. How much data to keep on disk.<br /><br />If not set, defaults to "default" pruning strategy. |
| `haltHeight` _[uint64](#uint64)_ | If set, block height at which to gracefully halt the chain and shutdown the node.<br /><br />Useful for testing or upgrades. |
| `tomlOverrides` _string_ | Custom app config toml.<br /><br />Values entered here take precedence over all other configuration.<br /><br />Must be valid toml.<br /><br />Important: all keys must be "kebab-case" which differs from config.toml. |


#### SelfHealSpec



SelfHealSpec is part of a CosmosFullNode but is managed by a separate controller, SelfHealingReconciler.
This is an effort to reduce complexity in the CosmosFullNodeReconciler.
The controller only modifies the CosmosFullNode's status subresource relying on the CosmosFullNodeReconciler
to reconcile appropriately.

_Appears in:_
- [FullNodeSpec](#fullnodespec)

| Field | Description |
| --- | --- |
| `pvcAutoScale` _[PVCAutoScaleSpec](#pvcautoscalespec)_ | Automatically increases PVC storage as they approach capacity.<br /><br /><br /><br /><br /><br />Your cluster must support and use the ExpandInUsePersistentVolumes feature gate. This allows volumes to<br /><br />expand while a pod is attached to it, thus eliminating the need to restart pods.<br /><br />If you cluster does not support ExpandInUsePersistentVolumes, you will need to manually restart pods after<br /><br />resizing is complete. |
| `heightDriftMitigation` _[HeightDriftMitigationSpec](#heightdriftmitigationspec)_ | Take action when a pod's height falls behind the max height of all pods AND still reports itself as in-sync. |


#### SelfHealingStatus





_Appears in:_
- [FullNodeStatus](#fullnodestatus)

| Field | Description |
| --- | --- |
| `pvcHealer` _object (keys:string, values:[PVCAutoScaleStatus](#pvcautoscalestatus))_ | PVC auto-scaling status. |


#### ServiceOverridesSpec



ServiceOverridesSpec allows some overrides for the created, single RPC service.

_Appears in:_
- [ServiceSpec](#servicespec)

| Field | Description |
| --- | --- |
| `metadata` _[Metadata](#metadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `type` _[ServiceType](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#servicetype-v1-core)_ | Describes ingress methods for a service.<br /><br />If not set, defaults to "ClusterIP". |
| `externalTrafficPolicy` _[ServiceExternalTrafficPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#serviceexternaltrafficpolicy-v1-core)_ | Sets endpoint and routing behavior.<br /><br />See: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#caveats-and-limitations-when-preserving-source-ips<br /><br />If not set, defaults to "Cluster". |
| `ports` _[ServicePort](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#serviceport-v1-core) array_ |  |


#### ServiceSpec





_Appears in:_
- [FullNodeSpec](#fullnodespec)

| Field | Description |
| --- | --- |
| `maxP2PExternalAddresses` _integer_ | Max number of external p2p services to create for CometBFT peer exchange.<br /><br />The public endpoint is set as the "p2p.external_address" in the config.toml.<br /><br />Controller creates p2p services for each pod so that every pod can peer with each other internally in the cluster.<br /><br />This setting allows you to control the number of p2p services exposed for peers outside of the cluster to use.<br /><br />If not set, defaults to 1. |
| `p2pTemplate` _[ServiceOverridesSpec](#serviceoverridesspec)_ | Overrides for all P2P services that need external addresses. |
| `rpcTemplate` _[ServiceOverridesSpec](#serviceoverridesspec)_ | Overrides for the single RPC service. |


#### Statesync





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `enable` _[bool](#bool)_ | which you enable stateSync<br /><br />If not set, defaults to false |
| `rpcServers` _string_ |  |
| `trustHeight` _[uint64](#uint64)_ |  |
| `trustHash` _string_ |  |
| `trustPeriod` _string_ | If not set, defaults to "168h0m0s" |
| `discoveryTime` _string_ | If not set, defaults to "15000ms"("15s") |
| `tempDir` _string_ |  |


#### Storage





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `discardAbciResponses` _[bool](#bool)_ | Set to true to discard ABCI responses from the state store, which can save a<br /><br />considerable amount of disk space. Set to false to ensure ABCI responses are<br /><br />persisted. ABCI responses are required for /block_results RPC queries, and to<br /><br />reindex events in the command-line tool.<br /><br /><br /><br /><br /><br />If not set, defaults to false |


#### SyncInfoPodStatus





_Appears in:_
- [FullNodeStatus](#fullnodestatus)

| Field | Description |
| --- | --- |
| `timestamp` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#time-v1-meta)_ | When consensus information was fetched. |
| `height` _[uint64](#uint64)_ | Latest height if no error encountered. |
| `inSync` _[bool](#bool)_ | If the pod reports itself as in sync with chain tip. |
| `error` _string_ | Error message if unable to fetch consensus state. |


#### TxIndex





_Appears in:_
- [CometBFTConfig](#cometconfig)

| Field | Description |
| --- | --- |
| `indexer` _string_ | It could be different depending on what chain you run.<br /><br />cosmos - "kv", namada - "null" |



## cosmos.bharvest/v1alpha1

Package v1alpha1 contains API Schema definitions for the cosmos v1alpha1 API group

### Resource Types
- [ScheduledVolumeSnapshot](#scheduledvolumesnapshot)
- [ScheduledVolumeSnapshotList](#scheduledvolumesnapshotlist)
- [StatefulJob](#statefuljob)
- [StatefulJobList](#statefuljoblist)



#### JobTemplateSpec



JobTemplateSpec is a subset of batchv1.JobSpec.

_Appears in:_
- [StatefulJobSpec](#statefuljobspec)

| Field | Description |
| --- | --- |
| `activeDeadlineSeconds` _integer_ | Specifies the duration in seconds relative to the startTime that the job<br /><br />may be continuously active before the system tries to terminate it; value<br /><br />must be positive integer.<br /><br />Do not set too short or you will run into PVC/VolumeSnapshot provisioning rate limits.<br /><br />Defaults to 24 hours. |
| `backoffLimit` _integer_ | Specifies the number of retries before marking this job failed.<br /><br />Defaults to 5. |
| `ttlSecondsAfterFinished` _integer_ | Limits the lifetime of a Job that has finished<br /><br />execution (either Complete or Failed). If this field is set,<br /><br />ttlSecondsAfterFinished after the Job finishes, it is eligible to be<br /><br />automatically deleted. When the Job is being deleted, its lifecycle<br /><br />guarantees (e.g. finalizers) will be honored. If this field is set to zero,<br /><br />the Job becomes eligible to be deleted immediately after it finishes.<br /><br />Defaults to 15 minutes to allow some time to inspect logs. |


#### LocalFullNodeRef





_Appears in:_
- [ScheduledVolumeSnapshotSpec](#scheduledvolumesnapshotspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the object, metadata.name |
| `namespace` _string_ | DEPRECATED: CosmosFullNode must be in the same namespace as the ScheduledVolumeSnapshot. This field is ignored. |
| `ordinal` _integer_ | Index of the pod to snapshot. If not provided, will do any pod in the CosmosFullNode.<br /><br />Useful when snapshots are local to the same node as the pod, requiring snapshots across multiple pods/nodes. |


#### ScheduledVolumeSnapshot



ScheduledVolumeSnapshot is the Schema for the scheduledvolumesnapshots API

_Appears in:_
- [ScheduledVolumeSnapshotList](#scheduledvolumesnapshotlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cosmos.bharvest/v1alpha1`
| `kind` _string_ | `ScheduledVolumeSnapshot`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br /><br />Servers may infer this from the endpoint the client submits requests to.<br /><br />Cannot be updated.<br /><br />In CamelCase.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br /><br />Servers should convert recognized schemas to the latest internal value, and<br /><br />may reject unrecognized values.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ScheduledVolumeSnapshotSpec](#scheduledvolumesnapshotspec)_ |  |
| `status` _[ScheduledVolumeSnapshotStatus](#scheduledvolumesnapshotstatus)_ |  |


#### ScheduledVolumeSnapshotList



ScheduledVolumeSnapshotList contains a list of ScheduledVolumeSnapshot



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cosmos.bharvest/v1alpha1`
| `kind` _string_ | `ScheduledVolumeSnapshotList`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br /><br />Servers may infer this from the endpoint the client submits requests to.<br /><br />Cannot be updated.<br /><br />In CamelCase.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br /><br />Servers should convert recognized schemas to the latest internal value, and<br /><br />may reject unrecognized values.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[ScheduledVolumeSnapshot](#scheduledvolumesnapshot) array_ |  |


#### ScheduledVolumeSnapshotSpec



ScheduledVolumeSnapshotSpec defines the desired state of ScheduledVolumeSnapshot
Creates recurring VolumeSnapshots of a PVC managed by a CosmosFullNode.
A VolumeSnapshot is a CRD (installed in GKE by default).
See: https://kubernetes.io/docs/concepts/storage/volume-snapshots/
This enables recurring, consistent backups.
To prevent data corruption, a pod is temporarily deleted while the snapshot takes place which could take
several minutes.
Therefore, if you create a ScheduledVolumeSnapshot, you must use replica count >= 2 to prevent downtime.
If <= 1 pod in a ready state, the controller will not temporarily delete the pod. The controller makes every
effort to prevent downtime.
Only 1 VolumeSnapshot is created at a time, so at most only 1 pod is temporarily deleted.
Multiple, parallel VolumeSnapshots are not supported.

_Appears in:_
- [ScheduledVolumeSnapshot](#scheduledvolumesnapshot)

| Field | Description |
| --- | --- |
| `fullNodeRef` _[LocalFullNodeRef](#localfullnoderef)_ | Reference to the source CosmosFullNode.<br /><br />This field is immutable. If you change the fullnode, you may encounter undefined behavior.<br /><br />The CosmosFullNode must be in the same namespace as the ScheduledVolumeSnapshot.<br /><br />Instead delete the ScheduledVolumeSnapshot and create a new one with the correct fullNodeRef. |
| `schedule` _string_ | A crontab schedule using the standard as described in https://en.wikipedia.org/wiki/Cron.<br /><br />See https://crontab.guru for format.<br /><br />Kubernetes providers rate limit VolumeSnapshot creation. Therefore, setting a crontab that's<br /><br />too frequent may result in rate limiting errors. |
| `volumeSnapshotClassName` _string_ | The name of the VolumeSnapshotClass to use when creating snapshots. |
| `deletePod` _boolean_ | If true, the controller will temporarily delete the candidate pod before taking a snapshot of the pod's associated PVC.<br /><br />This option prevents writes to the PVC, ensuring the highest possible data integrity.<br /><br />Once the snapshot is created, the pod will be restored. |
| `minAvailable` _integer_ | Minimum number of CosmosFullNode pods that must be ready before creating a VolumeSnapshot.<br /><br />In the future, this field will have no effect unless spec.deletePod=true.<br /><br />This controller gracefully deletes a pod while taking a snapshot. Then recreates the pod once the<br /><br />snapshot is complete.<br /><br />This way, the snapshot has the highest possible data integrity.<br /><br />Defaults to 2.<br /><br />Warning: If set to 1, you will experience downtime. |
| `limit` _integer_ | The number of recent VolumeSnapshots to keep.<br /><br />Defaults to 3. |
| `suspend` _boolean_ | If true, the controller will not create any VolumeSnapshots.<br /><br />This allows you to disable creation of VolumeSnapshots without deleting the ScheduledVolumeSnapshot resource.<br /><br />This pattern works better when using tools such as Kustomzie.<br /><br />If a pod is temporarily deleted, it will be restored. |


#### ScheduledVolumeSnapshotStatus



ScheduledVolumeSnapshotStatus defines the observed state of ScheduledVolumeSnapshot

_Appears in:_
- [ScheduledVolumeSnapshot](#scheduledvolumesnapshot)

| Field | Description |
| --- | --- |
| `observedGeneration` _integer_ | The most recent generation observed by the controller. |
| `status` _string_ | A generic message for the user. May contain errors. |
| `phase` _[SnapshotPhase](#snapshotphase)_ | The phase of the controller. |
| `createdAt` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#time-v1-meta)_ | The date when the CRD was created.<br /><br />Used as a reference when calculating the next time to create a snapshot. |
| `candidate` _[SnapshotCandidate](#snapshotcandidate)_ | The pod/pvc pair of the CosmosFullNode from which to make a VolumeSnapshot. |
| `lastSnapshot` _[VolumeSnapshotStatus](#volumesnapshotstatus)_ | The most recent volume snapshot created by the controller. |


#### SnapshotCandidate





_Appears in:_
- [ScheduledVolumeSnapshotStatus](#scheduledvolumesnapshotstatus)

| Field | Description |
| --- | --- |
| `podName` _string_ |  |
| `pvcName` _string_ |  |
| `podLabels` _object (keys:string, values:string)_ |  |


#### SnapshotPhase

_Underlying type:_ _string_



_Appears in:_
- [ScheduledVolumeSnapshotStatus](#scheduledvolumesnapshotstatus)



#### StatefulJob



StatefulJob is the Schema for the statefuljobs API

_Appears in:_
- [StatefulJobList](#statefuljoblist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cosmos.bharvest/v1alpha1`
| `kind` _string_ | `StatefulJob`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br /><br />Servers may infer this from the endpoint the client submits requests to.<br /><br />Cannot be updated.<br /><br />In CamelCase.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br /><br />Servers should convert recognized schemas to the latest internal value, and<br /><br />may reject unrecognized values.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[StatefulJobSpec](#statefuljobspec)_ |  |
| `status` _[StatefulJobStatus](#statefuljobstatus)_ |  |


#### StatefulJobList



StatefulJobList contains a list of StatefulJob



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cosmos.bharvest/v1alpha1`
| `kind` _string_ | `StatefulJobList`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br /><br />Servers may infer this from the endpoint the client submits requests to.<br /><br />Cannot be updated.<br /><br />In CamelCase.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br /><br />Servers should convert recognized schemas to the latest internal value, and<br /><br />may reject unrecognized values.<br /><br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[StatefulJob](#statefuljob) array_ |  |


#### StatefulJobSpec



StatefulJobSpec defines the desired state of StatefulJob

_Appears in:_
- [StatefulJob](#statefuljob)

| Field | Description |
| --- | --- |
| `selector` _object (keys:string, values:string)_ | The selector to target VolumeSnapshots. |
| `interval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#duration-v1-meta)_ | Interval at which the controller runs snapshot job with pvc.<br /><br />Expressed as a duration string, e.g. 1.5h, 24h, 12h.<br /><br />Defaults to 24h. |
| `jobTemplate` _[JobTemplateSpec](#jobtemplatespec)_ | Specification of the desired behavior of the job. |
| `podTemplate` _[PodTemplateSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#podtemplatespec-v1-core)_ | Specification of the desired behavior of the job's pod.<br /><br />You should include container commands and args to perform the upload of data to a remote location like an<br /><br />object store such as S3 or GCS.<br /><br />Volumes will be injected and mounted into every container in the spec.<br /><br />Working directory will be /home/operator.<br /><br />The chain directory will be /home/operator/cosmos and set as env var $CHAIN_HOME.<br /><br />If not set, pod's restart policy defaults to Never. |
| `volumeClaimTemplate` _[StatefulJobVolumeClaimTemplate](#statefuljobvolumeclaimtemplate)_ | Specification for the PVC associated with the job. |


#### StatefulJobStatus



StatefulJobStatus defines the observed state of StatefulJob

_Appears in:_
- [StatefulJob](#statefuljob)

| Field | Description |
| --- | --- |
| `observedGeneration` _integer_ | The most recent generation observed by the controller. |
| `status` _string_ | A generic message for the user. May contain errors. |
| `jobHistory` _[JobStatus](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#jobstatus-v1-batch) array_ | Last 5 job statuses created by the controller ordered by more recent jobs. |


#### StatefulJobVolumeClaimTemplate



StatefulJobVolumeClaimTemplate is a subset of a PersistentVolumeClaimTemplate

_Appears in:_
- [StatefulJobSpec](#statefuljobspec)

| Field | Description |
| --- | --- |
| `storageClassName` _string_ | The StorageClass to use when creating a temporary PVC for processing the data.<br /><br />On GKE, the StorageClass must be the same as the PVC's StorageClass from which the<br /><br />VolumeSnapshot was created. |
| `accessModes` _[PersistentVolumeAccessMode](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#persistentvolumeaccessmode-v1-core) array_ | The desired access modes the volume should have.<br /><br />More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1<br /><br />Defaults to ReadWriteOnce. |


#### VolumeSnapshotStatus





_Appears in:_
- [ScheduledVolumeSnapshotStatus](#scheduledvolumesnapshotstatus)

| Field | Description |
| --- | --- |
| `name` _string_ | The name of the created VolumeSnapshot. |
| `startedAt` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v/#time-v1-meta)_ | The time the controller created the VolumeSnapshot. |
| `status` _[VolumeSnapshotStatus](#volumesnapshotstatus)_ | The last VolumeSnapshot's status |


