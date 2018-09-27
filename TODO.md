# TODO

```text
F    Feature
I    Issue
T    Test
```

Relative priority is indicated by position not issue number.

## High

* F01. Support persistence of sidecar files for metadata. Consider using BoltDB

## Medium

* F04. Panic recovery that will not loose cached values (via Restore).
* F05. Create a validation function for buckets that matches the file system when storeNew() called.

* T01. Test that correct error messages are generated with bad user input
* T02. Test moving directories

## Low

* F06. Better documentation and setup examples.

## Done

Move done issues to the top of this list.

* F02. File system factory methods: New and MemNew
* F03. Use of Gob binary format to store values and allow more flexible data fetch.

* T02. Test that encode/decode do round trip data.

## Obsolete

* F07. Consider mmap to keep in sync, see for example [adventures with mmap](https://medium.com/@arpith/adventures-with-mmap-463b33405223).
* F08. Implement values as []byte arrays using Gob encoding.