# Makefile for Sigil.
# 
# Run 'make' to build locally, 'make install' to install binaries and other data, 'make package'
# to prepare a redistributable package and 'make uninstall' to remove installed data.

COMPILER = gc
PROGRAM  = sigil
REPO     = github.com/deuill/sigil

VERSION  = $(shell git describe --tags | cut -c3-)
SERVICES = $(shell find srv/* -maxdepth 1 -type d)
ENGINES  = $(shell find lib/engine/* -maxdepth 1 -type d)

.PHONY: $(PROGRAM) depend
all: $(PROGRAM)

depend:
	$(shell echo "package main"  > imports.go)
	$(foreach srv, $(SERVICES), $(shell echo "import _ \"$(REPO)/$(srv)\""  >> imports.go))
	$(foreach ngn, $(ENGINES),  $(shell echo "import _ \"$(REPO)/$(ngn)\""  >> imports.go))

$(PROGRAM): depend
	@echo "Building '$(PROGRAM)'..."

	@mkdir -p .tmp
	@go build -compiler $(COMPILER) -o .tmp/$(PROGRAM)

install:
	@echo "Installing '$(PROGRAM)' and data..."

	@install -Dm 644 dist/init/systemd/$(PROGRAM).service $(DESTDIR)/usr/lib/systemd/system/$(PROGRAM).service

	@install -dm 0750 $(DESTDIR)/etc/$(PROGRAM)
	@install -m 0640 dist/conf/* $(DESTDIR)/etc/$(PROGRAM)

	@install -Dsm 0755 .tmp/$(PROGRAM) $(DESTDIR)/usr/bin/$(PROGRAM)

package:
	@echo "Building package for '$(PROGRAM)'..."

	@mkdir -p .tmp/package
	@make DESTDIR=.tmp/package install
	@fakeroot -- tar -cJf $(PROGRAM)-$(VERSION).tar.xz -C .tmp/package .

uninstall:
	@echo "Uninstalling '$(PROGRAM)'..."

	@rm -f $(DESTDIR)/usr/bin/$(PROGRAM)
	@rm -f $(DESTDIR)/usr/lib/systemd/system/$(PROGRAM).service
	@mv -f $(DESTDIR)/etc/$(PROGRAM) $(DESTDIR)/etc/$(PROGRAM).save

	@echo -e "\033[1mConfiguration files moved to '$(DESTDIR)/etc/$(PROGRAM).save'.\033[0m"

clean:
	@echo "Cleaning '$(PROGRAM)'..."

	@go clean
	@rm -Rf .tmp imports.go