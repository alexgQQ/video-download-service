import os
import subprocess
from concurrent.futures import TimeoutError
from dataclasses import asdict, dataclass, fields
from typing import Optional

import redis
from google.cloud import pubsub_v1


@dataclass
class Config:
    """
    My hacky way to coalesce env vars into an object
    with a touch of validation
    """

    project: str = os.getenv("GCP_PROJECT")
    subscription: str = os.getenv("GCP_SUBSCRIBER_ID")
    bucket: str = os.getenv("GCP_BUCKET")
    redis_host: str = os.getenv("REDIS_HOST")
    redis_port: str = os.getenv("REDIS_PORT")
    download_host: str = os.getenv("DOWNLOAD_HOST")

    def __post_init__(self):
        for field in fields(self):
            attr = field.name
            # Checks for None or empty string
            if not getattr(self, attr):
                raise RuntimeError(f"Config value `{attr}` is not set")


config = Config()


@dataclass
class Download:
    """Represents downloads results to store in redis"""

    download_url: str
    original_url: str
    complete: bool

    def __str__(self) -> str:
        return "Download"


def insert_redis_entry(uuid: str, entry: Download):
    r = redis.Redis(
        host=config.redis_host, port=config.redis_port, decode_responses=True
    )
    mapping = asdict(entry)
    # Redis doesn't support straight boolean
    mapping["complete"] = int(mapping["complete"])
    r.hset(uuid, mapping=mapping)


def handle_message(message: pubsub_v1.subscriber.message.Message):
    message.ack()
    print("Message acknowledged")
    uuid = message.attributes.get("uuid")
    url = message.attributes.get("url")

    try:
        print(f"Downloading video {url}")
        # A few problems when tinkering with this
        # # 1 the fake gcs server I'm using for local emulation
        # doesn't work entirely with copying from stdin so it'll fail
        # cmd = f"youtube-dl -o - {url} | gsutil -o 'Credentials:gs_json_host=bucket' -o 'Credentials:gs_json_port=4443' -o 'Boto:https_validate_certificates=False' cp - gs://pixelpopart-public/{uuid}.mp4"
        cmd = f"youtube-dl -o - {url} | gsutil cp - gs://{config.bucket}/vdl-test/{uuid}.mp4"
        # # 2 usually we should split the command string but for some
        # reason it does not get parsed right with the bulky command
        # It's baffling as it runs fine when done manually in the related image
        # but only fails when running as a message callback ¯\_(ツ)_/¯
        subprocess.check_output(cmd, shell=True)
        success = True
        download_url = f"{config.download_host}/{uuid}.mp4"
    except Exception as err:
        print(f"Failed to download video: {err.cmd} {err}")
        success = False
        download_url = ""

    data = Download(download_url, url, success)
    insert_redis_entry(uuid, data)


def receive_messages(timeout: Optional[float] = None):
    """Receives messages from a pull subscription."""

    subscriber = pubsub_v1.SubscriberClient()
    subscription_path = subscriber.subscription_path(
        config.project, config.subscription
    )
    streaming_pull_future = subscriber.subscribe(
        subscription_path, callback=handle_message
    )
    print(f"Listening for messages on {subscription_path}..\n")

    with subscriber:
        try:
            streaming_pull_future.result(timeout=timeout)
        except TimeoutError:
            streaming_pull_future.cancel()  # Trigger the shutdown.
            streaming_pull_future.result()  # Block until the shutdown is complete.


if __name__ == "__main__":
    receive_messages()
