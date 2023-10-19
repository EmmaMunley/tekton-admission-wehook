.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/admission-webhook-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o bin/admission-webhook-linux-amd64 .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building tekton-webhook Docker image..."
	docker build -t tekton-webhook:latest .

# From this point `kind` is required
.PHONY: cluster
cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	minikube start

.PHONY: delete-cluster
delete-cluster:
	@echo "\n♻️  Deleting Kubernetes cluster..."
	minikube delete

.PHONY: push
push: docker-build
	@echo "\n📦 Pushing admission-webhook image into Minikube's Docker daemon..."
	minikube image load tekton-webhook:latest

.PHONY: deploy-config
deploy-config:
	@echo "\n⚙️  Applying cluster config..."
	kubectl apply -f dev/manifests/cluster-config/

.PHONY: delete-config
delete-config:
	@echo "\n♻️  Deleting Kubernetes cluster config..."
	kubectl delete -f dev/manifests/cluster-config/

.PHONY: deploy
deploy: push delete deploy-config
	@echo "\n🚀 Deploying tekton-webhook..."
	kubectl apply -f dev/manifests/webhook/

.PHONY: install-tekton
install-tekton:
	@echo "\n⚙️  Installing tekton..."
	kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml

.PHONY: delete
delete:
	@echo "\n♻️  Deleting tekton-webhook deployment if existing..."
	kubectl delete -f dev/manifests/webhook/ || true

.PHONY: valid-pipeline
valid-pipeline:
	@echo "\n🚀 Deploying \"valid\" pipeline..."
	kubectl apply -f dev/manifests/pipelines/valid-pipeline.yaml

.PHONY: delete-valid-pipeline
delete-valid-pipeline:
	@echo "\n🚀 Deleting \"valid\" pipeline..."
	kubectl delete -f dev/manifests/pipelines/valid-pipeline.yaml

.PHONY: invalid-pipeline
invalid-pipeline:
	@echo "\n🚀 Deploying \"invalid\" pipeline..."
	kubectl apply -f dev/manifests/pipelines/invalid-pipeline.yaml

.PHONY: delete-invalid-pipeline
delete-invalid-pipeline:
	@echo "\n🚀 Deleting \"invalid\" pipeline..."
	kubectl delete -f dev/manifests/pipelines/invalid-pipeline.yaml

.PHONY: valid-task
valid-task:
	@echo "\n🚀 Deploying \"valid\" task..."
	kubectl apply -f dev/manifests/tasks/valid-task.yaml

.PHONY: delete-valid-task
delete-valid-task:
	@echo "\n🚀 Deleting \"valid\" task..."
	kubectl delete -f dev/manifests/tasks/valid-task.yaml

.PHONY: invalid-task
invalid-task:
	@echo "\n🚀 Deploying \"invalid\" task..."
	kubectl apply -f dev/manifests/tasks/invalid-task.yaml

.PHONY: delete-invalid-task
delete-invalid-task:
	@echo "\n🚀 Deleting \"invalid\" task..."
	kubectl delete -f dev/manifests/tasks/invalid-task.yaml

.PHONY: taint
taint:
	@echo "\n🎨 Taining Kubernetes node.."
	kubectl taint nodes kind-control-plane "acme.com/lifespan-remaining"=4:NoSchedule

.PHONY: logs
logs:
	@echo "\n🔍 Streaming tekton-webhook logs..."
	kubectl logs -l app=tekton-webhook -f

.PHONY: delete-all
delete-all: delete delete-config delete-valid-pipeline delete-invalid-pipeline delete-valid-task delete-invalid-task
