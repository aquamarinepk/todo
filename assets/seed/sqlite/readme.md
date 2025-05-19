# Seeding

Seeding using `.sql` files can be convenient for inserting large volumes of
raw data, even when it's associated with persistent entities in the system.
However, this approach can be limiting when certain values are the result of
business logic that is hard or inconvenient to reproduce, for example,
digests, encrypted fields, or other derived data that is not readily
available at seeding time.

We prefer JSON-based seeding (`.json`) when we need to populate concrete
system entities with precomputed values and clearly defined relationships,
expressed directly and explicitly in the file.

More details on how to structure this JSON document will be added soon.
