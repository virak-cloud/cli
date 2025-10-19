// Package dns provides commands for managing DNS domains and records.
//
// The commands in this package allow users to create, delete, list, and show DNS domains,
// as well as create, delete, list, and update DNS records.
//
// Examples:
//
//  # Create a new DNS domain
//  virak-cli dns domain create --name example.com
//
//  # List all DNS domains
//  virak-cli dns domain list
//
//  # Create a new DNS record
//  virak-cli dns record create --domain-name example.com --name www --type A --content 192.0.2.1
//
//  # List all DNS records for a domain
//  virak-cli dns record list --domain-name example.com
package dns
