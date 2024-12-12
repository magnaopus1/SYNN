server:
  address: "mainnet.synnergy.com:8080"  # Mainnet server URL

nodes:
  bootstrap_node:
    ip: "http://192.168.1.101:8080"     # IP and port for Bootstrap Node
    type: "bootstrap"
    ledger_path: "/var/lib/synnergy/ledger_node"
    validator: true                     # Bootstrap Node is a validator
    max_connections: 100                # Max number of connections this node can handle
    connection_pool:
      max_idle_time: 300                # Max idle time for connections (in seconds)
      max_connections: 50               # Max number of active connections in the pool

  new_node_001:
    ip: "http://192.168.1.105:8080"     # IP and port for new Node
    type: "regular"
    ledger_path: "/var/lib/synnergy/ledger_new_node"
    validator: false                    # This is not a validator node
    max_connections: 50                 # Max number of connections this node can handle
    connection_pool:
      max_idle_time: 200                # Max idle time for connections (in seconds)
      max_connections: 20               # Max number of active connections in the pool

# Quorum settings for fault tolerance
quorum:
  min_nodes_required: 2                 # Minimum number of nodes required for quorum
