#!/bin/bash

status_code=$(curl --write-out %{http_code} --silent --output /dev/null https://pipeline-checker.example/v1/pipeline/status/project/proj/environment/test)

if [[ "$status_code" -ne 423 ]] ; then
  exit 0
else
  echo "Pipeline is locked!"
  exit 1
fi