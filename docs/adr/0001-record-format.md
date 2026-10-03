# ADR-0001: Log Record Format

## Status
Accepted

## Context
This entry documents the decision for the format of the entry in the database.

## Decision
[4 bytes CRC32][4 bytes timestamp][4 bytes key size][4 bytes value size][key bytes][value bytes]

## Consequences
