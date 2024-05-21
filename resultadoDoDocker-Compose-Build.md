 docker-compose up --build
WARN[0000] /Users/sandroluisschmidtke/Pos-Go/temperaturaComMetricas/docker-compose.yml: `version` is obsolete
[+] Building 1.4s (23/23) FINISHED       docker:desktop-linux
 => [serviceb internal] load build definition from Dock  0.0s
 => => transferring dockerfile: 332B                     0.0s
 => [servicea internal] load metadata for docker.io/lib  1.1s
 => [serviceb internal] load .dockerignore               0.0s
 => => transferring context: 2B                          0.0s
 => [servicea 1/9] FROM docker.io/library/golang:1.21-a  0.0s
 => [serviceb internal] load build context               0.0s
 => => transferring context: 11.45kB                     0.0s
 => CACHED [servicea 2/9] RUN apk add --no-cache git     0.0s
 => CACHED [servicea 3/9] WORKDIR /app                   0.0s
 => CACHED [serviceb 4/8] COPY go.mod ./                 0.0s
 => CACHED [serviceb 5/8] COPY go.sum ./                 0.0s
 => CACHED [serviceb 6/8] RUN go mod download            0.0s
 => CACHED [serviceb 7/8] COPY *.go ./                   0.0s
 => CACHED [serviceb 8/8] RUN go build -o /serviceb      0.0s
 => [serviceb] exporting to image                        0.0s
 => => exporting layers                                  0.0s
 => => writing image sha256:47f5b8fd0959b5eaf6d17d6def4  0.0s
 => => naming to docker.io/library/temperaturacommetric  0.0s
 => [servicea internal] load build definition from Dock  0.0s
 => => transferring dockerfile: 398B                     0.0s
 => [servicea internal] load .dockerignore               0.0s
 => => transferring context: 2B                          0.0s
 => [servicea internal] load build context               0.0s
 => => transferring context: 12.25kB                     0.0s
 => CACHED [servicea 4/9] COPY go.mod ./                 0.0s
 => CACHED [servicea 5/9] COPY go.sum ./                 0.0s
 => CACHED [servicea 6/9] RUN go mod download            0.0s
 => CACHED [servicea 7/9] COPY*.go ./                   0.0s
 => CACHED [servicea 8/9] COPY index.html ./             0.0s
 => CACHED [servicea 9/9] RUN go build -o /servicea      0.0s
 => [servicea] exporting to image                        0.0s
 => => exporting layers                                  0.0s
 => => writing image sha256:52a45d27ffd015dc2ca9d191230  0.0s
 => => naming to docker.io/library/temperaturacommetric  0.0s
[+] Running 5/3
[+] Running 7/6eraturacommetricas_default         Created0.0s
 ✔ Network temperaturacommetricas_default         Created0.0s
 ✔ Container temperaturacommetricas-otelcol-1     Created0.1s
 ✔ Container temperaturacommetricas-zipkin-1      Created0.1s
 ✔ Container temperaturacommetricas-serviceb-1    Created0.0s
 ✔ Container temperaturacommetricas-servicea-1    Created0.0s
 ✔ Container temperaturacommetricas-prometheus-1  Created0.0s
 ✔ Container temperaturacommetricas-grafana-1     Created0.0s
Attaching to grafana-1, otelcol-1, prometheus-1, servicea-1, serviceb-1, zipkin-1
otelcol-1     | 2024-05-21T03:12:40.659Z        info    service@v0.100.0/service.go:102       Setting up own telemetry...
otelcol-1     | 2024-05-21T03:12:40.659Z        info    service@v0.100.0/telemetry.go:103     Serving metrics {"address": ":8888", "level": "Normal"}
otelcol-1     | 2024-05-21T03:12:40.660Z        info    service@v0.100.0/service.go:169       Starting otelcol-contrib...  {"Version": "0.100.0", "NumCPU": 10}
otelcol-1     | 2024-05-21T03:12:40.660Z        info    extensions/extensions.go:34   Starting extensions...
otelcol-1     | 2024-05-21T03:12:40.660Z        warn    internal@v0.100.0/warning.go:42       Using the 0.0.0.0 address exposes this server to every network interface, which may facilitate Denial of Service attacks. Enable the feature gate to change the default and remove this warning.  {"kind": "receiver", "name": "otlp", "data_type": "traces", "documentation": "<https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/security-best-practices.md#safeguards-against-denial-of-service-attacks>", "feature gate ID": "component.UseLocalHostAsDefaultHost"}
otelcol-1     | 2024-05-21T03:12:40.660Z        info    otlpreceiver@v0.100.0/otlp.go:102     Starting GRPC server    {"kind": "receiver", "name": "otlp", "data_type": "traces", "endpoint": "0.0.0.0:4317"}
otelcol-1     | 2024-05-21T03:12:40.660Z        warn    internal@v0.100.0/warning.go:42       Using the 0.0.0.0 address exposes this server to every network interface, which may facilitate Denial of Service attacks. Enable the feature gate to change the default and remove this warning.  {"kind": "receiver", "name": "otlp", "data_type": "traces", "documentation": "<https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/security-best-practices.md#safeguards-against-denial-of-service-attacks>", "feature gate ID": "component.UseLocalHostAsDefaultHost"}
otelcol-1     | 2024-05-21T03:12:40.660Z        info    otlpreceiver@v0.100.0/otlp.go:152     Starting HTTP server    {"kind": "receiver", "name": "otlp", "data_type": "traces", "endpoint": "0.0.0.0:4318"}
otelcol-1     | 2024-05-21T03:12:40.661Z        info    service@v0.100.0/service.go:195       Everything is ready. Begin running and processing data.
otelcol-1     | 2024-05-21T03:12:40.661Z        warn    localhostgate/featuregate.go:63       The default endpoints for all servers in components will change to use localhost instead of 0.0.0.0 in a future version. Use the feature gate to preview the new default. {"feature gate ID": "component.UseLocalHostAsDefaultHost"}
serviceb-1    | 2024/05/21 03:12:40 Service B is running on port 8081
servicea-1    | 2024/05/21 03:12:40 Service A is running on port 8080
prometheus-1  | ts=2024-05-21T03:12:41.054Z caller=main.go:573 level=info msg="No time or size retention was set so using the default time retention" duration=15d
prometheus-1  | ts=2024-05-21T03:12:41.054Z caller=main.go:617 level=info msg="Starting Prometheus Server" mode=server version="(version=2.52.0, branch=HEAD, revision=879d80922a227c37df502e7315fad8ceb10a986d)"
prometheus-1  | ts=2024-05-21T03:12:41.054Z caller=main.go:622 level=info build_context="(go=go1.22.3, platform=linux/arm64, user=root@1b4f4c206e41, date=20240508-21:59:01, tags=netgo,builtinassets,stringlabels)"
prometheus-1  | ts=2024-05-21T03:12:41.054Z caller=main.go:623 level=info host_details="(Linux 6.6.26-linuxkit #1 SMP Sat Apr 27 04:13:19 UTC 2024 aarch64 fd102547b2dd (none))"
prometheus-1  | ts=2024-05-21T03:12:41.054Z caller=main.go:624 level=info fd_limits="(soft=1048576, hard=1048576)"
prometheus-1  | ts=2024-05-21T03:12:41.054Z caller=main.go:625 level=info vm_limits="(soft=unlimited, hard=unlimited)"
prometheus-1  | ts=2024-05-21T03:12:41.058Z caller=web.go:568 level=info component=web msg="Start listening for connections" address=0.0.0.0:9090
prometheus-1  | ts=2024-05-21T03:12:41.059Z caller=main.go:1129 level=info msg="Starting TSDB ..."
prometheus-1  | ts=2024-05-21T03:12:41.060Z caller=tls_config.go:313 level=info component=web msg="Listening on" address=[::]:9090
prometheus-1  | ts=2024-05-21T03:12:41.060Z caller=tls_config.go:316 level=info component=web msg="TLS is disabled." http2=false address=[::]:9090
prometheus-1  | ts=2024-05-21T03:12:41.061Z caller=head.go:616 level=info component=tsdb msg="Replaying on-disk memory mappable chunks if any"
prometheus-1  | ts=2024-05-21T03:12:41.061Z caller=head.go:703 level=info component=tsdb msg="On-disk memory mappable chunks replay completed" duration=8.75µs
prometheus-1  | ts=2024-05-21T03:12:41.061Z caller=head.go:711 level=info component=tsdb msg="Replaying WAL, this may take a while"
prometheus-1  | ts=2024-05-21T03:12:41.061Z caller=head.go:783 level=info component=tsdb msg="WAL segment loaded" segment=0 maxSegment=0
prometheus-1  | ts=2024-05-21T03:12:41.061Z caller=head.go:820 level=info component=tsdb msg="WAL replay completed" checkpoint_replay_duration=21.25µs wal_replay_duration=415.417µs wbl_replay_duration=125ns chunk_snapshot_load_duration=0s mmap_chunk_replay_duration=8.75µs total_replay_duration=457.583µs
prometheus-1  | ts=2024-05-21T03:12:41.063Z caller=main.go:1150 level=info fs_type=EXT4_SUPER_MAGIC
prometheus-1  | ts=2024-05-21T03:12:41.063Z caller=main.go:1153 level=info msg="TSDB started"
prometheus-1  | ts=2024-05-21T03:12:41.063Z caller=main.go:1335 level=info msg="Loading configuration file" filename=/etc/prometheus/prometheus.yml
prometheus-1  | ts=2024-05-21T03:12:41.066Z caller=main.go:1372 level=info msg="Completed loading of configuration file" filename=/etc/prometheus/prometheus.yml totalDuration=2.53375ms db_storage=1.833µs remote_storage=1.542µs web_handler=792ns query_engine=1.958µs scrape=292.542µs scrape_sd=41.917µs notify=1µs notify_sd=584ns rules=1.583µs tracing=10.125µs
prometheus-1  | ts=2024-05-21T03:12:41.066Z caller=main.go:1114 level=info msg="Server is ready to receive web requests."
prometheus-1  | ts=2024-05-21T03:12:41.066Z caller=manager.go:163 level=info component="rule manager" msg="Starting rule manager..."
grafana-1     | logger=settings t=2024-05-21T03:12:41.309184261Z level=info msg="Starting Grafana" version=10.4.3 commit=0bfd547800e6eb79dc98e55844ba28194b3df002 branch=v10.4.x compiled=2024-05-21T03:12:41Z
grafana-1     | logger=settings t=2024-05-21T03:12:41.309382094Z level=info msg="Config loaded from" file=/usr/share/grafana/conf/defaults.ini
grafana-1     | logger=settings t=2024-05-21T03:12:41.309392678Z level=info msg="Config loaded from" file=/etc/grafana/grafana.ini
grafana-1     | logger=settings t=2024-05-21T03:12:41.309394761Z level=info msg="Config overridden from command line" arg="default.paths.data=/var/lib/grafana"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309396428Z level=info msg="Config overridden from command line" arg="default.paths.logs=/var/log/grafana"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309397928Z level=info msg="Config overridden from command line" arg="default.paths.plugins=/var/lib/grafana/plugins"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309399344Z level=info msg="Config overridden from command line" arg="default.paths.provisioning=/etc/grafana/provisioning"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309400844Z level=info msg="Config overridden from command line" arg="default.log.mode=console"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309402428Z level=info msg="Config overridden from Environment variable" var="GF_PATHS_DATA=/var/lib/grafana"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309404136Z level=info msg="Config overridden from Environment variable" var="GF_PATHS_LOGS=/var/log/grafana"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309405553Z level=info msg="Config overridden from Environment variable" var="GF_PATHS_PLUGINS=/var/lib/grafana/plugins"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309407053Z level=info msg="Config overridden from Environment variable" var="GF_PATHS_PROVISIONING=/etc/grafana/provisioning"
grafana-1     | logger=settings t=2024-05-21T03:12:41.309408553Z level=info msg=Target target=[all]
grafana-1     | logger=settings t=2024-05-21T03:12:41.309413886Z level=info msg="Path Home" path=/usr/share/grafana
grafana-1     | logger=settings t=2024-05-21T03:12:41.309415386Z level=info msg="Path Data" path=/var/lib/grafana
grafana-1     | logger=settings t=2024-05-21T03:12:41.309417053Z level=info msg="Path Logs" path=/var/log/grafana
grafana-1     | logger=settings t=2024-05-21T03:12:41.309418761Z level=info msg="Path Plugins" path=/var/lib/grafana/plugins
grafana-1     | logger=settings t=2024-05-21T03:12:41.309420303Z level=info msg="Path Provisioning" path=/etc/grafana/provisioning
grafana-1     | logger=settings t=2024-05-21T03:12:41.309421719Z level=info msg="App mode production"
grafana-1     | logger=sqlstore t=2024-05-21T03:12:41.309670053Z level=info msg="Connecting to DB" dbtype=sqlite3
grafana-1     | logger=migrator t=2024-05-21T03:12:41.310217261Z level=info msg="Starting DB migrations"
grafana-1     | logger=migrator t=2024-05-21T03:12:41.323212219Z level=info msg="migrations completed" performed=0 skipped=551 duration=313.417µs
grafana-1     | logger=secrets t=2024-05-21T03:12:41.324372553Z level=info msg="Envelope encryption state" enabled=true currentprovider=secretKey.v1
grafana-1     | logger=plugin.store t=2024-05-21T03:12:41.338408428Z level=info msg="Loading plugins..."
grafana-1     | logger=local.finder t=2024-05-21T03:12:41.369880594Z level=warn msg="Skipping finding plugins as directory does not exist" path=/usr/share/grafana/plugins-bundled
grafana-1     | logger=plugin.store t=2024-05-21T03:12:41.369904594Z level=info msg="Plugins loaded" count=55 duration=31.496584ms
grafana-1     | logger=query_data t=2024-05-21T03:12:41.371973261Z level=info msg="Query Service initialization"
grafana-1     | logger=live.push_http t=2024-05-21T03:12:41.374302136Z level=info msg="Live Push Gateway initialization"
grafana-1     | logger=ngalert.migration t=2024-05-21T03:12:41.389415469Z level=info msg=Starting
grafana-1     | logger=ngalert.state.manager t=2024-05-21T03:12:41.398891386Z level=info msg="Running in alternative execution of Error/NoData mode"
grafana-1     | logger=infra.usagestats.collector t=2024-05-21T03:12:41.400362386Z level=info msg="registering usage stat providers" usageStatsProvidersLen=2
grafana-1     | logger=provisioning.alerting t=2024-05-21T03:12:41.400559469Z level=info msg="starting to provision alerting"
grafana-1     | logger=provisioning.alerting t=2024-05-21T03:12:41.400573178Z level=info msg="finished to provision alerting"
grafana-1     | logger=ngalert.state.manager t=2024-05-21T03:12:41.401606928Z level=info msg="Warming state cache for startup"
grafana-1     | logger=ngalert.state.manager t=2024-05-21T03:12:41.402238761Z level=info msg="State cache has been initialized" states=0 duration=630.625µs
grafana-1     | logger=grafanaStorageLogger t=2024-05-21T03:12:41.402675136Z level=info msg="Storage starting"
grafana-1     | logger=ngalert.multiorg.alertmanager t=2024-05-21T03:12:41.403164969Z level=info msg="Starting MultiOrg Alertmanager"
grafana-1     | logger=http.server t=2024-05-21T03:12:41.403339594Z level=info msg="HTTP Server Listen" address=[::]:3000 protocol=http subUrl= socket=
grafana-1     | logger=ngalert.scheduler t=2024-05-21T03:12:41.403567303Z level=info msg="Starting scheduler" tickInterval=10s maxAttempts=1
grafana-1     | logger=ticker t=2024-05-21T03:12:41.403629178Z level=info msg=starting first_tick=2024-05-21T03:12:50Z
grafana-1     | logger=provisioning.dashboard t=2024-05-21T03:12:41.439916886Z level=info msg="starting to provision dashboards"
grafana-1     | logger=provisioning.dashboard t=2024-05-21T03:12:41.439941469Z level=info msg="finished to provision dashboards"
zipkin-1      |
zipkin-1      |                   oo
zipkin-1      |                  oooo
zipkin-1      |                 oooooo
zipkin-1      |                oooooooo
zipkin-1      |               oooooooooo
zipkin-1      |              oooooooooooo
zipkin-1      |            ooooooo  ooooooo
zipkin-1      |           oooooo     ooooooo
zipkin-1      |          oooooo       ooooooo
zipkin-1      |         oooooo   o  o   oooooo
zipkin-1      |        oooooo   oo  oo   oooooo
zipkin-1      |      ooooooo  oooo  oooo  ooooooo
zipkin-1      |     oooooo   ooooo  ooooo  ooooooo
zipkin-1      |    oooooo   oooooo  oooooo  ooooooo
zipkin-1      |   oooooooo      oo  oo      oooooooo
zipkin-1      |   ooooooooooooo oo  oo ooooooooooooo
zipkin-1      |       oooooooooooo  oooooooooooo
zipkin-1      |           oooooooo  oooooooo
zipkin-1      |               oooo  oooo
zipkin-1      |
zipkin-1      |      ________ ____  _______   _
zipkin-1      |     |__/_ *|  _ \| |/ /* *| \ | |
zipkin-1      |       / / | || |*) | ' / | ||  \| |
zipkin-1      |      / /_| ||__/| . \ | || |\  |
zipkin-1      |     |*__*|__*|*|   |*|\_\___|*| \_|
zipkin-1      |
zipkin-1      | :: version 3.3.0 :: commit dfd8ee2 ::
zipkin-1      |
grafana-1     | logger=grafana.update.checker t=2024-05-21T03:12:41.695469553Z level=info msg="Update check succeeded" duration=293.184125ms
grafana-1     | logger=plugins.update.checker t=2024-05-21T03:12:41.710317928Z level=info msg="Update check succeeded" duration=307.325875ms
grafana-1     | logger=grafana-apiserver t=2024-05-21T03:12:41.753331428Z level=info msg="Adding GroupVersion playlist.grafana.app v0alpha1 to ResourceManager"
grafana-1     | logger=grafana-apiserver t=2024-05-21T03:12:41.754791136Z level=info msg="Adding GroupVersion featuretoggle.grafana.app v0alpha1 to ResourceManager"
zipkin-1      | 2024-05-21T03:12:42.392Z  INFO [/] 1 --- [oss-http-*:9411] c.l.a.s.Server                           : Serving HTTP at /[0:0:0:0:0:0:0:0%0]:9411 - <http://127.0.0.1:9411/>

## Execução dos serviços e resultado

curl -X POST <http://localhost:8080/cep> -H "Content-Type: application/json" -d '{"cep": "70766060"}'

{"city":"Brasília","temp_C":22,"temp_F":71.6,"temp_K":295}

## docker-compose down

WARN[0000] /Users/sandroluisschmidtke/Pos-Go/temperaturaComMetricas/docker-compose.yml: `version` is obsolete
[+] Running 7/0
 ✔ Container temperaturacommetricas-zipkin-1      Removed0.0s
 ✔ Container temperaturacommetricas-grafana-1     Removed0.0s
 ✔ Container temperaturacommetricas-prometheus-1  Removed0.0s
 ✔ Container temperaturacommetricas-servicea-1    Removed0.0s
 ✔ Container temperaturacommetricas-serviceb-1    Removed0.0s
 ✔ Container temperaturacommetricas-otelcol-1     Removed0.0s
 ✔ Network temperaturacommetricas_default         Removed0.0s
