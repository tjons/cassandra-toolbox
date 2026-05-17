[![Go Reference](https://pkg.go.dev/badge/github.com/tjons/cassandra-toolbox.svg)](https://pkg.go.dev/github.com/tjons/cassandra-toolbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/tjons/cassandra-toolbox)](https://goreportcard.com/report/github.com/tjons/cassandra-toolbox)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/tjons/cassandra-toolbox/blob/master/LICENSE)
![Main pipeline](https://github.com/tjons/cassandra-toolbox/actions/workflows/main.yml/badge.svg)

# cassandra-toolbox
A collection of Go packages for making Apache Cassandra more developer-friendly in Go projects.

## qb
`qb` is a fluent-query builder for Cassandra Query Language (CQL), in Go. It makes building dynamic CQL queries easy, without requiring string interpolation, templating, or counting just how many bind `?` placeholders are in a query. To learn more about `qb`, [check out the package documentation](./qb/README.md).

## migrate
`migrate` is a migrations library for creating, managing and migrating schemas in Apache Cassandra. `migrate` was written to support the `gocql/v2` driver. [Learn more about `migrate`](./migrate/README.md).

## pages
`pages` is a library providing a simple solution to a surprisingly hard problem - paginating through a result set. `pages` makes this easy, with support for paging tokens, common validation helpers, and [easy interoperability with the official Go Cassandra driver project, `gocql`](https://github.com/apache/cassandra-gocql-driver/v2). [Learn more about `pages`](./pages/README.md)
