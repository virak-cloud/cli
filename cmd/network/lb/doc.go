// Package lb provides commands for managing network load balancers.
//
// The commands in this package allow users to create, delete, list, assign, and deassign load balancer rules.
//
// Examples:
//
//  # List all load balancer rules for a network
//  virak-cli network lb list --network-id <network-id>
//
//  # Create a new load balancer rule
//  virak-cli network lb create --network-id <network-id> --name my-lb-rule --public-port 80 --private-port 8080
//
//  # Delete a load balancer rule
//  virak-cli network lb delete --network-id <network-id> --rule-id <rule-id>
package lb
