# Timestream Datasource

WIP: https://aws.amazon.com/timestream/

## Development

You need to install the following first:

* [Mage](https://magefile.org/)
* [Yarn](https://yarnpkg.com/)
* [Docker Compose](https://docs.docker.com/compose/)

```
mage watch
```

In another terminal
```
docker-compose up
```

To restart after backend changes:
`./scripts/restart-plugin.sh`
