#!/bin/bash
docker run -d \
  --name unify_auth_center \
  -p 8066:8066 \
  --network http-proxy-nginx-bridge \
  -e TZ="Asia/Shanghai" \
  anyingiit/unify_auth_center:V0.6 \
&& docker logs -f unify_auth_center
