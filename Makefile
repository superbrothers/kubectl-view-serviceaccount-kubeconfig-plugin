VIEW_KUBECONFIG := view-kubeconfig
ARCHIVE_FILE := view-kubeconfig.zip

$(VIEW_KUBECONFIG): main.go Gopkg.lock
		@go build -o $@ .

.PHONY: clean
clean:
		@$(RM) $(VIEW_KUBECONFIG) $(ARCHIVE_FILE)

.PHONY: archive
archive: $(VIEW_KUBECONFIG) LICENSE plugin.yaml
		@zip -r $(ARCHIVE_FILE) $^
