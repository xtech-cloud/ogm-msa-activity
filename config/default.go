package config

const defaultYAML string = `
service:
    name: omo.api.msa.activity
    address: :9602
    ttl: 15
    interval: 10
logger:
    level: trace
    dir: /var/log/msa/
database:
    lite: true
    mysql:
        address: 127.0.0.1:3306
        user: root
        password: mysql@OMO
        db: msa_activity
    sqlite:
        path: /tmp/msa-activity.db
`
