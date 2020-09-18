import FlowIdentityTable from 0xIDENTITYTABLEADDRESS

// This script returns the initial weight of a node
// You must fill in `{EPOCHPHASE}` with 
// `Current`, `Previous`, or `Proposed` with the correct phase

pub fun main(nodeID: String): UInt64 {
    let nodeInfo = FlowIdentityTable.getAll{EPOCHPHASE}NodeInfo()

    let node = nodeInfo[nodeID]
        ?? panic("Node with the specified nodeID does not exist")

    return node.initialWeight
}