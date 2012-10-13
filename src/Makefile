E       := @echo
Q       := @

INSTALLTOOL := install
GO      := go
GOFLAGS := 

CONF    := conf
INOTIFY := inotify

EGGD    := eggd
EGGD_GO := eggd.go
EGGD_LIBS := $(CONF) $(INOTIFY)

.PHONY: all install clean $(EGGD_LIBS)

all: $(EGGD)

install:
	$(E) "Installing $(EGGD) to $(INSTALL)..."
	$(Q) $(INSTALLTOOL) -d $(INSTALL)
	$(Q) $(INSTALLTOOL) $(EGGD) $(INSTALL)
	$(E) "Successfully installed $(EGGD)"

clean:
	$(Q) rm -f $(EGGD)

$(CONF):
	$(E) "  GO    $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) install $@

$(INOTIFY):
	$(E) "  GO    $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) install $@

$(EGGD): $(EGGD_GO) $(EGGD_LIBS)
	$(E) "  GO    $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) build $<
