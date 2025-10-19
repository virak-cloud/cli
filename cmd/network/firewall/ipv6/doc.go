// Package ipv6 provides commands for managing IPv6 firewall rules.
//
// The commands in this package allow users to create, delete, and list IPv6 firewall rules.
//
// Examples:
//
//  # List all IPv6 firewall rules for a network
//  virak-cli network firewall ipv6 list --network-id <network-id>
//
//  # Create a new IPv6 firewall rule
//  virak-cli network firewall ipv6 create --network-id <network-id> --protocol tcp --cidr 2001:db8::/32
//
//  # Delete an IPv6 firewall rule
//  virak-cli network firewall ipv6 delete --network-id <network-id> --rule-id <rule-id>
package ipv6
