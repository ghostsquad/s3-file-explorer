#!/bin/bash

success="false"

n=0
until [ "$n" -ge 5 ]
do
  if http --check-status --ignore-stdin --timeout=2.5 GET http://app:8080 &> /dev/null; then
    success="true"
    echo 'OK!' && break
  else
    case $? in
      2) echo 'Request timed out!' ;;
      3) echo 'Unexpected HTTP 3xx Redirection!' ;;
      4) echo 'HTTP 4xx Client Error!' ;;
      5) echo 'HTTP 5xx Server Error!' ;;
      6) echo 'Exceeded --max-redirects=<n> redirects!' ;;
      *) echo 'Other Error!' ;;
    esac
  fi
  n=$((n+1))
  sleep 0.1
done

if [ "${success}" == "true" ]; then
  exit 0
else
  exit 1
fi