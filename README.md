# This repo is a Google Cloud Function for DrawExact analytics

We want to track events when someone new tries out DrawExact, and to observe their subsequent usage patterns.
This cloud function forms the REST API to receive these events.

## User Identity

- We need a way to bind incoming events to a particular - but anonymised user.
- Anonymity is important to satisfy privacy requirements.
- And we are particularly interested in what happens when people try DrawExact for the first time.
- These new users won't be signed in to DrawExact. And we don't know who they are.
- It follows that their Google sign-in identity cannot be the root basis for the identity key.

So we will create a unique UUID identity proxy for a user the first time DrawExact runs.
We'll detect if it is the first invocation by storing a flag in browser local storage. With this identity
being device and browser specific - it is an imperfect proxy model for a user, but a reasonable, and good enough trade off between accuracy and simplicity for the intended purpose.

# References

Original source Material:

https://cloud.google.com/run/docs/write-functions

https://github.com/GoogleCloudPlatform/functions-framework-go

# Architecture

The cloud function run entry point is the `./cmd.main()`

The function business logic entry point is `.functions/injestEvent()`

The code depends on, conforms to, and is orchestrated by Google's Functions-as-a-Service framework.

Note the routing / binding of a declared function name to the go function in the `init()` function for
the `./functions` module.

We use Google's *build from source* option - which lets you use `gcloud` to delegate the (Docker) build process to Google in the cloud.

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
- This bypasses the docker build, and tests your function handler embedded in a 
  development server, using a native go workflow.
- See the references

# Deploy it as a cloud function

Use the `gcloud` CLI - see Makefile `deploy` target. 

# Check it is running and working

Use `curl`  - see Makefile `triger` target. 

# CORS
The current code does not specify CORS headers, and that is fine for curl requests of course.
But for browser/javascript requests, CORS will need to be coded in the function. (See references).



