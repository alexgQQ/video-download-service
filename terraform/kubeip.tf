# Hallo!
# Here we need resources to make kubeip operate

# We'll need a single external ip address from google with:
# * a label property as kubeip=primary-cluster
# * name prefixed with "kubeip-"
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address.html

# We'll also want a service account for kubeip along with a role for it's specific perms with a binding to the service account
# The specific permissions the role needs are:
# * compute.addresses.list
# * compute.instances.addAccessConfig
# * compute.instances.deleteAccessConfig
# * compute.instances.get
# * compute.instances.list
# * compute.projects.get
# * container.clusters.get
# * container.clusters.list
# * resourcemanager.projects.get
# * compute.networks.useExternalIp
# * compute.subnetworks.useExternalIp
# * compute.addresses.use
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam_custom_role.html
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam.html#google_project_iam_binding