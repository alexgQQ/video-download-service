# Sup!
# Here we'll setup the GKE cluster and a firewall rule to allow traffic

# We'll need a service account for automated GKE actions, no binding or roles needed here tho.

# We'll need a single cluster with two nodes for a node pool. They should all be within `gke_zone` location. Name the cluster "primary-cluster".
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster#example-usage---with-a-separately-managed-node-pool-recommended
# The nodes should use the service account made above and should be named "primary-node" and "kubeip-node" respectivley
# Both nodes require these oauth scopes:
# * https://www.googleapis.com/auth/cloud-platform
# * https://www.googleapis.com/auth/service.management
# * https://www.googleapis.com/auth/servicecontrol
# * https://www.googleapis.com/auth/compute
# * https://www.googleapis.com/auth/devstorage.read_only
# * https://www.googleapis.com/auth/logging.write
# * https://www.googleapis.com/auth/monitoring
# Both nodes should allow auto_repair and auto_upgrade
# Both nodes should be small and cheap. Let's make them:
# * preemptible
# * 10GB disk size with a standard disk type
# * e2-small machine type for primary and e2-micro for kubeip
# * containerd as the image type

# Finally we need an ingress firewall rule on the default network.
# The rule should allow all incoming TCP traffic to our gke_node_port
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_firewall.html
