// Package bucket provides commands for managing object storage buckets.
//
// The commands in this package allow users to create, delete, list, show, and update object storage buckets.
// All commands require a zone ID, which can be provided via the --zoneId flag or by setting a default zone in the configuration.
//
// Examples:
//
//   # Create a new bucket
//   virak-cli bucket create --name my-bucket --policy Private
//
//   # List all buckets in a zone
//   virak-cli bucket list
//
//   # Show details of a bucket
//   virak-cli bucket show --bucketId <bucket-id>
//
//   # Update a bucket's policy
//   virak-cli bucket update --bucketId <bucket-id> --policy Public
//
//   # Delete a bucket
//   virak-cli bucket delete --bucketId <bucket-id>
package bucket
