# cf-metrics config & setup
the only configuration requirement is that the cf-metrics binary is run from a shell that has recently logged into cloud foundry via `cf login` and that said logged-in user has proper permissions.

# output
the binary will output a csv file for each org and space in the foundry inside of a directory called "output"

The binary currently hits the events endpoint for app creates, starts, and updates on a per org basis, and on a per space basis
It will also hit the apps endpoint on a per org basis and a per space basis.

It will then output this info for each org and each space into separate csv files located in the output directory
