import FlowIDTableStaking from 0xIDENTITYTABLEADDRESS

// This script returns the balance of unstaked tokens of a node

pub fun main(nodeID: String): UFix64 {
<<<<<<< HEAD
    return FlowIDTableStaking.getNodeUnStakedBalance(nodeID)!
=======
    let nodeInfo = FlowIDTableStaking.NodeInfo(nodeID: nodeID)
    return nodeInfo.tokensUnstaked
>>>>>>> struct
}