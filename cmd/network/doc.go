// Package network provides commands for managing network resources.
//
// The commands in this package allow users to create, delete, list, and show networks.
// It also provides commands for managing network firewall rules, load balancers, public IPs, and VPNs.
// All commands require a zone ID, which can be provided via the --zoneId flag or by setting a default zone in the configuration.
//
// Examples:
//
//  # Create a new L3 network
//  virak-cli network create l3 --name my-network --display-text my-network
//
//  # List all networks in a zone
//  virak-cli network list
//
//  # Show details of a network
//  virak-cli network show --id <network-id>
//
//  # Delete a network
//  virak-cli network delete --id <network-id>
package network
