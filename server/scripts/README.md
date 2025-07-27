# Development Scripts

This directory contains UV scripts for development tasks.

## Scripts

### `clear_db.py`
Clears all data from the proompt database while preserving the schema.

```bash
uv run scripts/clear_db.py
# or if executable:
./scripts/clear_db.py
```

**Warning**: This will delete ALL data from the database. Use with caution.

### `seed_data.py`
Seeds the database with realistic sample data for development and testing.

```bash
uv run scripts/seed_data.py
# or if executable:
./scripts/seed_data.py
```

**Features**:
- Creates sample prompts with realistic content
- Generates code snippets and documentation examples
- Adds notes, tags, and relationships
- Easily customizable by editing the script

**Customization**:
To modify what data gets generated, edit the following in `seed_data.py`:
- `SAMPLE_COUNTS`: Adjust quantities
- `generate_prompt_data()`: Modify predefined prompts
- `generate_snippet_data()`: Modify predefined snippets
- Tag lists and generation logic

### `view_db.py`
Pretty prints the entire database contents in a nicely formatted, readable way.

```bash
uv run scripts/view_db.py
# or if executable:
./scripts/view_db.py
```

**Features**:
- Rich, colorized output with tables and panels
- Shows database overview with record counts
- Displays prompts, snippets, notes, tags, and relationships
- Truncates long content for readability
- Customizable display options

**Customization**:
Edit `DISPLAY_OPTIONS` in the script to customize:
- `max_content_length`: Truncate long text
- `max_rows_per_table`: Limit rows shown per table
- `show_empty_tables`: Show/hide empty tables
- `show_relationships`: Show/hide prompt links
- `show_metadata`: Show/hide table summaries

## Requirements

These scripts use UV (Python package manager) with inline script metadata. They will automatically install required dependencies when run.

Dependencies:
- `sqlite3` (built-in)
- `faker` (for generating realistic test data)
- `rich` (for colorized output and formatting)
- `tabulate` (for table formatting)

## Usage Notes

- Run these scripts from the server directory (`/server`)
- The database path is configured to use `./data/proompt.db`
- Make sure the database has been initialized before seeding
- Both scripts include confirmation prompts for safety