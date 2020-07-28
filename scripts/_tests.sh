exit 0

./scripts/run.ps1 client --netgroup work --server ::1 --external-ip-nohttp --server-ca=/root/.wireguard-config-server/ca.cert.pem

./scripts/run.ps1  server --server-name=test --ip-nohttp
