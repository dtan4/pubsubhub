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

```sh-session
$ bin/pubsubhub example-project
2019/04/18 21:18:59 Retrieving subscriptions...
2019/04/18 21:19:03 Start subscribing subscription-a...
2019/04/18 21:19:03 Start subscribing subscription-b...
2019/04/18 21:19:54 [subscription:subscription-a message:951063581465360] body:"foobarbaz"
2019/04/18 21:19:54 [subscription:subscription-a message:951063581465360] ACK
2019/04/18 21:19:59 [subscription:subscription-b] Stopping...
2019/04/18 21:20:00 [subscription:subscription-a] Stopping...
2019/04/18 21:20:00 Retrieving subscriptions...
2019/04/18 21:20:02 Start subscribing subscription-a...
2019/04/18 21:20:02 Start subscribing subscription-b...
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))
