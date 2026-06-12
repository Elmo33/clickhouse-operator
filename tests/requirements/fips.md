# QA-SRS ClickHouse Operator FIPS 140-3
# Software Requirements Specification

(c) 2026 Altinity Inc. All Rights Reserved.

**Document status:** Confidential

**Author:** Saba Momtselidze

**Date:** May 29, 2026

## Table of Contents

* 1 [Introduction](#introduction)
* 2 [Configuration Requirements](#configuration-requirements)
    * 2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Config.ExternalTLS](#rqsrs-026clickhouseoperatorfipsconfigexternaltls)
* 3 [Build Verification](#build-verification)
    * 3.1 [Shipped Binaries](#shipped-binaries)
        * 3.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries](#rqsrs026clickhouseoperatorfipsbuildshippedbinaries)
            * 3.1.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries.GOFIPS140](#rqsrs026clickhouseoperatorfipsbuildshippedbinariesgofips140)
            * 3.1.1.2 [RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries.FIPSIdentity](#rqsrs026clickhouseoperatorfipsbuildshippedbinariesfipsidentity)
            * 3.1.1.3 [RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries.FIPSVersion](#rqsrs026clickhouseoperatorfipsbuildshippedbinariesfipsversion)
            * 3.1.1.4 [RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries.FIPSEnabled](#rqsrs026clickhouseoperatorfipsbuildshippedbinariesfipsenabled)
            * 3.1.1.5 [RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries.StartupBanner](#rqsrs026clickhouseoperatorfipsbuildshippedbinariesstartupbanner)
* 4 [GODEBUG Strict Mode Smoke Test](#godebug-strict-mode-smoke-test)
    * 4.1 [RQ.SRS-026.ClickHouseOperator.FIPS.GODEBUG.StrictMode](#rqsrs-026clickhouseoperatorfipsgodebugstrictmode)
* 5 [FIPS 140-3 Valid TLS Cipher Suites](#fips-140-3-valid-tls-cipher-suites)
    * 5.1 [Approved TLS Cipher Suites](#approved-tls-cipher-suites)
        * 5.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.TLS.ApprovedCiphers](#rqsrs-026clickhouseoperatorfipstlsapprovedciphers)
    * 5.2 [Rejected Cipher Suites and Protocols](#rejected-cipher-suites-and-protocols)
        * 5.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.TLS.RejectedCiphers](#rqsrs-026clickhouseoperatorfipstlsrejectedciphers)
* 6 [ClickHouse Server and Keeper FIPS Configurations](#clickhouse-server-and-keeper-fips-configurations)
    * 6.1 [ClickHouse Server](#clickhouse-server)
        * 6.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.FIPSConfig](#rqsrs026clickhouseoperatorfipsdataplanechfipsconfig)
        * 6.1.2 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHIDeploy](#rqsrs026clickhouseoperatorfipsdataplanechideploy)
        * 6.1.3 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.NoPlainHTTP](#rqsrs026clickhouseoperatorfipsdataplanechnoplainhttp)
        * 6.1.4 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.NoPlainNative](#rqsrs026clickhouseoperatorfipsdataplanechnoplainnative)
        * 6.1.5 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.NoUnexpectedPorts](#rqsrs026clickhouseoperatorfipsdataplanechnounexpectedports)
        * 6.1.6 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.InternodeTLS](#rqsrs026clickhouseoperatorfipsdataplanechinternodetls)
        * 6.1.7 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.ScaleUp](#rqsrs026clickhouseoperatorfipsdataplanechscaleup)
        * 6.1.8 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.ScaleDown](#rqsrs026clickhouseoperatorfipsdataplanechscaledown)
        * 6.1.9 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.ConfigUpdate](#rqsrs026clickhouseoperatorfipsdataplanechconfigupdate)
    * 6.2 [ClickHouse Keeper](#clickhouse-keeper)
        * 6.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.FIPSConfig](#rqsrs026clickhouseoperatorfipsdataplanechkfipsconfig)
        * 6.2.2 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHKDeploy](#rqsrs026clickhouseoperatorfipsdataplanechkdeploy)
        * 6.2.3 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.NoPlainClientPort](#rqsrs026clickhouseoperatorfipsdataplanechknoplainclientport)
        * 6.2.4 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.NoUnexpectedPorts](#rqsrs026clickhouseoperatorfipsdataplanechknounexpectedports)
        * 6.2.5 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.RaftTLS](#rqsrs026clickhouseoperatorfipsdataplanechkrafttls)
        * 6.2.6 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.ScaleUp](#rqsrs026clickhouseoperatorfipsdataplanechkscaleup)
        * 6.2.7 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.ScaleDown](#rqsrs026clickhouseoperatorfipsdataplanechkscaledown)
        * 6.2.8 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CHK.ConfigUpdate](#rqsrs026clickhouseoperatorfipsdataplanechkconfigupdate)
    * 6.3 [ClickHouse Backup Sidecar](#clickhouse-backup-sidecar)
        * 6.3.0 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.CH.VersionString](#rqsrs026clickhouseoperatorfipsdataplanechversionstring)
        * 6.3.1 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.FIPSBinary](#rqsrs026clickhouseoperatorfipsdataplanebackupfipsbinary)
        * 6.3.2 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.GOFIPS140](#rqsrs026clickhouseoperatorfipsdataplanebackupgofips140)
        * 6.3.3 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.OnlyTLSPorts](#rqsrs026clickhouseoperatorfipsdataplanebackuponlytlsports)
        * 6.3.4 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.HTTPSAPI](#rqsrs026clickhouseoperatorfipsdataplanebackuphttpsapi)
        * 6.3.5 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.ClickHouseOverTLS](#rqsrs026clickhouseoperatorfipsdataplanebackupclickhouseovertls)
        * 6.3.6 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.RestoreRoundTrip](#rqsrs026clickhouseoperatorfipsdataplanebackuprestoreroundtrip)
        * 6.3.7 [RQ.SRS-026.ClickHouseOperator.FIPS.DataPlane.Backup.RemoteUploadTLS](#rqsrs026clickhouseoperatorfipsdataplanebackupremoteuploadtls)
* 7 [FIPS Enforcement Mode](#fips-enforcement-mode)
    * 7.1 [Security Coercion](#security-coercion)
        * 7.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.SecurityCoercion](#rqsrs026clickhouseoperatorfipsenforcedsecuritycoercion)
        * 7.1.2 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectInsecureKubeconfig](#rqsrs026clickhouseoperatorfipsenforcedrejectinsecurekubeconfig)
        * 7.1.3 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectVerifyNoneCHI](#rqsrs026clickhouseoperatorfipsenforcedrejectverifynonechi)
        * 7.1.4 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectVerifyNoneZK](#rqsrs026clickhouseoperatorfipsenforcedrejectverifynonezk)
        * 7.1.5 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectInvalidMinVersion](#rqsrs026clickhouseoperatorfipsenforcedrejectinvalidminversion)
        * 7.1.6 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectExternalZookeeper](#rqsrs026clickhouseoperatorfipsenforcedrejectexternalzookeeper)
        * 7.1.7 [RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectCHKBypass](#rqsrs026clickhouseoperatorfipsenforcedrejectchkbypass)
    * 7.2 [Image Policy](#image-policy)
        * 7.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.RejectCHI](#rqsrs026clickhouseoperatorfipsimagesrequiredrejectchi)
        * 7.2.2 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.AcceptCHI](#rqsrs026clickhouseoperatorfipsimagesrequiredacceptchi)
        * 7.2.3 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.RejectCHK](#rqsrs026clickhouseoperatorfipsimagesrequiredrejectchk)
        * 7.2.4 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.RuntimeVersion](#rqsrs026clickhouseoperatorfipsimagesrequiredruntimeversion)
        * 7.2.5 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.Permissive](#rqsrs026clickhouseoperatorfipsimagespermissive)
        * 7.2.6 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.ShortCircuit](#rqsrs026clickhouseoperatorfipsimagesrequiredshortcircuit)
    * 7.3 [Image Tag Detection](#image-tag-detection)
        * 7.3.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.FIPSSuffix](#rqsrs026clickhouseoperatorfipsimagestagdetectionfipssuffix)
        * 7.3.2 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.AltinityFIPS](#rqsrs026clickhouseoperatorfipsimagestagdetectionaltinityfips)
        * 7.3.3 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.DigestOnly](#rqsrs026clickhouseoperatorfipsimagestagdetectiondigestonly)
        * 7.3.4 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.RegistryPath](#rqsrs026clickhouseoperatorfipsimagestagdetectionregistrypath)
        * 7.3.5 [RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.CaseInsensitive](#rqsrs026clickhouseoperatorfipsimagestagdetectioncaseinsensitive)
* 8 [Operator External Connections](#operator-external-connections)
    * 8.1 [Operator Runtime Listener Verification](#operator-runtime-listener-verification)
        * 8.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.Listeners](#rqsrs026clickhouseoperatorfipsconnectoperatorlisteners)
    * 8.2 [Operator to Kubernetes API](#operator-to-kubernetes-api)
        * 8.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.Kubernetes](#rqsrs026clickhouseoperatorfipsconnectoperatorkubernetes)
    * 8.3 [Operator to ClickHouse Server](#operator-to-clickhouse-server)
        * 8.3.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.ClickHouse](#rqsrs026clickhouseoperatorfipsconnectoperatorclickhouse)
    * 8.4 [Operator to ZooKeeper/Keeper](#operator-to-zookeeperkeeper)
        * 8.4.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.Zookeeper](#rqsrs026clickhouseoperatorfipsconnectoperatorzookeeper)
    * 8.5 [Operator to metrics-exporter IPC](#operator-to-metrics-exporter-ipc)
        * 8.5.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.IPCSecure](#rqsrs026clickhouseoperatorfipsconnectoperatoripcsecure)
    * 8.6 [Operator Prometheus Metrics](#operator-prometheus-metrics)
        * 8.6.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Gap.OperatorMetricsTLS](#rqsrs026clickhouseoperatorfipsgapoperatormetricstls)
* 9 [Exporter External Connections](#exporter-external-connections)
    * 9.1 [Exporter to Kubernetes API](#exporter-to-kubernetes-api)
        * 9.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Exporter.Kubernetes](#rqsrs026clickhouseoperatorfipsconnectexporterkubernetes)
    * 9.2 [Exporter to ClickHouse Server](#exporter-to-clickhouse-server)
        * 9.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Exporter.ClickHouse](#rqsrs026clickhouseoperatorfipsconnectexporterclickhouse)
    * 9.3 [Exporter Prometheus Metrics](#exporter-prometheus-metrics)
        * 9.3.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Gap.ExporterMetricsTLS](#rqsrs026clickhouseoperatorfipsgapexportermetricstls)
* 10 [Integrity Check Failure](#integrity-check-failure)
    * 10.1 [Operator Integrity Tampering](#operator-integrity-tampering)
        * 10.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Integrity.OperatorMismatch](#rqsrs026clickhouseoperatorfipsintegrityoperatormismatch)
    * 10.2 [Exporter Integrity Tampering](#exporter-integrity-tampering)
        * 10.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Integrity.ExporterMismatch](#rqsrs026clickhouseoperatorfipsintegrityexportermismatch)
* 11 [CAST Failure](#cast-failure)
    * 11.1 [Operator CAST Failure](#operator-cast-failure)
        * 11.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.CAST.OperatorFail](#rqsrs026clickhouseoperatorfipscastoperatorfail)
    * 11.2 [Exporter CAST Failure](#exporter-cast-failure)
        * 11.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.CAST.ExporterFail](#rqsrs026clickhouseoperatorfipscastexporterfail)
* 12 [Synthetic TLS Cipher Validation](#synthetic-tls-cipher-validation)
    * 12.1 [Approved cipher matrix](#approved-cipher-matrix)
        * 12.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Synthetic.ApprovedCiphers](#rqsrs-026clickhouseoperatorfipssyntheticapprovedciphers)
    * 12.2 [Rejected cipher matrix](#rejected-cipher-matrix)
        * 12.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.Synthetic.RejectedCiphers](#rqsrs-026clickhouseoperatorfipssyntheticrejectedciphers)
* 13 [CI/CD Image and Policy Verification](#cicd-image-and-policy-verification)
    * 13.1 [RQ.SRS-026.ClickHouseOperator.FIPS.CICD.OperatorImageBuild](#rqsrs-026clickhouseoperatorfipscicdoperatorimagebuild)
    * 13.2 [RQ.SRS-026.ClickHouseOperator.FIPS.CICD.ExporterImageBuild](#rqsrs-026clickhouseoperatorfipscicdexporterimagebuild)
    * 13.3 [RQ.SRS-026.ClickHouseOperator.FIPS.CICD.VulnerabilityScan](#rqsrs-026clickhouseoperatorfipscicdvulnerabilityscan)
* 14 [AI Static Code Review](#ai-static-code-review)
    * 14.1 [Operator Source Review](#operator-source-review)
        * 14.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Operator.Tree](#rqsrs-026clickhouseoperatorfipsaireviewoperatortree)
        * 14.1.2 [RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Operator.SharedPkg](#rqsrs-026clickhouseoperatorfipsaireviewoperatorsharedpkg)
        * 14.1.3 [RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Operator.RegressionGate](#rqsrs-026clickhouseoperatorfipsaireviewoperatorregressiongate)
    * 14.2 [Exporter Source Review](#exporter-source-review)
        * 14.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Exporter.Tree](#rqsrs-026clickhouseoperatorfipsaireviewexportertree)
        * 14.2.2 [RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Exporter.SharedPkg](#rqsrs-026clickhouseoperatorfipsaireviewexportersharedpkg)
        * 14.2.3 [RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Exporter.RegressionGate](#rqsrs-026clickhouseoperatorfipsaireviewexporterregressiongate)
* 15 [ACVP Algorithm Validation](#acvp-algorithm-validation)
    * 15.1 [Operator ACVP Validation](#operator-acvp-validation)
        * 15.1.1 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.WrapperIntegration](#rqsrs026clickhouseoperatorfipsacvpoperatorwrapperintegration)
        * 15.1.2 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.ConfigGeneration](#rqsrs026clickhouseoperatorfipsacvpoperatorconfiggeneration)
        * 15.1.3 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.ExpectedOutputReplay](#rqsrs026clickhouseoperatorfipsacvpoperatorexpectedoutputreplay)
        * 15.1.4 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.SuiteCount](#rqsrs026clickhouseoperatorfipsacvpoperatorsuitecount)
    * 15.2 [Exporter ACVP Validation](#exporter-acvp-validation)
        * 15.2.1 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.WrapperIntegration](#rqsrs026clickhouseoperatorfipsacvpexporterwrapperintegration)
        * 15.2.2 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.ConfigGeneration](#rqsrs026clickhouseoperatorfipsacvpexporterconfiggeneration)
        * 15.2.3 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.ExpectedOutputReplay](#rqsrs026clickhouseoperatorfipsacvpexporterexpectedoutputreplay)
        * 15.2.4 [RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.SuiteCount](#rqsrs026clickhouseoperatorfipsacvpexportersuitecount)
* 16 [Terminology](#terminology)
    * 16.1 [SRS](#srs)
    * 16.2 [FIPS 140-3](#fips-140-3)
    * 16.3 [clickhouse-operator](#clickhouse-operator)
    * 16.4 [metrics-exporter](#metrics-exporter)
    * 16.5 [CHI](#chi)
    * 16.6 [CHK](#chk)
    * 16.7 [ACVP](#acvp)
    * 16.8 [CMVP](#cmvp)
    * 16.9 [CAVP](#cavp)

## Introduction

This specification describes FIPS 140-3 compatibility requirements for the
[clickhouse-operator] and [metrics-exporter] binaries built with Go FIPS support.

The goal is to verify that FIPS-enabled builds of the operator and metrics-exporter:
- Operate correctly under FIPS constraints
- Properly enforce cryptographic restrictions
- Use FIPS-compliant TLS for all inbound and outbound connections

Autotests that trace to these requirements live in
[`tests/e2e/test_operator_fips.py`](../e2e/test_operator_fips.py) and
[`tests/e2e/test_acvp.py`](../e2e/test_acvp.py).

**Boundary:** The operator and metrics-exporter run in the same pod. Internal IPC between
them is localhost HTTP and is not subject to FIPS TLS requirements. The Prometheus metrics
endpoints (operator `:9999` and metrics-exporter `:8888`) are also served over plain HTTP
and remain outside the FIPS TLS scope as a known gap.

## Configuration Requirements

Plain HTTP/TCP on any external connection is a configuration error for FIPS compliance.
TLS must be enabled for all connections to:

- Kubernetes API
- ClickHouse Server
- ZooKeeper/Keeper
- Prometheus scrape endpoints

### RQ.SRS-026.ClickHouseOperator.FIPS.Config.HTTP
version: 1.0

All external connections SHALL require TLS with FIPS-compliant settings, except for localhost IPC between the operator
and metrics-exporter and the Prometheus metrics endpoints: `:9999` and :`8888`.

## Build Verification

### RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries
version: 1.0

Each shipped pod binary — `clickhouse-operator` and `metrics-exporter` — SHALL satisfy all of the following:

* Both binaries SHALL be built with `GOFIPS140=v1.0.0` (or `certified`); `go version -m` on each binary SHALL show the `GOFIPS140` build setting when the binary is inspectable.
* Each binary SHALL identify itself as a FIPS build via `--version` output, `--fips-info`, or startup logs containing a FIPS indicator.
* Each binary SHALL report `crypto/fips140.Version()` equal to `v1.0.0` (for example via `--fips-info` or in-process inspection).
* Each binary SHALL report `crypto/fips140.Enabled()` equal to `true` when FIPS mode is active per `GODEBUG=fips140`.

Examples:
* `go version -m clickhouse-operator` contains `GOFIPS140=v1.0.0`
* `go version -m metrics-exporter` contains `GOFIPS140=v1.0.0`
* `clickhouse-operator --fips-info` reports:

  ```yaml
  fips_module:
    version: v1.0.0
    enabled: true
  ```

* `metrics-exporter --fips-info` reports:

  ```yaml
  fips_module:
    version: v1.0.0
    enabled: true
  ```

### RQ.SRS-026.ClickHouseOperator.FIPS.Build.ShippedBinaries.StartupLogs
version: 1.0

At startup, each binary SHALL emit a FIPS startup banner in logs indicating build and runtime FIPS state.

when GODEBUG=fips140=only:

```text
FIPS: chopconf.fips.enforced=true \
build.linked=true \
module.active=true \
runtime.enforced=true \
module=v1.0.0
```

## Approved TLS Cipher Suites

### RQ.SRS-026.ClickHouseOperator.FIPS.TLS.ApprovedCiphers
version: 1.0

TLS-enforced external connections for [clickhouse-operator] and [metrics-exporter]
SHALL negotiate only TLS 1.3 with the following approved cipher suites.

* TLS_AES_128_GCM_SHA256
* TLS_AES_256_GCM_SHA384
* TLS_CHACHA20_POLY1305_SHA256 (not accepted by default, needs to be specified explicitly in all openssl configs)

Any other cipher suite or protocol version SHALL be rejected by operator in a FIPS-compliant configuration.


## ClickHouse Server

#### RQ.SRS-026.ClickHouseOperator.FIPS.CH.FIPSConfig
version: 1.0

Deploying a `ClickHouseInstallation` with FIPS TLS OpenSSL settings SHALL start a FIPS-compliant ClickHouse server and client.

```yaml
  configuration:
    clusters:
      - name: default
        secure: "yes"
        insecure: "no"
        layout:
          shardsCount: 1
          replicasCount: 2
    zookeeper:
      nodes:
        - host: chk-test-030003-keeper-0-0
          port: 2281
          secure: "yes"
    settings:
      http_port: _removed_
      tcp_port: _removed_
      interserver_http_port: _removed_
      mysql_port: _removed_
      postgresql_port: _removed_
      https_port: 8443
      tcp_port_secure: 9440
      interserver_https_port: 9010
    files:
      openssl.xml: |
        <yandex>
          <openSSL>
            <server>
              <certificateFile>/etc/clickhouse-server/secrets.d/server.crt/clickhouse-certs/server.crt</certificateFile>
              <privateKeyFile>/etc/clickhouse-server/secrets.d/server.key/clickhouse-certs/server.key</privateKeyFile>
              <dhParamsFile>/etc/clickhouse-server/secrets.d/dhparam.pem/clickhouse-certs/dhparam.pem</dhParamsFile>
              <!-- Server-auth TLS only: clients validate this certificate; the server does not require client certificates (not mTLS). -->
              <verificationMode>none</verificationMode>
              <disableProtocols>sslv2,sslv3,tlsv1,tlsv1_1</disableProtocols>
              <cipherSuites>TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384</cipherSuites>
            </server>
            <client>
              <caConfig>/etc/clickhouse-server/secrets.d/ca.crt/clickhouse-certs/ca.crt</caConfig>
              <loadDefaultCAFile>false</loadDefaultCAFile>
              <verificationMode>strict</verificationMode>
              <disableProtocols>sslv2,sslv3,tlsv1,tlsv1_1</disableProtocols>
              <cipherSuites>TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384</cipherSuites>
            </client>
          </openSSL>
        </yandex>
```

The deployed ClickHouse server SHALL use only the following ports:

* HTTPS API port 8443 (instead of 8123)
* Secure native TCP port 9440 (instead of 9000)
* Interserver HTTPS port 9010 (instead of interserver HTTP port 9009)
* Backup sidecar HTTPS API port 7171 (instead of 7180), when backups are enabled

Each exposed port SHALL support TLS communication using only FIPS-compliant protocol versions and cipher suites.

#### RQ.SRS-026.ClickHouseOperator.FIPS.CH.FIPSConfig.ExternalClient
version: 1.0

External clients connecting to the ClickHouse server SHALL be able to use any enabled TLS protocol version, including TLS 1.2.

#### RQ.SRS-026.ClickHouseOperator.FIPS.CH.Rescale
version: 1.0

Adding or removing a replica from a FIPS-configured `ClickHouseInstallation` SHALL reconcile successfully and result in the expected number of running pods.

After rescaling, all replicas SHALL continue to run the FIPS ClickHouse binary and maintain the configured TLS-only OpenSSL settings.

#### RQ.SRS-026.ClickHouseOperator.FIPS.CH.ConfigUpdate
version: 1.0

Updating TLS settings on a running CHI SHALL reload ClickHouse with the new FIPS-compliant configuration.


## ClickHouse Keeper

#### RQ.SRS-026.ClickHouseOperator.FIPS.CHK.FIPSConfig
version: 1.0

Deploying a `ClickHouseKeeperInstallation` with FIPS TLS OpenSSL settings SHALL start a FIPS-compliant ClickHouse Keeper server and client.

```yaml
  configuration:
    clusters:
      - name: keeper
        secure: "yes"
        insecure: "no"
        layout:
          replicasCount: 2
    settings:
      keeper_server/log_storage_path: /var/lib/clickhouse/coordination/log
      keeper_server/snapshot_storage_path: /var/lib/clickhouse/coordination/snapshots
      keeper_server/raft_configuration/server/port: 9444
    files:
      openssl.xml: |
        <clickhouse>
          <openSSL>
              <server>
                <certificateFile>/etc/clickhouse-server/secrets.d/server.crt/clickhouse-certs/server.crt</certificateFile>
                <privateKeyFile>/etc/clickhouse-server/secrets.d/server.key/clickhouse-certs/server.key</privateKeyFile>
                <!-- Server-auth TLS only: clients validate this certificate; the server does not require client certificates (not mTLS). -->
                <verificationMode>none</verificationMode>
                <disableProtocols>sslv2,sslv3,tlsv1,tlsv1_1</disableProtocols>
                <cipherSuites>TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384</cipherSuites>
              </server>
              <client>
                <caConfig>/etc/clickhouse-server/secrets.d/ca.crt/clickhouse-certs/ca.crt</caConfig>
                <loadDefaultCAFile>false</loadDefaultCAFile>
                <verificationMode>strict</verificationMode>
                <disableProtocols>sslv2,sslv3,tlsv1,tlsv1_1</disableProtocols>
                <cipherSuites>TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384</cipherSuites>
              </client>
          </openSSL>
        </clickhouse>
```

The deployed ClickHouse Keeper cluster SHALL use only the following ports:

* Secure client port 2281 (instead of 2181)
* Secure Raft communication port 9444

Each exposed port SHALL support TLS communication using only FIPS-compliant protocol versions and cipher suites.

#### RQ.SRS-026.ClickHouseOperator.FIPS.CHK.Rescale
version: 1.0

Adding or removing a node from a FIPS-configured `ClickHouseKeeperInstallation` SHALL reconcile successfully and result 
in the expected number of running pods.

After rescaling, all Keeper nodes SHALL continue to run the FIPS ClickHouse Keeper binary and maintain the configured 
TLS-only OpenSSL settings.

#### RQ.SRS-026.ClickHouseOperator.FIPS.CHK.ConfigUpdate
version: 1.0

Updating TLS settings on a running CHK SHALL reload ClickHouse with the new FIPS-compliant configuration.


## ClickHouse Backup Sidecar

#### RQ.SRS-026.ClickHouseOperator.FIPS.Backup.FIPSBinary
version: 1.0

The `clickhouse-backup` sidecar SHALL run a FIPS-built binary.

The sidecar binary SHALL satisfy all of the following:

* `clickhouse-backup --version` contains `fips` (case-insensitive)
* When inspectable, `go version -m` reports `GOFIPS140=v1.0.0`

#### RQ.SRS-026.ClickHouseOperator.FIPS.Backup.FIPSConfig
version: 1.0

Deploying a `ClickHouseInstallation` with a FIPS-configured backup sidecar SHALL start `clickhouse-backup` with a FIPS-compliant TLS configuration.

The deployed backup sidecar SHALL only add the following listener ports to the clickhouse container:

* HTTPS API port 7171 (instead of 7180)

Each exposed port SHALL support TLS communication using only FIPS-compliant protocol versions and cipher suites.

The `clickhouse-backup` sidecar SHALL connect to ClickHouse using secure native TCP with TLS enabled.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Backup.RestoreRoundTrip
version: 1.0

Creating a backup and restoring it through the HTTPS API SHALL succeed over TLS.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Backup.RemoteUploadTLS
version: 1.0

Uploading backups to remote object storage SHALL use FIPS-compliant TLS communication.


### Security Coercion

#### RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.SecurityCoercion
version: 1.0

When `security.fips.enforced: "true"` is set in the [ClickHouseOperatorConfiguration], the operator SHALL coerce unset or relaxed security settings as follows:

* Unset TLS verify SHALL be coerced to Strict for ClickHouse, ZooKeeper/Keeper, and Kubernetes clients.
* Unset TLS `minVersion` SHALL be coerced to `"1.3"` for the operator's outbound TLS clients (`security.clickhouse.tls`, `security.zookeeper.tls`, and `security.kubernetes.tls`).
* Explicit `minVersion: "1.2"` for those TLS clients SHALL be coerced to `"1.3"`.
* Unset IPC mode SHALL be coerced to Secure.

Example configuration with explicit `minVersion: "1.2"`:

```yaml
spec:
  security:
    fips:
      enforced: "true"
    clickhouse:
      tls:
        minVersion: "1.2"
    zookeeper:
      tls:
        minVersion: "1.2"
    kubernetes:
      tls:
        minVersion: "1.2"
```

After operator configuration normalization, the effective `minVersion` for each TLS client listed above SHALL be `"1.3"`.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectInsecureKubeconfig
version: 1.0

The operator SHALL refuse to start when kubeconfig uses `TLSClientConfig.Insecure=true` under strict/FIPS mode.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.RejectNonCompliantSpecs
version: 1.0

When `security.fips.enforced: "true"` is set in the [ClickHouseOperatorConfiguration], the operator SHALL reject 
non-compliant CHI and CHK specifications with `FIPSValidationFailed` and SHALL NOT create workload StatefulSets for:

* CHI referencing plain external ZooKeeper nodes, including when `secure` is explicitly set to `"false"`.
* CHI with `clickhouse.tls.verify=None` at spec or cluster level.
* CHI with `zookeeper.tls.verify=None`.
* CHI with invalid `clickhouse.tls.minVersion`.
* CHK with TLS verify bypass at spec level.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Enforced.MinVersionScope
version: 1.0

The `minVersion` coercion SHALL apply only to TLS clients created and managed by the operator.
They SHALL NOT require ClickHouse Server or ClickHouse Keeper listener endpoints to reject TLS 1.2.

### Image Policy

#### RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.RejectNonFIPS
version: 1.0

With `security.fips.images.policy=Required`, non-FIPS images SHALL be rejected with `FIPSImagePolicyViolation` as follows:

* CHI with non-FIPS image tag SHALL be rejected at admission.
* CHK with non-FIPS Keeper image SHALL be rejected at admission.
* CHI with multiple non-FIPS hosts SHALL produce a single policy violation error.
* Digest-only image references SHALL NOT be detected as FIPS at admission.
* Registry hostname containing `fips` SHALL NOT satisfy FIPS tag detection.
* CHI admitted with a FIPS-tagged image whose running binary lacks `fips` in `SELECT version()` SHALL fail at runtime.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Images.Required.AcceptCHI
version: 1.0

With image policy Required, CHI with FIPS-tagged image SHALL reconcile normally.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Images.Permissive
version: 1.0

With permissive image policy, non-FIPS CHI images SHALL reconcile (default).


### Image Tag Detection

#### RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.FIPSSuffix
version: 1.0

Image tags containing `fips` (case-insensitive) SHALL be detected as FIPS.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.AltinityFIPS
version: 1.0

Image tags containing `altinityfips` SHALL be detected as FIPS.

#### RQ.SRS-026.ClickHouseOperator.FIPS.Images.TagDetection.CaseInsensitive
version: 1.0

Image tags such as `25.3.FIPS` or `25.3.Fips` SHALL be detected as FIPS (case-insensitive match on the tag).


### Operator Runtime Listener Verification

In a FIPS deployment, workload containers deployed by the operator (ClickHouse, Keeper, and sidecar containers) SHALL expose only expected TLS listener ports. Verification reads `/proc/net/tcp` and `/proc/net/tcp6` inside each container and parses ports in LISTEN state (`0A`):

```bash
kubectl exec <pod> -c clickhouse -- sh -c 'cat /proc/net/tcp /proc/net/tcp6'
```

E2e coverage: [`test_020011`](../e2e/test_operator_fips.py#L200).

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.Listeners
version: 1.0

FIPS workload pods (ClickHouse, Keeper, and sidecar containers) SHALL listen only on expected TLS ports. 
Plaintext service ports (8123, 9000, 2181) SHALL NOT be open when FIPS transport hardening applies. 
The clickhouse-operator pod network namespace SHALL expose only the expected Prometheus listener ports: `:8888` for 
metrics-exporter and `:9999` for clickhouse-operator, because both containers share the same pod network namespace.


### Operator to Kubernetes API

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.Kubernetes
version: 1.0

The operator SHALL connect to the Kubernetes API using FIPS-approved TLS ciphers.


### Operator to ClickHouse Server

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.ClickHouse
version: 1.0

The operator SHALL connect to ClickHouse using FIPS-approved TLS ciphers.


### Operator to ZooKeeper/Keeper

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.Zookeeper
version: 1.0

The operator SHALL connect to ZooKeeper/Keeper using FIPS-approved TLS ciphers.


### Operator to metrics-exporter IPC

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Operator.IPCSecure
version: 1.0

Operator IPC with `security.ipc.mode=Secure` SHALL work over localhost HTTP with token auth.


### Operator Prometheus Metrics

#### RQ.SRS-026.ClickHouseOperator.FIPS.Gap.OperatorMetricsTLS
version: 1.0

Operator Prometheus metrics on :9999 currently expose a known FIPS gap (HTTP-only).


## Exporter External Connections

**Objective:** Verify all **metrics-exporter** inbound and outbound connections use FIPS-compliant TLS.

### Exporter to Kubernetes API

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Exporter.Kubernetes
version: 1.0

The exporter SHALL connect to the Kubernetes API using FIPS-approved TLS ciphers.


### Exporter to ClickHouse Server

#### RQ.SRS-026.ClickHouseOperator.FIPS.Connect.Exporter.ClickHouse
version: 1.0

The exporter SHALL query ClickHouse using FIPS-approved TLS when configured for HTTPS.


### Exporter Prometheus Metrics

#### RQ.SRS-026.ClickHouseOperator.FIPS.Gap.ExporterMetricsTLS
version: 1.0

Exporter Prometheus metrics on :8888 currently expose a known FIPS gap (HTTP-only).


## Integrity Check Failure

**Objective:** Verify FIPS integrity self-test detects binary tampering for each shipped binary independently.


### Operator Integrity Tampering

#### RQ.SRS-026.ClickHouseOperator.FIPS.Integrity.OperatorMismatch
version: 1.0

Tampering with `clickhouse-operator` `.go.fipsinfo` SHALL panic with `fips140: verification mismatch`.


### Exporter Integrity Tampering

#### RQ.SRS-026.ClickHouseOperator.FIPS.Integrity.ExporterMismatch
version: 1.0

Tampering with `metrics-exporter` `.go.fipsinfo` SHALL panic with `fips140: verification mismatch`.


## CAST Failure

**Objective:** Verify FIPS Cryptographic Algorithm Self-Test (CAST) detects failures in each binary independently.


### Operator CAST Failure

#### RQ.SRS-026.ClickHouseOperator.FIPS.CAST.OperatorFail
version: 1.0

Running `clickhouse-operator` with `GODEBUG=failfipscast=<name>` SHALL terminate with a CAST error.


### Exporter CAST Failure

#### RQ.SRS-026.ClickHouseOperator.FIPS.CAST.ExporterFail
version: 1.0

Running `metrics-exporter` with `GODEBUG=failfipscast=<name>` SHALL terminate with a CAST error.


## Synthetic TLS Cipher Validation

**Objective:** Validate FIPS cipher enforcement on all external (to the pod) connections using `openssl s_client` and `openssl s_server`.

Use `openssl` to simulate connections with specific ciphers and verify the operator/exporter accepts FIPS-approved ciphers and rejects non-approved ones.

```bash
# Operator as TLS client against server offering only approved cipher
openssl s_server -accept 8443 -cert server.crt -key server.key \
  -ciphersuites TLS_AES_256_GCM_SHA384

# Operator as TLS client against server offering non-approved cipher
openssl s_server -accept 8443 -cert server.crt -key server.key \
  -cipher ECDHE-RSA-CHACHA20-POLY1305

# Inbound connection to operator/exporter metrics endpoint
openssl s_client -connect localhost:9999 -cipher ECDHE-RSA-AES256-GCM-SHA384
```

### Approved cipher matrix

#### RQ.SRS-026.ClickHouseOperator.FIPS.Synthetic.ApprovedCiphers
version: 1.0

For each external connection listed below, when exercised as a TLS **client** with `openssl s_server` offering only [approved ciphers](#rqsrs-026clickhouseoperatorfipstlsapprovedciphers), or as a TLS **server** with `openssl s_client` using only approved ciphers, the connection SHALL succeed:

| Connection | Role | Tool |
|------------|------|------|
| Operator to Kubernetes API | Client | `openssl s_server` |
| Operator to ClickHouse Server | Client | `openssl s_server` |
| Operator to ZooKeeper/Keeper | Client | `openssl s_server` |
| Operator metrics :9999 | Server | `openssl s_client` |
| Exporter to Kubernetes API | Client | `openssl s_server` |
| Exporter to ClickHouse Server | Client | `openssl s_server` |
| Exporter metrics :8888 | Server | `openssl s_client` |


### Rejected cipher matrix

#### RQ.SRS-026.ClickHouseOperator.FIPS.Synthetic.RejectedCiphers
version: 1.0

For each external connection listed below, when the peer offers only [rejected ciphers or protocols](#rqsrs-026clickhouseoperatorfipstlsrejectedciphers), the connection SHALL be rejected:

| Connection | Role | Tool |
|------------|------|------|
| Operator to Kubernetes API | Client | `openssl s_server` |
| Operator to ClickHouse Server | Client | `openssl s_server` |
| Operator to ZooKeeper/Keeper | Client | `openssl s_server` |
| Operator metrics :9999 | Server | `openssl s_client` |
| Exporter to Kubernetes API | Client | `openssl s_server` |
| Exporter to ClickHouse Server | Client | `openssl s_server` |
| Exporter metrics :8888 | Server | `openssl s_client` |


## CI/CD Image and Policy Verification

**Objective:** Add CI/CD jobs to validate FIPS image build and supply-chain checks.

### RQ.SRS-026.ClickHouseOperator.FIPS.CICD.OperatorImageBuild
version: 1.0

CI SHALL build the [clickhouse-operator] FIPS image successfully.

### RQ.SRS-026.ClickHouseOperator.FIPS.CICD.ExporterImageBuild
version: 1.0

CI SHALL build the [metrics-exporter] FIPS image successfully.

### RQ.SRS-026.ClickHouseOperator.FIPS.CICD.VulnerabilityScan
version: 1.0

FIPS images SHALL pass vulnerability scanning with no Critical, High, or Medium findings.


### Operator Source Review

#### RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Operator.Tree
version: 1.0

Static review of operator-scoped paths SHALL produce no Critical findings; Warning-level findings SHALL be documented.

#### RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Operator.SharedPkg
version: 1.0

Review of shared packages reachable from `cmd/operator` SHALL produce no Critical findings.

#### RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Operator.RegressionGate
version: 1.0

A signed-off review artifact SHALL be stored with the build record before release.

### Exporter Source Review

#### RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Exporter.Tree
version: 1.0

Static review of exporter-scoped paths SHALL produce no Critical findings; Warning-level findings SHALL be documented.

#### RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Exporter.SharedPkg
version: 1.0

Review of shared packages reachable from `cmd/metrics_exporter` SHALL produce no Critical findings.

#### RQ.SRS-026.ClickHouseOperator.FIPS.AIReview.Exporter.RegressionGate
version: 1.0

A signed-off review artifact SHALL be stored with the build record before release.

## ACVP Algorithm Validation

**Objective:** Reproduce ACVP expected-output checks for each FIPS binary using the tracked public-scope config in [`pkg/util/fips/acvp/`](../../../pkg/util/fips/acvp/).


### Operator ACVP Validation

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.WrapperIntegration
version: 1.0

Building clickhouse-operator with `-tags acvp_wrapper` SHALL expose a working ACVP responder via argv0 dispatch.

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.ConfigGeneration
version: 1.0

The clickhouse-operator ACVP responder SHALL answer `getConfig` with supported capabilities.

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.ExpectedOutputReplay
version: 1.0

`bash pkg/util/fips/acvp/run.sh` SHALL match all configured expected outputs for the operator.

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Operator.SuiteCount
version: 1.0

The tracked ACVP config SHALL report 38 matched expectations for clickhouse-operator.


### Exporter ACVP Validation

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.WrapperIntegration
version: 1.0

Building metrics-exporter with `-tags acvp_wrapper` SHALL expose a working ACVP responder.

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.ConfigGeneration
version: 1.0

The metrics-exporter ACVP responder SHALL answer `getConfig` with supported capabilities.

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.ExpectedOutputReplay
version: 1.0

`BINARY=metrics-exporter bash pkg/util/fips/acvp/run.sh` SHALL match all expected outputs.

#### RQ.SRS-026.ClickHouseOperator.FIPS.ACVP.Exporter.SuiteCount
version: 1.0

The tracked ACVP config SHALL report 38 matched expectations for metrics-exporter.

## Terminology

### SRS

Software Requirements Specification.

### FIPS 140-3

Federal Information Processing Standard for cryptographic module validation.

### clickhouse-operator

The Altinity ClickHouse Operator Kubernetes controller binary.

### metrics-exporter

The Prometheus metrics exporter sidecar binary shipped in the operator pod.

### CHI

ClickHouseInstallation custom resource.

### CHK

ClickHouseKeeperInstallation custom resource.

### ACVP

Automated Cryptographic Validation Protocol.

### CMVP

Cryptographic Module Validation Program.

### CAVP

Cryptographic Algorithm Validation Program.

[clickhouse-operator]: #clickhouse-operator
[metrics-exporter]: #metrics-exporter
[Kubernetes API]: https://kubernetes.io/docs/reference/kubernetes-api/
[ClickHouse Server]: #clickhouse-server
[ZooKeeper/Keeper]: #clickhouse-keeper
