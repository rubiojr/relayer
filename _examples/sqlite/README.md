relayer sqlite
=============

  - a basic relay implementation based on relayer.
  - uses SQLite
  - it has some antispam limits, tries to delete old stuff so things don't get out of control, and some other small optimizations.

running
-------

    SQLITE_DATABASE=relayer.db ./relayer-sqlite

compiling
---------

    make
