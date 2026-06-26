.PHONY: build install uninstall clean

BINARY  := coremanager
ALIAS   := cman
PREFIX  ?= /usr/local
BINDIR  := $(PREFIX)/bin

build:
	go build -o $(BINARY) .

install: build
	install -Dm755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	ln -sf $(BINDIR)/$(BINARY) $(DESTDIR)$(BINDIR)/$(ALIAS)

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)
	rm -f $(DESTDIR)$(BINDIR)/$(ALIAS)

clean:
	rm -f $(BINARY)
