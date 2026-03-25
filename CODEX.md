# Project: terraform-provider-netskopesdwan

## Goal
Build a Terraform provider in Go for Netskope SD-WAN API v2.

## Constraints
- Use Go
- Use Terraform Plugin Framework
- Keep implementation simple and explicit
- Prefer handwritten client code over code generation initially
- Do not invent API fields; if uncertain, leave a TODO and explain the assumption
- Keep Terraform UX clean and stable
- Small diffs only
- Always run formatting and build steps after changes when possible

## Provider identity
- Terraform provider type name: netskopesdwan
- Local source namespace during development: virtualmo/netskopesdwan

## Initial scope
1. provider configuration
2. shared API client
3. data source: netskopesdwan_gateways
4. data source: netskopesdwan_gateway

## Provider configuration
- base_url
- api_token
- optional timeout
- optional insecure placeholder

## Environment variables
- NETSKOPESDWAN_BASE_URL
- NETSKOPESDWAN_API_TOKEN
- NETSKOPESDWAN_TIMEOUT
- NETSKOPESDWAN_INSECURE

## Engineering rules
- Keep API DTOs separate from Terraform schema models when useful
- Prefer clear error messages
- Add TODOs for unknown API behavior
- Never rename public Terraform schema fields without explicit instruction
- Summarize assumptions in every response

## Working style
- Make one logical change per task
- After code changes, explain:
  1. what changed
  2. why
  3. any assumptions
  4. any next recommended task
