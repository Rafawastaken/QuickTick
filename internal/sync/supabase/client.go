package supabase

import (
	"github.com/nedpals/supabase-go"
)

type Client struct {
	client *supabase.Client
	token  string
}

func NewClient(url, key, token string) *Client {
	client := supabase.CreateClient(url, key)
	// If token exists, use it for Auth.
	// But `client.Auth.User(ctx, token)` is to get user.
	// To perform DB operations *as* user, usually libraries have a way to set header.
	// nedpals/supabase-go: Client struct has no direct header manipulation exposed easily?
	// It uses `client.DB` which is likely using rest-go or similar.
	// Wait, `CreateClient` returns `*Client`.
	// Does it expose `Header`? No?

	// Actually supabase-go (nedpals) usually allows setting headers on requests but maybe not globally easily.
	// EXCEPT: If we look at how `postgrest-go` (which it likely wraps) works.
	// But wait, `supabase.CreateClient` is high level.

	// Let's check if there's `AuthToken` or similar.
	// If not, we might be stuck unless we find how to pass it.

	// Alternative: `client.DB.From(...).Auth(token)` ?
	// Let's check `config.go` or documentation if we could.
	// I don't have internet access to docs.
	// But `supabase-go` typically has `Auth` method on the query builder.

	// Let's modify `client` struct to store token and use it in `Sync`.
	// But `Sync` uses `client.client.DB.From(...)`.
	// If `From(...)` returns a builder that has `.Auth(token)`, we should use it there!

	// So `NewClient` just stores it?
	// The `Client` struct I defined in `client.go` only has `client *supabase.Client`.
	// I should add `token string` to `Client` struct and use it in `Sync`!

	return &Client{
		client: client,
		token:  token, // We need to add this field
	}
}
