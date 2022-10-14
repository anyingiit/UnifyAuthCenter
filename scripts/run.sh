docker run -t --name unify_auth_center \
  -p 8066:8066 \
  --network http-proxy-nginx-bridge \
  -e TZ="Asia/Shanghai" \
  unify_auth_center:V0.1 \
\
& docker logs -f unify_auth_center