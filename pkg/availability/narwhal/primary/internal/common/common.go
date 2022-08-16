package common

import t "github.com/filecoin-project/mir/pkg/types"

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self   t.ModuleID // id of this module
	Worker t.ModuleID // id of worker modules on worker nodes
	Net    t.ModuleID
	Crypto t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
type ModuleParams struct {
	InstanceUID  []byte     // unique identifier for this instance of BCB, used to prevent cross-instance replay attacks.
	PrimaryNodes []t.NodeID // the list of IDs of all primary nodes in the system.
	F            int        // the maximum number of primary failures tolerated. Must be at most (len(PrimaryNodes)-1) / 3.
	WorkerNodes  []t.NodeID // the list of IDs of worker nodes in this validator.
	NodeID       t.NodeID   // ID of this node
}
