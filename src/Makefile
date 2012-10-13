E       := @echo
Q       := @

INSTALLTOOL := install
GO      := go
GOFLAGS := 

INOTIFY := inotify

EGGD    := eggd
EGGD_GO := eggd.go

.PHONY: all install clean $(INOTIFY)

all: $(EGGD)

install:
	$(E) "Installing $(EGGD) to $(INSTALL)..."
	$(Q) $(INSTALLTOOL) -d $(INSTALL)
	$(Q) $(INSTALLTOOL) $(EGGD) $(INSTALL)
	$(E) "Successfully installed $(EGGD)"

clean:
	$(Q) rm -f $(EGGD)

$(INOTIFY):
	$(E) "  GO    $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) install $@

$(EGGD): $(EGGD_GO) $(INOTIFY)
	$(E) "  GO    $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) build $<
