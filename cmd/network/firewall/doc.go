// Package firewall provides commands for managing network firewalls.
//
// The commands in this package allow users to manage IPv4 and IPv6 firewall rules.
//
// Examples:
//
//  # List all IPv4 firewall rules for a network
//  virak-cli network firewall ipv4 list --network-id <network-id>
//
//  # Create a new IPv4 firewall rule
//  virak-cli network firewall ipv4 create --network-id <network-id> --protocol tcp --cidr 192.168.1.0/24
//
//  # Delete an IPv4 firewall rule
//  virak-cli network firewall ipv4 delete --network-id <network-id> --rule-id <rule-id>
package firewall
