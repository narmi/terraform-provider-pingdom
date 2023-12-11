PACKAGE_VERSION=1.0.0

PINGDOM_SPEC_URL=https://docs.pingdom.com/api/API_3.1.yaml
PINGDOM_SPEC_FILE=specs/pingdom-v3.1.yaml
PINGDOM_GO_FILE=pingdom/pingdom.gen.go

GIT_ORG=narmi
GIT_REPO=go-pingdom


default: pingdom_spec

$(PINGDOM_SPEC_FILE):
	@echo "Downloading Pingdom spec file"
	@curl -s -o $(PINGDOM_SPEC_FILE) $(PINGDOM_SPEC_URL)

pingdom_spec: $(PINGDOM_SPEC_FILE)
	@echo "Generating Pingdom client"
	@oapi-codegen -package pingdom $(PINGDOM_SPEC_FILE) > $(PINGDOM_GO_FILE)

.PHONY: pingdom_spec_download pingdom_spec
