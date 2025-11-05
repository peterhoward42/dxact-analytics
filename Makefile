
# Deploy the cloud run function
#
# The service-injest-event parameter is the name you want GCP to adopt for the `service` that implmeents your function.
# The InjestEvent parameter should be the name you used in your code to register the function.
# The base-image parameter tells GCP which "build pack" it should use to build the Docker image around your code.
# The region parameter specifies where you want the function to run.
.PHONY: deploy
deploy:
	gcloud run deploy service-injest-event --source . --function InjestEvent --base-image go125 --region europe-west2

# Make a test curl request to the google cloud function.
.PHONY: trigger
trigger:
	curl -X POST https://service-injest-event-65030510907.europe-west2.run.app -d '{ \
	"EventULID": "01G65Z755AFWAKHE12NY0CQ9FH", \
	"ProxyUserID": "9e61fcda-ddf5-4294-9b1b-36263317c99f", \
	"TimeUTC": "2023-09-24T15:30:00Z", \
	"Visit": 3, \
	"Event": "this event", \
	"Parameters": "these params" \
	}' \


# Download the entire telemetry event bucket to the local file system.
.PHONY: download
download:
	gcloud storage cp --recursive gs://drawexact-telemetry ~/scratch

# Perform an analysis on the downloaded local copy of the telemetry events
.PHONY: report
report:
	cd analysis/cmd/walker; go run main.go
