# Hello!
# In here we'll need some simple A records to link the domain to the ip and a cloudflare port route ruleset.

# Before doing this you'll want to setup the `google_compute_address` resources in the kubeip file
# so it's ip value can be used dynamically here

# Do not do a root record, but instead use vdl.alexgrand.dev and www.vdl.alexgrand.dev for the domain targets
# We'll want the records to be proxied and the TTL should be short, lets say 10 mins
# https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/resources/record

# The ruleset should apply to any request coming to vdl.alexgrand.dev or www.vdl.alexgrand.dev
# and route the port to out GKE nodeport
# https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/resources/ruleset#nested-schema-for-rules
