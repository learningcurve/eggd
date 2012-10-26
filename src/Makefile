E       := @echo
Q       := @

INSTALLTOOL := install
UPSTART := /etc/init
GO      := /usr/local/go/bin/go
GOFLAGS := 

CONF    := conf
INOTIFY := inotify

EGGD    := eggd
EGGD_GO := eggd.go
EGGD_LIBS := $(CONF) $(INOTIFY)
EGGD_CONF := eggd.conf

.PHONY: all install clean $(EGGD_LIBS)

all: $(EGGD)

install: $(EGGD_CONF)
	$(Q) $(INSTALLTOOL) -d $(INSTALL)
	$(Q) $(INSTALLTOOL) -t $(INSTALL) $(EGGD)
	$(E) "  Successfully installed $(EGGD) to $(INSTALL)"
	$(Q) $(INSTALLTOOL) -t $(UPSTART) $(EGGD_CONF)
	$(E) "  Successfully installed $(EGGD_CONF) to $(UPSTART)"

clean:
	$(Q) rm -f $(EGGD)
	$(Q) rm -f $(EGGD_CONF)

$(EGGD): $(EGGD_GO) $(EGGD_LIBS)
	$(E) "  GO $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) build $<

$(CONF):
	$(E) "  GO $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) install $@

$(INOTIFY):
	$(E) "  GO $@"
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) install $@

$(EGGD_CONF):
	$(E) "  Generating $@"
	$(Q) rm -f $(EGGD_CONF)
	$(E) "start on startup" >> $(EGGD_CONF)
	$(E) "exec $(INSTALL)/$(EGGD)" >> $(EGGD_CONF)
	$(E) "console none" >> $(EGGD_CONF)
