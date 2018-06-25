# cf-metrics config & setup
## place a config.yaml file in the same directory as the cf-metrics binary with the following fields
- apiAddress: sourced via running `cf api` from a shell that has targeted the cf instance via `cf target`
- username: admin creds sourced via `https://YOUR_PCF_ADDRESS/api/v0/deployed/products/cf-f7ff4e780edd70b48001/credentials/.uaa.admin_credentials
- password: admin creds sourced via the above link (`https://YOUR_PCF_ADDRESS/api/v0/deployed/products/cf-f7ff4e780edd70b48001/credentials/.uaa.admin_credentials)
