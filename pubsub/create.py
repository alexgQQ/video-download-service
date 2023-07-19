"""
Creates a single topic and pull subscription for a GCP pubsub instance.
Intended to be used with the pubsub emulator while building a local dev env.
"""

import os
import argparse

from google.cloud import pubsub_v1
from google.api_core.exceptions import AlreadyExists


def create_topic(project_id: str, topic_id: str):
    publisher = pubsub_v1.PublisherClient()
    topic_path = publisher.topic_path(project_id, topic_id)

    try:
        topic = publisher.create_topic(request={"name": topic_path})
        print(f"Created topic: {topic.name}")
    except AlreadyExists:
        print(f"Topic: {topic_path} already exists")
    except Exception as err:
        print(f"Topic creation failed: {err}")


def create_subscription(project_id: str, topic_id: str, subscription_id: str):
    publisher = pubsub_v1.PublisherClient()
    subscriber = pubsub_v1.SubscriberClient()
    topic_path = publisher.topic_path(project_id, topic_id)
    subscription_path = subscriber.subscription_path(project_id, subscription_id)

    try:
        with subscriber:
            subscription = subscriber.create_subscription(
                request={"name": subscription_path, "topic": topic_path}
            )
        print(f"Subscription created: {subscription}")
    except AlreadyExists:
        print(f"Subscription: {subscription_path} already exists")
    except Exception as err:
        print(f"Subscription creation failed: {err}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "project_id", nargs="?", default=os.environ.get("GCP_PROJECT"), type=str
    )
    parser.add_argument(
        "topic_id", nargs="?", default=os.environ.get("GCP_TOPIC_ID"), type=str
    )
    parser.add_argument(
        "subscription_id", nargs="?", default=os.environ.get("GCP_SUBSCRIBER_ID"), type=str
    )
    args = parser.parse_args()

    create_topic(args.project_id, args.topic_id)
    create_subscription(args.project_id, args.topic_id, args.subscription_id)
