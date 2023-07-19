# Hiii!
# Here we'll need a cloud bucket and iam access control for it.

# We'll need a bucket with a name of "vdl-primary-bucket" and should use our gke_zone and gke_region
# The bucket should also have lifecycle management to delete anything after 1 day.
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket#example-usage---life-cycle-settings-for-storage-bucket-objects

# We'll also want a service account for along with a role for it's perms with a binding to the service account
# Give the service account a descriptive name but make sure it's role hase these perms:
# * pubsub.subscriptions.consume
# * pubsub.topics.publish
# * storage.objects.create
# * storage.objects.get
# * storage.objects.list
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam_custom_role.html
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam.html#google_project_iam_binding
