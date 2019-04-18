# pubsubhub

Subscribe all Cloud Pub/Sub in the project and refresh subscription list periodically

1. Retrieve all [pull subscriptions](https://cloud.google.com/pubsub/docs/pull) in the given GCP project
2. Start receiving message from all subscriptions
3. After 1 minute, stop receiving message
4. Return to 1.

## Usage

```bash
pubsubhub <GCP_PROJECT_ID>
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))
