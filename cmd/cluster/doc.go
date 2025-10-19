// Package cluster provides commands for managing Kubernetes clusters.
//
// The commands in this package allow users to create, delete, list, show, scale, start, stop, and update Kubernetes clusters.
// All commands require a zone ID, which can be provided via the --zoneId flag or by setting a default zone in the configuration.
//
// Examples:
//
//  # Create a new Kubernetes cluster
//  virak-cli cluster create --name my-cluster --version 1.28 --node-pool-size 3
//
//  # List all Kubernetes clusters in a zone
//  virak-cli cluster list
//
//  # Show details of a Kubernetes cluster
//  virak-cli cluster show --cluster-id <cluster-id>
//
//  # Scale a Kubernetes cluster
//  virak-cli cluster scale --cluster-id <cluster-id> --node-pool-size 5
//
//  # Delete a Kubernetes cluster
//  virak-cli cluster delete --cluster-id <cluster-id>
package cluster
