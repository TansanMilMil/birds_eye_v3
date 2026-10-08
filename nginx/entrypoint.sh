#!/bin/sh
set -eu

if [ -z "${BIRDSEYE_CLOUDFRONT_SECRET:-}" ]; then
  echo "FATAL: BIRDSEYE_CLOUDFRONT_SECRET must be set and non-empty" >&2
  exit 1
fi

case "$BIRDSEYE_CLOUDFRONT_SECRET" in
  *[!A-Za-z0-9_-]*) echo "FATAL: BIRDSEYE_CLOUDFRONT_SECRET contains invalid characters (only A-Za-z0-9_- allowed)" >&2; exit 1 ;;
esac

envsubst '$BIRDSEYE_CLOUDFRONT_SECRET' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf
exec "$@"
