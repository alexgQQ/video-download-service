
# Video Download Service

A web application for downloading just about any video on the internet. Kinda like those youtube video downloader services but better (maybe). The reason behind this is an exploration of tailwind and htmx while using golang for application code. This also provides an opportunity to build some architecture with k8s, GCP and terraform.

## Local Dev Setup

This requires docker and the gcloud cli setup with a project and a bucket with a service account having access to it.
I'd like to remove the need for an actual bucket for a dev setup but there isn't support for bucket emulation like some of other GCP services.

Set the bucket name as the `GCP_BUCKET` value in the `.env` file.

Create a service account key in the `/secrets` dir
```bash
gcloud iam service-accounts keys create secrets/key.json --iam-account service_acct_email
```

As a note, I don't like having to mount a service account key. However the [GCP bucket URL signing](https://cloud.google.com/storage/docs/samples/storage-generate-signed-url-v4) is only supported through service accounts and [impersonating a service account](https://cloud.google.com/docs/authentication/use-service-account-impersonation) with my user creds throws some random auth errors.

Download the [standalone tailwind cli](https://tailwindcss.com/blog/standalone-cli) to the `/client` dir and process the css.
```bash
cd client
# Download whichever binary is appropriate for your machine
curl -o tailwindcss -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss
# Add `--watch` flag for active reloading
./tailwindcss -i templates/css/layout.css -o public/style.css
```

Start up the services

```bash
docker compose up --build -d
```

It'll take a couple seconds for everything to init but the app should be available locally on port 8080,
The client code will auto reload on changes. If any subscriber listener changes are needed then you will need to rebuild it.
```bash
docker compose up --build -d subscriber
```

Tear it down
```bash
docker compose down --remove-orphans -v
```

## Build Images

Build the client image and subscriber listener image then push to the container repository.
Make sure to build the latest css for the client app.
```bash
cd client
docker build -t us.gcr.io/video-download-service/vdl-client:latest .
docker push us.gcr.io/video-download-service/vdl-client:latest

cd ../subscriber
docker build -t us.gcr.io/video-download-service/vdl-subscriber:latest .
docker push us.gcr.io/video-download-service/vdl-subscriber:latest
```

## Terraform Setup

Obviously this requires [terraform](https://developer.hashicorp.com/terraform/downloads) but you'll also need to supply a cloudflare api token and zone id to the tfvars file.

```bash
cd terraform
terraform init
```

Then make changes and plan/apply as needed.

