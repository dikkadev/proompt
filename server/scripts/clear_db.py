#!/usr/bin/env -S uv run --script
# /// script
# dependencies = []
# ///

"""
Clear all data from the proompt database.
This script will delete all records from all tables while preserving the schema.
"""

import sqlite3
import os
import sys
from pathlib import Path

def get_db_path():
    """Get the database path from the config or use default."""
    # Default path based on the config file
    return "./data/proompt.db"

def clear_database():
    """Clear all data from the database tables."""
    db_path = get_db_path()
    
    if not os.path.exists(db_path):
        print(f"Database file not found at: {db_path}")
        print("Make sure you're running this from the server directory.")
        sys.exit(1)
    
    print(f"Clearing database at: {db_path}")
    
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Disable foreign key constraints temporarily
        cursor.execute("PRAGMA foreign_keys = OFF")
        
        # Get all table names, excluding FTS internal tables
        cursor.execute("""
            SELECT name FROM sqlite_master 
            WHERE type='table' 
            AND name NOT LIKE 'sqlite_%'
            AND name NOT LIKE '%_fts_data'
            AND name NOT LIKE '%_fts_idx'
            AND name NOT LIKE '%_fts_docsize'
            AND name NOT LIKE '%_fts_config'
        """)
        tables = cursor.fetchall()
        
        if not tables:
            print("No tables found to clear.")
            return
        
        print(f"Found {len(tables)} tables to clear:")
        for table in tables:
            print(f"  - {table[0]}")
        
        # Clear each table
        for table in tables:
            table_name = table[0]
            cursor.execute(f"DELETE FROM {table_name}")
            print(f"Cleared table: {table_name}")
        
        # Rebuild FTS tables to maintain integrity
        fts_tables = ['prompts_fts', 'snippets_fts', 'notes_fts']
        for fts_table in fts_tables:
            try:
                cursor.execute(f"INSERT INTO {fts_table}({fts_table}) VALUES('rebuild')")
                print(f"Rebuilt FTS table: {fts_table}")
            except sqlite3.Error:
                # FTS table might not exist yet, that's okay
                pass
        
        # Re-enable foreign key constraints
        cursor.execute("PRAGMA foreign_keys = ON")
        
        # Commit changes
        conn.commit()
        print("\nDatabase cleared successfully!")
        
    except sqlite3.Error as e:
        print(f"Database error: {e}")
        sys.exit(1)
    finally:
        if conn:
            conn.close()

if __name__ == "__main__":
    print("🗑️  Proompt Database Cleaner")
    print("=" * 30)
    
    # Confirm action
    response = input("Are you sure you want to clear ALL data from the database? (yes/no): ")
    if response.lower() not in ['yes', 'y']:
        print("Operation cancelled.")
        sys.exit(0)
    
    clear_database()