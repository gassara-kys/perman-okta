# perman-okta

## export ldap env

```bash
# LDAP
$ export LDAP_HOST="localhost"
$ export BASE_DN="dc=example,dc=com"
$ export FILTER_STRING="(uid=hogehoge)"

# Okta
$ export OKTA_FQDN="example.okta.com"
$ export OKTA_APIKEY="xxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

## run

```bash
$ cd path/to/app
$ ./run.sh
```

