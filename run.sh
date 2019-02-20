#!/bin/sh

if [ "$DOCKER" = true ]; then
	ETC_PATH=/etc
	ETC_GEN_PATH=/etc
else
	ETC_PATH=/opt/multiproxy
	ETC_GEN_PATH=/tmp/multiproxy

	rm -rf "$ETC_GEN_PATH"
	mkdir "$ETC_GEN_PATH"
fi

mkdir "$ETC_GEN_PATH/s6.d" "$ETC_GEN_PATH/s6.d/.s6-svscan" "$ETC_GEN_PATH/instance_config"

if [ "$DOCKER" = true ]; then
	cat <<EOF > "$ETC_GEN_PATH/s6.d/.s6-svscan/finish"
#!/bin/sh

exit 0
EOF
	chmod +x "$ETC_GEN_PATH/s6.d/.s6-svscan/finish"
else
	cat <<EOF > "$ETC_GEN_PATH/s6.d/.s6-svscan/finish"
#!/bin/sh

rm -rf "$ETC_GEN_PATH"
exit 0
EOF
	chmod +x "$ETC_GEN_PATH/s6.d/.s6-svscan/finish"
fi

if [ "$DOCKER" = true ]; then
	mkdir "$ETC_GEN_PATH/s6.d/syslogd"
	cat <<EOF > "$ETC_GEN_PATH/s6.d/syslogd/run"
#!/bin/sh

exec syslogd -nO -
EOF
	chmod +x "$ETC_GEN_PATH/s6.d/syslogd/run"
fi

while IFS='' read -r line || [ -n "$line" ]; do
	[ -n "$line" ] || continue
	echo "$line" | grep -qv "^#" || continue

	MARK=$(echo "$line" | awk -F '\t' '{print $3}')
	PUID=$(echo "$line" | awk -F '\t' '{print $5}')
	PGID=$(echo "$line" | awk -F '\t' '{print $5}')
	PORT=$(echo "$line" | awk -F '\t' '{print $6}')
	USER=$(echo "$line" | awk -F '\t' '{print $7}')
	PASSWORD=$(echo "$line" | awk -F '\t' '{print $8}')

	cp "$ETC_PATH/3proxy.cfg" "$ETC_GEN_PATH/instance_config/$MARK.cfg"
	sed -i "s/@PORT@/$PORT/g" "$ETC_GEN_PATH/instance_config/$MARK.cfg"
	sed -i "s/@USER@/$USER/g" "$ETC_GEN_PATH/instance_config/$MARK.cfg"
	sed -i "s/@PASSWORD@/$PASSWORD/g" "$ETC_GEN_PATH/instance_config/$MARK.cfg"

	mkdir "$ETC_GEN_PATH/s6.d/$MARK"
	cat <<EOF > "$ETC_GEN_PATH/s6.d/$MARK/run"
#!/bin/sh

$([ "$DOCKER" = true ] && echo "s6-svwait \"$ETC_GEN_PATH/s6.d/syslogd\"")
exec su-exec "$PUID:$PGID" 3proxy "$ETC_GEN_PATH/instance_config/$MARK.cfg" 
EOF
	chmod +x "$ETC_GEN_PATH/s6.d/$MARK/run"
done < "$ETC_PATH/instances"

exec s6-svscan "$ETC_GEN_PATH/s6.d"
