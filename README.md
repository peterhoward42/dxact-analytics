# This repo is in support of DrawExact analytics

We hope to detect DrawExact usage patterns that tell us something about DrawExact's user adoption and obstacles to conversion. For example:
- the proportion of new users who complete the training cage process
- the proportion of new users who choose to sign in after the training

The DrawExact app emits suitable telemetry as HTTP POST requests to a Google Cloud Run Function.

This repo implements that Cloud Run Function, and exports the event schema payload schema (as the ./lib/EventPayload type).

## Anonymity

- We need a way to bind incoming events to the originating DrawExact user - but ANONYMOUSLY.
- Anonymity is important to satisfy privacy requirements.
- And we are particularly interested in what happens when people try DrawExact for the first time.
- During this phase of usage, people are unlikely to be signed in to DrawExact. So we don't know who they are.
- It follows that their Google sign-in identity cannot be the root basis for the identity key.

So we will create an anonymous unique UUID identity proxy for a user the first time DrawExact runs.
We'll detect if it is the first invocation by storing a flag in their browser local storage. With this identity
being device and browser specific - it is an imperfect proxy model for a user, but a reasonable, and good enough trade off between accuracy and simplicity for the intended purpose.

# References for Google Cloud Functions

Original source Material:

https://cloud.google.com/run/docs/write-functions

https://github.com/GoogleCloudPlatform/functions-framework-go

# Architecture

The cloud function run entry point is the `./cmd.main()`

The function business logic entry point is `.functions/injestEvent()`

The code depends on, conforms to, and is orchestrated by Google's Functions-as-a-Service framework.

Note the routing / binding of a declared function name to the go function in the `init()` function for
the `./functions` module.

We use Google's *build from source* option - which lets you use `gcloud` to bypass needing to do the Docker build process locally, and instead delegate that to Google in the cloud.

# Configure / Build and Run

There are four steps.

For the details - see the references above.

1. Create and configure a GCP project to own the function. 
2. Configure the `gcloud` CLI tool locally - (see references)
3. Test your function locally
4. Deploy it as a *Google Cloud Run* function.

# Configure the GCP Project
In overview this is to:
- create a new project
- enable the Cloud Run API
- enable set of other required APIs
- grant a set of operational permissions to the deployer identity (your email address)
- grant a set of operational permissions to the default google cloud run service account identity

# Configuring gcloud
In overview this is to:

- update gcloud
- set config for:
    - your identity (gmail address)
    - your current gcp project
    - the region in which to run Cloud Run functions

# Test your function locally
- This uses the open-source `functions-framework` provided by Google.
- The framework lets us test the function handler, (without a Docker build process) directly -  embedded in a 
  local development server, using a native go workflow.

```
FUNCTION_TARGET=InjestEvent  go run cmd/main.go

// In another terminal:

curl -X POST \
localhost:8080 \
-d '{
"ProxyUserId": "thisuuid", 
"TimeUTC": "thisTime",
"Visit": 3, 
"Event": "this event", 
"Parameters": 
"these params"
}'

```

# Deploy it as a cloud function

See the Makefile `deploy` target. 

# Check it is running and working

See the Makefile `triger` target. 


