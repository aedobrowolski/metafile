# metafile

A file system abstraction that adds persisted user metadata to entries.

## Introduction

Metafile is a file system with metadata.  The abstraction that we wrap 
is billy.  Every directory or file in the file system implements the 
Metadata interface, allowing client code to treat it as a key-value 
map where both are strings.
