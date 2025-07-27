#!/usr/bin/env -S uv run --script
# /// script
# dependencies = [
#   "rich",
#   "tabulate",
# ]
# ///

"""
Pretty print the entire proompt database contents.
This script displays all data in a nicely formatted, readable way.

Customize the output by editing:
- DISPLAY_OPTIONS: Control what gets shown and how
- format_*() functions: Modify how data is displayed
- Table styling and truncation settings
"""

import sqlite3
import os
import sys
import json
from datetime import datetime
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.text import Text
from rich.columns import Columns
from rich.tree import Tree
from rich import box
from tabulate import tabulate

# Configuration - Edit these to customize the output
DISPLAY_OPTIONS = {
    'max_content_length': 100,  # Truncate long content
    'max_rows_per_table': 50,   # Limit rows shown per table
    'show_empty_tables': True,  # Show tables even if empty
    'show_relationships': True, # Show foreign key relationships
    'show_metadata': True,      # Show table counts and timestamps
    'color_scheme': 'auto',     # 'auto', 'light', 'dark', 'none'
}

console = Console()

def get_db_path():
    """Get the database path from the config or use default."""
    return "./data/proompt.db"

def truncate_text(text, max_length=None):
    """Truncate text to specified length with ellipsis."""
    if max_length is None:
        max_length = DISPLAY_OPTIONS['max_content_length']
    
    if text is None:
        return ""
    
    text = str(text)
    if len(text) <= max_length:
        return text
    
    return text[:max_length-3] + "..."

def format_json_field(json_str):
    """Format JSON fields for display."""
    if not json_str:
        return ""
    
    try:
        data = json.loads(json_str)
        if isinstance(data, list):
            return ", ".join(str(item) for item in data)
        elif isinstance(data, dict):
            return ", ".join(f"{k}: {v}" for k, v in data.items())
        else:
            return str(data)
    except (json.JSONDecodeError, TypeError):
        return str(json_str)

def format_datetime(dt_str):
    """Format datetime strings for display."""
    if not dt_str:
        return ""
    
    try:
        # Try to parse and reformat
        dt = datetime.fromisoformat(dt_str.replace('Z', '+00:00'))
        return dt.strftime('%Y-%m-%d %H:%M')
    except (ValueError, AttributeError):
        return str(dt_str)

def get_table_info(cursor):
    """Get information about all tables in the database."""
    cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
    tables = cursor.fetchall()
    
    table_info = {}
    for table in tables:
        table_name = table[0]
        
        # Get row count
        cursor.execute(f"SELECT COUNT(*) FROM {table_name}")
        count = cursor.fetchone()[0]
        
        # Get column info
        cursor.execute(f"PRAGMA table_info({table_name})")
        columns = cursor.fetchall()
        
        table_info[table_name] = {
            'count': count,
            'columns': columns
        }
    
    return table_info

def display_table_summary(table_info):
    """Display a summary of all tables."""
    console.print("\n📊 [bold blue]Database Overview[/bold blue]")
    
    summary_table = Table(box=box.ROUNDED)
    summary_table.add_column("Table", style="cyan", no_wrap=True)
    summary_table.add_column("Records", justify="right", style="magenta")
    summary_table.add_column("Columns", justify="right", style="green")
    
    total_records = 0
    for table_name, info in table_info.items():
        summary_table.add_row(
            table_name,
            str(info['count']),
            str(len(info['columns']))
        )
        total_records += info['count']
    
    summary_table.add_row("", "", "", style="dim")
    summary_table.add_row("[bold]Total[/bold]", f"[bold]{total_records}[/bold]", "", style="bold")
    
    console.print(summary_table)

def display_prompts(cursor):
    """Display prompts table with rich formatting."""
    cursor.execute("""
        SELECT id, title, content, type, use_case, model_compatibility_tags, 
               temperature_suggestion, created_at, updated_at
        FROM prompts 
        ORDER BY updated_at DESC
        LIMIT ?
    """, (DISPLAY_OPTIONS['max_rows_per_table'],))
    
    rows = cursor.fetchall()
    if not rows and not DISPLAY_OPTIONS['show_empty_tables']:
        return
    
    console.print(f"\n🎯 [bold green]Prompts[/bold green] ({len(rows)} shown)")
    
    if not rows:
        console.print("  [dim]No prompts found[/dim]")
        return
    
    for row in rows:
        id_short = row[0][:8] + "..." if len(row[0]) > 8 else row[0]
        
        # Create a panel for each prompt
        content_preview = truncate_text(row[2], 150)
        models = format_json_field(row[5])
        
        prompt_info = f"""[bold]{row[1]}[/bold]
[dim]ID:[/dim] {id_short}
[dim]Type:[/dim] {row[3]} | [dim]Use Case:[/dim] {row[4] or 'N/A'}
[dim]Models:[/dim] {models}
[dim]Temperature:[/dim] {row[6] or 'N/A'} | [dim]Updated:[/dim] {format_datetime(row[8])}

{content_preview}"""
        
        console.print(Panel(prompt_info, border_style="green", padding=(0, 1)))

def display_snippets(cursor):
    """Display snippets table with rich formatting."""
    cursor.execute("""
        SELECT id, title, content, description, created_at, updated_at
        FROM snippets 
        ORDER BY updated_at DESC
        LIMIT ?
    """, (DISPLAY_OPTIONS['max_rows_per_table'],))
    
    rows = cursor.fetchall()
    if not rows and not DISPLAY_OPTIONS['show_empty_tables']:
        return
    
    console.print(f"\n📝 [bold yellow]Snippets[/bold yellow] ({len(rows)} shown)")
    
    if not rows:
        console.print("  [dim]No snippets found[/dim]")
        return
    
    for row in rows:
        id_short = row[0][:8] + "..." if len(row[0]) > 8 else row[0]
        content_preview = truncate_text(row[2], 100)
        
        snippet_info = f"""[bold]{row[1]}[/bold]
[dim]ID:[/dim] {id_short} | [dim]Updated:[/dim] {format_datetime(row[5])}
[dim]Description:[/dim] {row[3] or 'N/A'}

[cyan]{content_preview}[/cyan]"""
        
        console.print(Panel(snippet_info, border_style="yellow", padding=(0, 1)))

def display_notes(cursor):
    """Display notes with their associated prompts."""
    cursor.execute("""
        SELECT n.id, n.prompt_id, n.title, n.body, n.created_at, p.title as prompt_title
        FROM notes n
        LEFT JOIN prompts p ON n.prompt_id = p.id
        ORDER BY n.created_at DESC
        LIMIT ?
    """, (DISPLAY_OPTIONS['max_rows_per_table'],))
    
    rows = cursor.fetchall()
    if not rows and not DISPLAY_OPTIONS['show_empty_tables']:
        return
    
    console.print(f"\n📋 [bold cyan]Notes[/bold cyan] ({len(rows)} shown)")
    
    if not rows:
        console.print("  [dim]No notes found[/dim]")
        return
    
    table = Table(box=box.SIMPLE)
    table.add_column("Note", style="cyan", width=30)
    table.add_column("Prompt", style="green", width=25)
    table.add_column("Content", style="white", width=50)
    table.add_column("Created", style="dim", width=12)
    
    for row in rows:
        table.add_row(
            truncate_text(row[2], 28),
            truncate_text(row[5], 23),
            truncate_text(row[3], 48),
            format_datetime(row[4])
        )
    
    console.print(table)

def display_tags(cursor):
    """Display tags and their usage."""
    console.print(f"\n🏷️  [bold magenta]Tags[/bold magenta]")
    
    # Prompt tags
    cursor.execute("""
        SELECT tag_name, COUNT(*) as count
        FROM prompt_tags
        GROUP BY tag_name
        ORDER BY count DESC, tag_name
    """)
    prompt_tags = cursor.fetchall()
    
    # Snippet tags
    cursor.execute("""
        SELECT tag_name, COUNT(*) as count
        FROM snippet_tags
        GROUP BY tag_name
        ORDER BY count DESC, tag_name
    """)
    snippet_tags = cursor.fetchall()
    
    if prompt_tags or snippet_tags:
        columns = []
        
        if prompt_tags:
            prompt_table = Table(title="Prompt Tags", box=box.SIMPLE)
            prompt_table.add_column("Tag", style="cyan")
            prompt_table.add_column("Count", justify="right", style="magenta")
            
            for tag, count in prompt_tags:
                prompt_table.add_row(tag, str(count))
            columns.append(prompt_table)
        
        if snippet_tags:
            snippet_table = Table(title="Snippet Tags", box=box.SIMPLE)
            snippet_table.add_column("Tag", style="yellow")
            snippet_table.add_column("Count", justify="right", style="magenta")
            
            for tag, count in snippet_tags:
                snippet_table.add_row(tag, str(count))
            columns.append(snippet_table)
        
        console.print(Columns(columns, equal=True, expand=True))
    else:
        console.print("  [dim]No tags found[/dim]")

def display_relationships(cursor):
    """Display prompt links and relationships."""
    if not DISPLAY_OPTIONS['show_relationships']:
        return
    
    cursor.execute("""
        SELECT pl.from_prompt_id, pl.to_prompt_id, pl.link_type, 
               p1.title as from_title, p2.title as to_title
        FROM prompt_links pl
        LEFT JOIN prompts p1 ON pl.from_prompt_id = p1.id
        LEFT JOIN prompts p2 ON pl.to_prompt_id = p2.id
        ORDER BY pl.created_at DESC
        LIMIT ?
    """, (DISPLAY_OPTIONS['max_rows_per_table'],))
    
    rows = cursor.fetchall()
    if not rows and not DISPLAY_OPTIONS['show_empty_tables']:
        return
    
    console.print(f"\n🔗 [bold blue]Prompt Relationships[/bold blue] ({len(rows)} shown)")
    
    if not rows:
        console.print("  [dim]No relationships found[/dim]")
        return
    
    table = Table(box=box.SIMPLE)
    table.add_column("From", style="green", width=25)
    table.add_column("Type", style="yellow", width=12)
    table.add_column("To", style="cyan", width=25)
    
    for row in rows:
        table.add_row(
            truncate_text(row[3], 23),
            row[2],
            truncate_text(row[4], 23)
        )
    
    console.print(table)

def display_raw_table(cursor, table_name, table_info):
    """Display a raw table in tabular format (fallback)."""
    cursor.execute(f"SELECT * FROM {table_name} LIMIT ?", (DISPLAY_OPTIONS['max_rows_per_table'],))
    rows = cursor.fetchall()
    
    if not rows and not DISPLAY_OPTIONS['show_empty_tables']:
        return
    
    console.print(f"\n📄 [bold white]{table_name.title()}[/bold white] ({len(rows)} shown)")
    
    if not rows:
        console.print("  [dim]No data found[/dim]")
        return
    
    # Get column names
    columns = [col[1] for col in table_info['columns']]
    
    # Format data for display
    formatted_rows = []
    for row in rows:
        formatted_row = []
        for i, value in enumerate(row):
            if columns[i] in ['content', 'body']:
                formatted_row.append(truncate_text(value, 50))
            elif 'json' in columns[i].lower() or columns[i] in ['model_compatibility_tags', 'other_parameters']:
                formatted_row.append(format_json_field(value))
            elif 'created_at' in columns[i] or 'updated_at' in columns[i]:
                formatted_row.append(format_datetime(value))
            else:
                formatted_row.append(truncate_text(str(value) if value is not None else "", 30))
        formatted_rows.append(formatted_row)
    
    # Truncate column names for display
    display_columns = [truncate_text(col, 15) for col in columns]
    
    table_str = tabulate(formatted_rows, headers=display_columns, tablefmt="grid", maxcolwidths=30)
    console.print(table_str)

def view_database():
    """View and pretty print the entire database."""
    db_path = get_db_path()
    
    if not os.path.exists(db_path):
        console.print(f"[red]Database file not found at: {db_path}[/red]")
        console.print("Make sure you're running this from the server directory.")
        sys.exit(1)
    
    console.print(f"[bold green]📊 Proompt Database Viewer[/bold green]")
    console.print(f"[dim]Database: {db_path}[/dim]")
    
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Get table information
        table_info = get_table_info(cursor)
        
        if not table_info:
            console.print("\n[yellow]No tables found in database.[/yellow]")
            return
        
        # Display overview
        if DISPLAY_OPTIONS['show_metadata']:
            display_table_summary(table_info)
        
        # Display each table with custom formatting
        display_prompts(cursor)
        display_snippets(cursor)
        display_notes(cursor)
        display_tags(cursor)
        display_relationships(cursor)
        
        # Display any other tables not handled above
        handled_tables = {'prompts', 'snippets', 'notes', 'prompt_tags', 'snippet_tags', 'prompt_links'}
        for table_name, info in table_info.items():
            if table_name not in handled_tables:
                display_raw_table(cursor, table_name, info)
        
        console.print(f"\n[dim]Total tables: {len(table_info)}[/dim]")
        
    except sqlite3.Error as e:
        console.print(f"[red]Database error: {e}[/red]")
        sys.exit(1)
    finally:
        if conn:
            conn.close()

def customize_display():
    """Show customization options."""
    console.print("\n💡 [bold]Customization Options:[/bold]")
    console.print("Edit DISPLAY_OPTIONS in this script to customize:")
    console.print("  • max_content_length: Truncate long text")
    console.print("  • max_rows_per_table: Limit rows shown")
    console.print("  • show_empty_tables: Show/hide empty tables")
    console.print("  • show_relationships: Show/hide prompt links")
    console.print("  • show_metadata: Show/hide table summaries")

if __name__ == "__main__":
    # Check for help flag
    if len(sys.argv) > 1 and sys.argv[1] in ['-h', '--help', 'help']:
        console.print("🔍 [bold]Proompt Database Viewer[/bold]")
        console.print("\nDisplays all database contents in a nicely formatted way.")
        console.print("\n[bold]Usage:[/bold]")
        console.print("  uv run scripts/view_db.py")
        console.print("  ./scripts/view_db.py")
        customize_display()
        sys.exit(0)
    
    view_database()