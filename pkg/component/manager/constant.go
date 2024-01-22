package manager

const (
	BackendEntrypointKey = "backend_entrypoint.sh"
	EntrypointKey        = "entrypoint.sh"

	BackendEntrypointValue = `/pulsar-manager/pulsar-manager/bin/pulsar-manager \
      --sync.cluster.interval=60000 \
      --spring.datasource.initialization-mode=never \
      --spring.datasource.driver-class-name=org.postgresql.Driver \
      --spring.datasource.url=jdbc:postgresql://127.0.0.1:5432/pulsar_manager \
      --spring.datasource.username=pulsar \
      --spring.datasource.password=pulsar \
      --pagehelper.helperDialect=postgresql \
      --bookie.host="http://%s:8000" \
      --bookie.enable=true \
      --redirect.scheme=http \
      --redirect.port=80 \
      --redirect.host=admin.test001.test.pulsar.example.local \
      --default.environment.name=%s \
      --default.environment.service_url=http://%s:8080 \
      --tls.enabled=false \
      --pulsar.peek.message=true `

	EntrypointValue = `
apk add --update openssl && rm -rf /var/cache/apk/*;
mkdir conf;
echo 'Starting PostGreSQL Server';
addgroup pulsar;
adduser --disabled-password --ingroup pulsar pulsar;
mkdir -p /run/postgresql;
chown -R pulsar:pulsar /run/postgresql/;
chown -R pulsar:pulsar /data;
chown pulsar:pulsar /pulsar-manager/init_db.sql;
chmod 750 /data;
su - pulsar -s /bin/sh /pulsar-manager/startup.sh;
echo 'Starting Pulsar Manager Front end';
nginx;
echo 'Starting Pulsar Manager Back end';
export JAVA_OPTS=${JAVA_OPTS};
chmod +x /pulsar-manager/pulsar-backend-entrypoint.sh;
/pulsar-manager/pulsar-backend-entrypoint.sh;`
)
