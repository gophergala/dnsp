// Package dnsp implements a simple DNS proxy.
//
// Queries are blocked or resolved based on a blacklist or a whitelist.
// Wildcard host patterns are supported (e.g. "*.com") as well as hosted,
// community-managed hosts files.
package dnsp
