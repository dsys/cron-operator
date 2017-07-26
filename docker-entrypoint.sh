#! /bin/sh

set -euo pipefail

WEB=${WEB:-true}
SCHEDULER=${SCHEDULER:-true}

if [[ "$WEB" = true ]]; then
  kubectl proxy \
    --address=0.0.0.0 \
    --port=80 \
    --www=/build/www \
    --www-prefix=/ \
    --api-prefix=/k8s-api \
    --accept-hosts='^(.*)$' &
fi

if [[ "$SCHEDULER" = true ]]; then
  /build/operator &
fi

for job in `jobs -p`; do
  echo $job
  wait $job
done
