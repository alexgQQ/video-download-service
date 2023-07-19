# Howdy!
# Here we need a single PubSub topic and subscriber

# The topic is super simple and we only need it's name set as "vdl-primary-topic"
# All other options can be left for default for now.
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_topic

# The subscription is also pretty simple. It should be a pull subscription named "vdl-primary-subscription"
# and be bound to the topic above.
# We don't need to retain acked messages and we don't need to use message ordering.
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_subscription#example-usage---pubsub-subscription-pull
