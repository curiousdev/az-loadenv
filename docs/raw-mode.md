# Raw mode

## Overview

Some `.env` file parsers and SDK libraries do not handle escaped or quoted values well. The `--raw` flag writes values exactly as they come from Azure, with no quoting or escaping applied.

## Usage

```bash
az-loadenv --app my-api --rg my-resource-group --raw
```

## Comparison

Given an Azure app setting `DB_URL` with value `Server=db.example.com;Password=has spaces`:

**Default (escaped):**
```
DB_URL="Server=db.example.com;Password=has spaces"
```

**Raw mode:**
```
DB_URL=Server=db.example.com;Password=has spaces
```

## When to use raw mode

Use `--raw` when:
- Your `.env` parser doesn't support quoted values or interprets quotes as literal characters
- You're piping values into a tool that expects unquoted input
- Your values contain characters that your parser misinterprets when escaped

## Caveats

With `--raw`, values containing newlines, `#` characters, or leading/trailing whitespace may not parse correctly depending on the consumer. Test with your specific toolchain before committing to raw mode for production workflows.
