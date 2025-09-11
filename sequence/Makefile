# Fixed registry
REGISTRY := docker.io/changyuan
IMAGE_NAME := kn-event-step-demo
IMAGE := $(REGISTRY)/$(IMAGE_NAME):latest

# Namespace (default: default), override via: make all NAMESPACE=my-ns
NAMESPACE ?= default

# Curl image version for sending CloudEvents
CURL_IMG := curlimages/curl:8.10.1

# ------- Main entry point -------
.PHONY: all
all: login docker-build docker-push deploy show

# Optional: Use environment variables for automatic docker login (skipped if already logged in)
.PHONY: login
login:
	@if [ -n "$$HARBOR_USER" ] && [ -n "$$HARBOR_PASS" ]; then \
	  echo "$$HARBOR_PASS" | docker login $(REGISTRY) -u "$$HARBOR_USER" --password-stdin ; \
	else \
	  echo "Skipping docker login (set HARBOR_USER/HARBOR_PASS to enable)"; \
	fi

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) .

.PHONY: docker-push
docker-push:
	docker push $(IMAGE)

# Use envsubst to inject namespace into YAML URIs
.PHONY: deploy
deploy:
	@which envsubst >/dev/null || (echo "Please install 'envsubst' (gettext)"; exit 1)
	env NAMESPACE="$(NAMESPACE)" envsubst < ksvc-sequence.yaml | kubectl apply -n $(NAMESPACE) -f -
	@echo "Waiting 10 seconds for Knative Services to become ready..."
	@sleep 10

.PHONY: undeploy
undeploy:
	@which envsubst >/dev/null || (echo "Please install 'envsubst' (gettext)"; exit 1)
	- env NAMESPACE="$(NAMESPACE)" envsubst < ksvc-sequence.yaml | kubectl delete -n $(NAMESPACE) -f -

.PHONY: show
show:
	@echo "----- Knative Services -----"
	kubectl get ksvc -n $(NAMESPACE)
	@echo "----- Sequence Address -----"
	kubectl get sequence five-step-seq -n $(NAMESPACE) -o jsonpath='{.status.address.url}'; echo

# --------- Invocation / Integration ---------

# Send a CloudEvent to Sequence (from within the cluster)
.PHONY: send
send:
	@URL=$$(kubectl get sequence five-step-seq -n $(NAMESPACE) -o jsonpath='{.status.address.url}'); \
	echo "Sending CloudEvent to $$URL"; \
	kubectl run ce-sender --image=$(CURL_IMG) -n $(NAMESPACE) --rm -it --restart=Never -- \
	  -v -X POST "$$URL" \
	  -H 'Content-Type: application/json' \
	  -H 'Ce-Id: demo-001' \
	  -H 'Ce-Specversion: 1.0' \
	  -H 'Ce-Source: cli.make' \
	  -H 'Ce-Type: dev.demo' \
	  -d '{"hello":"sequence"}'

# Test any single Knative Service: make test-step STEP=3 (default: 1)
STEP ?= 1
.PHONY: test-step
test-step:
	@URL="http://step$(STEP).$(NAMESPACE).svc.cluster.local/process"; \
	echo "POST $$URL"; \
	kubectl run curl-step$(STEP) --image=$(CURL_IMG) -n $(NAMESPACE) --rm -it --restart=Never -- \
	  -v -X POST "$$URL" \
	  -H 'Content-Type: application/json' \
	  -H 'Ce-Id: single-$(STEP)' \
	  -H 'Ce-Specversion: 1.0' \
	  -H 'Ce-Source: cli.make' \
	  -H 'Ce-Type: dev.single' \
	  -d "$$(printf '{"only":"step%s"}' $(STEP))"

# Tail logs: user-container for all steps
.PHONY: logs
logs:
	kubectl logs -n $(NAMESPACE) -l 'serving.knative.dev/service in (step1,step2,step3,step4,step5)' -c user-container --tail=100 -f

.PHONY: clean
clean:
	@echo "Cleaning resources in namespace $(NAMESPACE) ..."
	# 1) Delete by original manifest (ksvc + sequence), fallback even if YAML changed
	- @which envsubst >/dev/null || (echo "Please install 'envsubst' (gettext)"; exit 1)
	- env NAMESPACE="$(NAMESPACE)" envsubst < ksvc-sequence.yaml | kubectl delete -n $(NAMESPACE) -f - --ignore-not-found

	# 2) Explicit fallback deletion (removes resources even if YAML names changed)
	- kubectl -n $(NAMESPACE) delete sequence five-step-seq --ignore-not-found
	- kubectl -n $(NAMESPACE) delete ksvc step1 step2 step3 step4 step5 final-sink --ignore-not-found

	# 3) Delete InMemoryChannel/Subscription created by Sequence (by label and name)
	- kubectl -n $(NAMESPACE) delete inmemorychannel,subscription -l 'flows.knative.dev/sequence=five-step-seq' --ignore-not-found
	- kubectl -n $(NAMESPACE) get inmemorychannels,subscriptions -o name \
	    | grep -E 'five-step-seq|kn-sequence' \
	    | xargs -r kubectl -n $(NAMESPACE) delete

	# 4) Delete temporary sender/test pods (ce-sender, curl-step1..5)
	- kubectl -n $(NAMESPACE) delete pod -l run=ce-sender --ignore-not-found
	- for i in 1 2 3 4 5; do kubectl -n $(NAMESPACE) delete pod -l run=curl-step$$i --ignore-not-found; done

	@echo "Clean done."
