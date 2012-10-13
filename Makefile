E       := @echo
Q       := @

BASE    := $(abspath .)
INSTALL ?= /usr/local/bin
GOPATH  := $(BASE)

EGGD    := src/eggd

.PHONY: all install clean

all: $(EGGD)

install:
	$(Q) INSTALL=$(INSTALL) make -C src install

clean:
	$(Q) make -C src clean

$(EGGD):
	$(Q) GOPATH=$(GOPATH) make -C src $@
