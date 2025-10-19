// Package staticnat provides commands for managing static NAT for public IPs.
//
// The commands in this package allow users to enable and disable static NAT for public IPs.
//
// Examples:
//
//  # Enable static NAT for a public IP
//  virak-cli network public-ip staticnat enable --network-id <network-id> --public-ip-id <public-ip-id> --instance-id <instance-id>
//
//  # Disable static NAT for a public IP
//  virak-cli network public-ip staticnat disable --network-id <network-id> --public-ip-id <public-ip-id>
package staticnat
