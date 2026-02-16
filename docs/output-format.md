# Output format

## Structure

Settings are written one per line in `KEY=VALUE` format, sorted alphabetically by key:

```
APP_NAME=my-api
DATABASE_URL="Server=db.example.com;Password=s3cret"
LOG_LEVEL=info
REDIS_HOST=redis.example.com
```

## Quoting and escaping

By default, values that contain special characters are automatically double-quoted and escaped. This ensures compatibility with most `.env` file parsers.

A value is quoted when it contains any of the following:
- Spaces or tabs
- Newlines or carriage returns
- Double quotes (`"`)
- Single quotes (`'`)
- Hash (`#`) — would otherwise start a comment
- Backticks (`` ` ``)
- Dollar signs (`$`) — would otherwise trigger variable expansion
- Backslashes (`\`)
- Leading or trailing whitespace

### Escape sequences

Inside quoted values, these characters are escaped:

| Character | Escaped as |
|---|---|
| `\` | `\\` |
| `"` | `\"` |
| newline | `\n` |
| carriage return | `\r` |

### Example

An Azure setting with value `Server=db.example.com;Password=has "quotes" & spaces` would be written as:

```
DB_CONNECTION="Server=db.example.com;Password=has \"quotes\" & spaces"
```

## File permissions

The output file is created with `0600` permissions (owner read/write only). This prevents other users on shared systems from reading secrets in the file.

## Atomic writes

The file is written atomically using a temp file + rename strategy:

1. A temporary file is created in the same directory as the output file
2. All content is written to the temp file
3. The temp file is renamed to the final path

This means the output file is never in a partially written state. If the process is interrupted, the previous `.env` file (if any) remains untouched.

## Failed secret resolution

If a Key Vault secret fails to resolve, the original `@Microsoft.KeyVault(SecretUri=...)` reference is written to the output file as-is. This allows the rest of the settings to still be usable. The error is logged to stderr.
