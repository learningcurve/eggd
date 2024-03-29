E       := @echo
Q       := @

BASE    := $(abspath .)
INSTALL ?= /usr/local/bin
GOPATH  := $(BASE)

EGGD    := src/eggd

.PHONY: all install clean $(EGGD)

all: $(EGGD)

install:
	$(Q) INSTALL=$(INSTALL) make -C src install

clean:
	$(Q) make -C src clean
	$(Q) rm -rf pkg

$(EGGD):
	git submodule init
	git submodule update
	GOPATH=$(GOPATH) make -C src
