import FlowIDTableStaking from 0xIDENTITYTABLEADDRESS

// This script returns the balance of staked tokens of a delegator

pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
<<<<<<< HEAD
    return FlowIDTableStaking.getDelegatorStakedBalance(nodeID, delegatorID: delegatorID)!
=======
    let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
    return delInfo.tokensStaked
>>>>>>> struct
}