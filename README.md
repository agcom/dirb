# DirB

File-based approach; DirB is a very simple database written from scratch. Word "DirB" is a pun that mimic word "DB", mixed with an abbreviation of "Directory".

This project was made as a homework for a **database design principles** college course.

Currently, the only DirB's user interface is command line.

## CLI usage example

[![CLI usage example](https://asciinema.org/a/440994.svg)](https://asciinema.org/a/440994)

## Design

Design principles.

### Single directory

DirB works upon a single directory. Each file in the directory corresponds to an instance, and vice versa.

CLI: you can set the directory through `-d path` flag or it'll default to the working directory (e.g. the terminal's current directory).

### Schema-less

Not enforcing any schema for instances, other than being a **json object**; enforcing json makes querying possible.

### CRUD

Supports create, read, update, and delete.

CLI: `dirb create json [-d path]`, `dirb read name [-d path] [-p [bool]]`, `dirb update name json [-d path]`, and `dirb delete name [-d path]`.

### Query

Supports limited query operations.

CLI: `dirb ls [-d path]`, and `dirb find l op r [-l [bool]] [-r [bool]] [-d path]`.

### Dirty

TL;DR: the project probably contains bugs and unexpected behavior.

Currently, the main package (root directory of the project's src) contains a cluster of copy pastas and dirty codes (as to honor the deadline); "It works! But at what cost...", said the author; also, there are no tests at all.

### ACID

Supports ACID transactions on a single instance.

The strategy is to, per time, allow a single write and multiple reads. Note that if the underlying file-system doesn't support atomicity for common file operations (e.g. create, remove, and rename), then DirB can't guarantee what's discussed in this section.

### Daemon-less

"Fire... and... we're done", said DirB after each interaction.

