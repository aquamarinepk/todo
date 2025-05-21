# Seeding

Seeding using `.sql` files can be convenient for inserting large volumes of
raw data, even when it's associated with persistent entities in the system.
However, this approach can be limiting when certain values are the result of
business logic that is hard or inconvenient to reproduce, for example,
digests, encrypted fields, or other derived data that is not readily
available at seeding time.

We prefer JSON-based seeding (`.json`) when we need to populate concrete
system entities with business-related precomputed values and/or clearly defined relationships, expressed directly and explicitly in the file.

## Example SQL seed file

To add SQL-based seed data, create a file in:

    assets/seed/{engine}/yyyymmddhhmmss-name.sql

For example:

    assets/seed/sqlite/20250325183118-add-sample-data.sql

```
-- Superadmin user
INSERT INTO users (id, username, email_enc, name, password_enc, slug, created_by, updated_by, created_at, updated_at, last_login_at, last_login_ip, is_active)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'superadmin', X'73616d706c655f656d61696c', 'Super Admin', X'70617373776f7264313233', 'superadmin', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '127.0.0.1', 1);

-- Admin user
INSERT INTO users (id, username, email_enc, name, password_enc, slug, created_by, updated_by, created_at, updated_at, last_login_at, last_login_ip, is_active)
VALUES 
    ('00000000-0000-0000-0000-000000000002', 'admin', X'61646d696e5f656d61696c', 'Admin', X'70617373776f7264313233', 'admin', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '127.0.0.1', 1);

-- ...
```

## Example JSON seed file

To add JSON-based seed data, create a file in:

    assets/seed/{engine}/yyyymmddhhmmss-feat-name.json

For example:

    assets/seed/sqlite/2025052000321-auth-add-sample-data.json

- `{engine}` is your database engine (e.g., `sqlite`).
- The filename should start with a timestamp, followed by the feature (e.g., `auth`), and a descriptive name.
- The file should contain a JSON object matching the expected structure for the feature.

### JSON Seeding Structure

A JSON seed file is a single JSON object (`{ ... }`).

A more representative JSON seed file might look like this:

```json
{
  "users": [
    {
      "ref": "user-superadmin",
      "username": "superadmin",
      "email": "superadmin@example.com",
      "name": "Super Admin",
      "password": "password123",
      "is_active": true
    }
    // ... more users ...
  ],
  "orgs": [
    {
      "ref": "org-aquamarine",
      "name": "Aquamarine",
      "short_description": "Aquamarine Org",
      "description": "Main organization for Aquamarine"
    }
    // ... more orgs ...
  ],
  "org_owners": [
    { "org_ref": "org-aquamarine", "user_ref": "user-superadmin" }
    // ... more org owners ...
  ]
  // ... more entity arrays ...
}
```

- The property names (e.g., `username`, `name`) correspond to the JSON tags of the business objects being seeded.
- The `ref` property is a special, human-friendly identifier used only within the seed file to reference entities. It is not persisted in the database, but allows you to establish relationships between entities in a readable way.
- IDs (UUIDs) are generated during insertion. You cannot use them directly to connect entities in the seed file, so use `ref` for associations.
- For relationships, use properties like `xxx_ref` (e.g., `org_ref`, `user_ref`). These fields are only used for mapping relationships in the JSON; at runtime, the seeder will resolve each `xxx_ref` to the corresponding generated ID and set it in the appropriate `xxx_id` field. The database always stores traditional IDs (UUIDs), not the `ref` values.
- For example, in an `org_owners` array, `org_ref` and `user_ref` are used to map to the correct IDs at seeding time, and the resulting database records will use the resolved `org_id` and `user_id` UUIDs.

This approach makes it easy to define business-related data and relationships in a clear, maintainable, and human-friendly way.
