## Notes

### Initial Thoughts and Considerations

Basic breakdown of the application workflow:
* Enter a video url on a webpage
* Download video source to an available location
* Present user with the download location or any error

Straightforward but lots of gaps to fill. Some initial questions that drove decisions.

**How to process the video download?**

I personally use [youtube-dl](https://github.com/ytdl-org/youtube-dl/tree/master) for downloading video for stupid meme edits and what not. It works well, has lots of config options and despite the name has support for other services like tiktok, twitter and twitch. This makes me consider it a great candidate for the actual downloading process.

The download requests could come in at any time and could be of videos of varying sizes. This creates an async workflow for the client and wouldn't be ideal to run on an application server itself. So I think I should opt to use some messaging service to control the download flow.

A little unknown to think about is limiting the file size. I'm imagining someone trying to download a multi-hour stream vod and how that should be handled.

**Where to host the video file?**

For this I just need some accessible networked file service. Most cloud buckets can fit my needs here and there are lots of CDN services so I'll leave it flexible for now. File storage is short term. Something that I should lookout for is some type of lifecycle management. It would be very beneficial to leverage a service to remove files after a certain amount of time instead of writing some cronjob to handle it.

**How to return the download to the client?**

With a messaging service the client is decoupled from the actual downloading and I'll need a way for the client to get the download status and location. I think I'll need some sort of service to store results from the messaging service. A database feels like overkill, I don't need the storage to be long term or anything.


So at a high level this will just be a web client that submits download jobs to a messaging service. That service downloads the video to a networked location and stores results to a memory instance for the client to poll.

Lets refine this a bit.

I think the messaging service is the most important choice and there are lots of options. I'm choosing to use GCP Pub/Sub as it's simple, cheap, familiar and guarantees an "always once" delivery. This also locks me into the GCP ecosystem, so something to consider, but I don't forsee a major issue there. I'll need to have some listener running to process messages as they come in but this can be done with various GCP services.

For a memory service Redis is an easy choice. It is flexible, familiar and lots of ways to host. I can also leverage TTL capabilities for automatic cleanup. One note is that I am avoiding using the GCP's hosted Redis product as it incurs a flat rate usage price. I considered Datastore but it has a lot of quirks and charges by a query quota so no dice there.

If I'm not going to use the dedicated GCP Redis service then I need to think about hosting it. The two options I see are to use a Compute Instance or a GKE cluster to run it.

Going in theme with GCP products, I think their buckets should be sufficient. Simple blob storage with lifecycle management and access control. There are options to front it with a load balancer to make it a real CDN but I don't think I'll need it initially.

For the app itself I'd like to use Go for the application code along with Tailwind and HTMX for the frontend. These are some technologies I'd like to tinker with and this is a good opportunity. Tailwind and HTMX have a lot of features I like and I can always use more Golang practice. I'd like to avoid using any Go frameworks as this is simple enough to execute with the default http package.

I have a lot of options for hosting the app. It would be good to stay within GCP so from there I can use Cloud Run, App Engine, Compute instance or GKE. Now since I'm a fan of k8s and can use a GKE cluster to run the client, redis and the pubsub listener so I'm going to go with that.

Finally I think it is always good practice to have infra as code so lets orchestrate it all with terraform.

So here's the plan:
    * Web app written in Golang and bundled with Tailwind and HTMX
    * GCP PubSub instance with a single topic and pull subscription
    * GCP Bucket for downloadable video files
    * Redis instance to store download results
    * GKE cluster to run redis, webapp and subscription listener
    * Manage infrastructure with terraform

One last thing that I'll do to avoid expensive GCP load balancers is to use cloudflare for DNS and leverage it's port routing ruleset to direct traffic to our GKE node.
