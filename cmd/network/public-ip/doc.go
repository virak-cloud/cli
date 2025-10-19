// Package publicip provides commands for managing public IPs.
//
// The commands in this package allow users to list, associate, and disassociate public IPs.
// It also provides commands for managing static NAT.
//
// Examples:
//
//  # List all public IPs for a network
//  virak-cli network public-ip list --network-id <network-id>
//
//  # Associate a public IP to a network
//  virak-cli network public-ip associate --network-id <network-id>
//
//  # Disassociate a public IP from a network
//  virak-cli network public-ip disassociate --network-id <network-id> --ip-address <ip-address>
package publicip
