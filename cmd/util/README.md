# Encryption Utility

This is a temporary utility program used to generate encrypted values for database seeding.

## Purpose

The utility generates encrypted values for sensitive data (like emails and passwords) that need to be inserted into the database during seeding. This is a temporary solution until a more robust seeding mechanism is implemented.

## Current Implementation

The current implementation:
- Takes raw values (email, password)
- Encrypts them using the application's encryption key
- Outputs the values in SQL-compatible format (hex-encoded)

## Future Improvements

This utility is a temporary solution. Future improvements should include:
- Integration with the main seeding mechanism
- Support for pre-encryption of values during seed generation
- More robust handling of sensitive data
- Automated testing of encrypted values

## Usage

Run the utility to generate encrypted values for seeding:

```bash
go run main.go
```

The output will be in SQL-compatible format that can be directly used in seed files.

## Note

This is a temporary solution and should be replaced with a more robust seeding mechanism that handles encryption natively. 