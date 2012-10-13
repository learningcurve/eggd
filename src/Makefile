E       := @echo
Q       := @

INSTALLTOOL := install
GO      := go build
GOFLAGS := 

EGGD    := eggd
SOURCES := eggd.go

.PHONY: all install clean

all: $(EGGD)

install:
	$(E) "Installing $(EGGD) to $(INSTALL)..."
	$(Q) $(INSTALLTOOL) -d $(INSTALL)
	$(Q) $(INSTALLTOOL) $(EGGD) $(INSTALL)
	$(E) "Successfully installed $(EGGD)"

clean:
	$(Q) rm -f $(EGGD)

$(EGGD): $(SOURCES)
	$(E) "  GO    $@"
	$(Q) make -C inotify
	$(Q) GOPATH=$(GOPATH) $(GO) $(GOFLAGS) $<
	$(E) "Built $@"
