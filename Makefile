OUT_DIR ?= ./_out
DIST_DIR ?= ./_dist
$(shell mkdir -p _dist)

.PHONY: build
build:
		go build -o view-kubeconfig .

.PHONY: cross
build-cross: $(OUT_DIR)/linux-amd64/view-kubeconfig $(OUT_DIR)/darwin-amd64/view-kubeconfig

.PHONY: dist
dist: $(DIST_DIR)/view-kubeconfig-linux-amd64.zip $(DIST_DIR)/view-kubeconfig-darwin-amd64.zip

.PHONY: checksum
checksum:
		for f in _dist/*.zip; do \
			shasum -a 256 "$${f}"  | awk '{print $$1}' > "$${f}.sha256" ; \
		done

.PHONY: clean
clean:
		rm -rf $(OUT_DIR) $(DIST_DIR) view-kubeconfig

$(OUT_DIR)/%-amd64/view-kubeconfig:
		GOOS=$* GOARCH=amd64 go build -o $@ .

$(DIST_DIR)/view-kubeconfig-%-amd64.zip: $(OUT_DIR)/%-amd64/view-kubeconfig
		( \
			cd $(OUT_DIR)/$*-amd64/ && \
			cp ../../version.txt . && \
			cp ../../LICENSE . && \
			cp ../../README.md . && \
			cp ../../plugin.yaml . && \
			zip -r ../../$@ * \
		)
