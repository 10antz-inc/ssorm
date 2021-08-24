#!/bin/bash

gcloud config configurations create spanner-emulator
gcloud config set auth/disable_credentials true
gcloud config set project spanner-emulator
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/

gcloud beta emulators spanner start &>/dev/null
gcloud spanner instances create test --config spanner-emulator --description "test" --nodes=1
gcloud spanner databases create test --instance test

SPANNER_EMULATOR_HOST=localhost:9010 spanner-cli -p spanner-emulator -i test -d test -f ./tests/ddl/data.sql




