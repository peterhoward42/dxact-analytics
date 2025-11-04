# This repo is in support of DrawExact analytics

We hope to detect DrawExact usage patterns that will help us infer information about DrawExact's user adoption and satisfaction. For example:
- the proportion of new users who complete the training cage process
- the proportion of new users who choose to sign in after the training

The DrawExact app emits suitable telemetry as HTTP POST requests to a Google Cloud Run Function.

This repo implements that Cloud Run Function, and exports the event schema payload schema for use by the 
DrawExact app code (as the ./lib/EventPayload type).

## Anonymity

- We need a way to associate incoming events to the originating DrawExact user - but ANONYMOUSLY.
- Anonymity is important to satisfy privacy requirements.
- And we are particularly interested in what happens when people try DrawExact for the first time.
- During this phase of usage, people are unlikely to be signed in to DrawExact. So we don't know who they are.
- From the previous points, it follows that a user's Google sign-in identity cannot be the root basis for the identity key.

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

The code depends on Google's Functions-as-a-Service framework.

We use Google's *build from source* option to build and deploy the cloud function - which lets you use `gcloud` to bypass needing to do the Docker build process locally, and instead delegate that to Google in the cloud.

# Configure / Build and Run

There are three steps.

For the details - see the references above.

1. Create and configure a GCP project to own the function. 
2. Configure the `gcloud` CLI tool locally - (see references)
4. Deploy it as a *Google Cloud Run* function.

Note that testing the function locally is no longer convenient, because the cloud function now writes to Google Cloud Storage, and that depends on the default service account identity being known - which it is not when run locally.

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

```

# Deploy it as a cloud function

See the Makefile `deploy` target. 

# Check it is running and working

See the Makefile `triger` target. 


