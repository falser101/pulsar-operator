package proxy

import (
	"fmt"

	"github.com/falser101/pulsar-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeConfigMap(c *v1alpha1.Pulsar) *v1.ConfigMap {
	var configData = map[string]string{
		"PULSAR_GC":                            PulsarGC,
		"PULSAR_MEM":                           PulsarMem,
		"authenticationProviders":              "org.apache.pulsar.broker.authentication.AuthenticationProviderToken",
		"authorizationEnabled":                 "false",
		"brokerClientAuthenticationParameters": "file:///pulsar/tokens/proxy/token",
		"brokerClientAuthenticationPlugin":     "org.apache.pulsar.client.impl.auth.AuthenticationToken",
		"brokerServiceURL":                     fmt.Sprintf("pulsar://%s:6650", MakeServiceName(c)),
		"brokerWebServiceURL":                  fmt.Sprintf("http://%s:8080", MakeServiceName(c)),
		"clusterName":                          "pulsar-cluster",
		"forwardAuthorizationCredentials":      "true",
		"httpNumThreads":                       "8",
		"servicePort":                          "6650",
		"statusFilePath":                       "/pulsar/status",
		"superUserRoles":                       "broker-admin,proxy-admin,ws-admin,admin,console-admin",
		"tokenPublicKey":                       "file:///pulsar/keys/token/public.key",
		"webServicePort":                       "8080",
		"log4j2.yaml": `
Configuration:
status: INFO
monitorInterval: 30
name: pulsar
packages: io.prometheus.client.log4j2

Properties:
  Property:
    - name: "pulsar.log.dir"
      value: "logs"
    - name: "pulsar.log.file"
      value: "pulsar.log"
    - name: "pulsar.log.appender"
      value: "RoutingAppender"
    - name: "pulsar.log.root.level"
      value: "info"
    - name: "pulsar.log.level"
      value: "info"
    - name: "pulsar.routing.appender.default"
      value: "Console"

# Example: logger-filter script
Scripts:
  ScriptFile:
    name: filter.js
    language: JavaScript
    path: ./conf/log4j2-scripts/filter.js
    charset: UTF-8

Appenders:

  # Console
  Console:
    name: Console
    target: SYSTEM_OUT
    PatternLayout:
      Pattern: "%d{HH:mm:ss.SSS} [%t] %-5level %logger{36} - %msg%n"

  # Rolling file appender configuration
  RollingFile:
    name: RollingFile
    fileName: "${sys:pulsar.log.dir}/${sys:pulsar.log.file}"
    filePattern: "${sys:pulsar.log.dir}/${sys:pulsar.log.file}-%d{MM-dd-yyyy}-%i.log.gz"
    immediateFlush: false
    PatternLayout:
      Pattern: "%d{HH:mm:ss.SSS} [%t] %-5level %logger{36} - %msg%n"
    Policies:
      TimeBasedTriggeringPolicy:
        interval: 1
        modulate: true
      SizeBasedTriggeringPolicy:
        size: 1 GB
    # Delete file older than 30days
    DefaultRolloverStrategy:
        Delete:
          basePath: ${sys:pulsar.log.dir}
          maxDepth: 2
          IfFileName:
            glob: "*/${sys:pulsar.log.file}*log.gz"
          IfLastModified:
            age: 30d

  Prometheus:
    name: Prometheus

  # Routing
  Routing:
    name: RoutingAppender
    Routes:
      pattern: "$${ctx:function}"
      Route:
        -
          Routing:
            name: InstanceRoutingAppender
            Routes:
              pattern: "$${ctx:instance}"
              Route:
                -
                  RollingFile:
                    name: "Rolling-${ctx:function}"
                    fileName : "${sys:pulsar.log.dir}/functions/${ctx:function}/${ctx:functionname}-${ctx:instance}.log"
                    filePattern : "${sys:pulsar.log.dir}/functions/${sys:pulsar.log.file}-${ctx:instance}-%d{MM-dd-yyyy}-%i.log.gz"
                    PatternLayout:
                      Pattern: "%d{ABSOLUTE} %level{length=5} [%thread] [instance: %X{instance}] %logger{1} - %msg%n"
                    Policies:
                      TimeBasedTriggeringPolicy:
                        interval: 1
                        modulate: true
                      SizeBasedTriggeringPolicy:
                        size: "20MB"
                      # Trigger every day at midnight that also scan
                      # roll-over strategy that deletes older file
                      CronTriggeringPolicy:
                        schedule: "0 0 0 * * ?"
                    # Delete file older than 30days
                    DefaultRolloverStrategy:
                        Delete:
                          basePath: ${sys:pulsar.log.dir}
                          maxDepth: 2
                          IfFileName:
                            glob: "*/${sys:pulsar.log.file}*log.gz"
                          IfLastModified:
                            age: 30d
                - ref: "${sys:pulsar.routing.appender.default}"
                  key: "${ctx:function}"
        - ref: "${sys:pulsar.routing.appender.default}"
          key: "${ctx:function}"

Loggers:

  # Default root logger configuration
  Root:
    level: "${sys:pulsar.log.root.level}"
    additivity: true
    AppenderRef:
      - ref: "${sys:pulsar.log.appender}"
        level: "${sys:pulsar.log.level}"
      - ref: Prometheus
        level: info

  Logger:
    - name: org.apache.bookkeeper.bookie.BookieShell
      level: info
      additivity: false
      AppenderRef:
        - ref: Console

    - name: verbose
      level: info
      additivity: false
      AppenderRef:
        - ref: Console

  # Logger to inject filter script
#     - name: org.apache.bookkeeper.mledger.impl.ManagedLedgerImpl
#       level: debug
#       additivity: false
#       AppenderRef:
#         ref: "${sys:pulsar.log.appender}"
#         ScriptFilter:
#           onMatch: ACCEPT
#           onMisMatch: DENY
#           ScriptRef:
#             ref: filter.js
`,
	}
	if c.Spec.Broker.Authentication.Enabled {
		configData["authenticationEnabled"] = "true"
	}
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MakeConfigMapName(c),
			Namespace: c.Namespace,
		},
		Data: configData,
	}
}

func MakeConfigMapName(c *v1alpha1.Pulsar) string {
	return fmt.Sprintf("%s-proxy-configmap", c.GetName())
}
