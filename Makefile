.PHONY: all docker-build hilink-reconnect install install-failover install-all uninstall uninstall-failover uninstall-all

all: docker-build hilink-reconnect

docker-build:
	docker build -t multiproxy .

hilink-reconnect:
	make -C hilink-reconnect docker-build

install:
	if [ -d /lib/dhcpcd ]; then \
		cp dhcpcd/dhcpcd-hook /lib/dhcpcd/dhcpcd-hooks/99-multiproxy; \
	fi

	cp systemd/multiproxy-setup-firewall.service /etc/systemd/system
	cp systemd/multiproxy-setup-tables.service /etc/systemd/system

	if systemctl is-enabled -q dhcpcd 2>/dev/null; then \
		systemctl restart dhcpcd; \
	fi

	systemctl enable --now multiproxy-setup-firewall multiproxy-setup-tables

	docker-compose up -d

install-failover:
	if [ -d /lib/dhcpcd ]; then \
		cp dhcpcd/dhcpcd-hook-fallback /lib/dhcpcd/dhcpcd-hooks/99-multiproxy-fallback; \
	fi

	cp systemd/multiproxy-simplefailover.service /etc/systemd/system
	systemctl enable --now multiproxy-simplefailover

install-all: install install-failover

uninstall:
	docker-compose down

	systemctl disable --now multiproxy-setup-firewall multiproxy-setup-tables

	rm /etc/systemd/system/multiproxy-setup-firewall.service
	rm /etc/systemd/system/multiproxy-setup-tables.service

	while IFS='' read -r line || [ -n "$$line" ]; do \
		[ -n "$$line" ] || continue; \
		echo "$$line" | grep -qv "^#" || continue; \
		\
		GATEWAY=$$(echo "$$line" | awk -F '\t' '{print $$1}'); \
		INTERFACE=$$(echo "$$line" | awk -F '\t' '{print $$2}'); \
		TABLE=$$(echo "$$line" | awk -F '\t' '{print $$4}'); \
		\
		if [ -n "$$GATEWAY" -a -n "$$INTERFACE" ]; then \
			ip route del default via "$$GATEWAY" dev "$$INTERFACE" table "$$TABLE"; \
		elif [ -n "$$GATEWAY" ]; then \
			ip route del default via "$$GATEWAY" table "$$TABLE"; \
		elif [ -n "$$INTERFACE" ]; then \
			ip route del default dev "$$INTERFACE" table "$$TABLE"; \
		fi; \
	done < /opt/multiproxy/instances

	if [ -d /lib/dhcpcd ]; then \
		rm /lib/dhcpcd/dhcpcd-hooks/99-multiproxy; \
	fi

	if systemctl is-enabled -q dhcpcd 2>/dev/null; then \
		systemctl restart dhcpcd; \
	fi

	docker image rm -f multiproxy
	docker image rm -f hilink-reconnect

uninstall-failover:
	while IFS='' read -r line || [ -n "$$line" ]; do \
		[ -n "$$line" ] || continue; \
		echo "$$line" | grep -qv "^#" || continue; \
		\
		GATEWAY=$$(echo "$$line" | awk -F '\t' '{print $$1}'); \
		INTERFACE=$$(echo "$$line" | awk -F '\t' '{print $$2}'); \
		TABLE=$$(echo "$$line" | awk -F '\t' '{print $$4}'); \
		\
		if [ -n "$$GATEWAY" -a -n "$$INTERFACE" ]; then \
			ip route del default via "$$GATEWAY" dev "$$INTERFACE" metric "$$TABLE"; \
		elif [ -n "$$GATEWAY" ]; then \
			ip route del default via "$$GATEWAY" metric "$$TABLE"; \
		elif [ -n "$$INTERFACE" ]; then \
			ip route del default dev "$$INTERFACE" metric "$$TABLE"; \
		fi; \
	done < /opt/multiproxy/instances

	systemctl disable --now multiproxy-simplefailover
	rm /etc/systemd/system/multiproxy-simplefailover.service

	if [ -d /lib/dhcpcd ]; then \
		rm /lib/dhcpcd/dhcpcd-hooks/99-multiproxy-fallback; \
	fi

	if systemctl is-enabled -q dhcpcd 2>/dev/null; then \
		systemctl restart dhcpcd; \
	fi

uninstall-all: uninstall uninstall-failover
