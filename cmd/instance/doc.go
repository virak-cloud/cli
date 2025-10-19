// Package instance provides commands for managing virtual machine instances.
//
// The commands in this package allow users to create, delete, list, show, start, stop, reboot, and rebuild instances.
// It also provides commands for managing instance snapshots and volumes.
// All commands require a zone ID, which can be provided via the --zoneId flag or by setting a default zone in the configuration.
//
// Examples:
//
//  # Create a new instance
//  virak-cli instance create --name my-instance --vm-image-id <image-id> --service-offering-id <offering-id>
//
//  # List all instances in a zone
//  virak-cli instance list
//
//  # Show details of an instance
//  virak-cli instance show --id <instance-id>
//
//  # Stop an instance
//  virak-cli instance stop --id <instance-id>
//
//  # Start an instance
//  virak-cli instance start --id <instance-id>
//
//  # Delete an instance
//  virak-cli instance delete --id <instance-id>
package instance
