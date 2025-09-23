#!/usr/bin/env python3
"""
Simple script to convert plain text poems to JSON format for the Alaska Hoffman Poetry Archive.
"""

import json
import os
from datetime import datetime

def create_poem_json():
    print("=== Alaska Hoffman Poetry Archive - Poem to JSON Converter ===\n")
    
    # Get poem details
    poem_id = input("Enter poem ID number: ").strip()
    title = input("Enter poem title: ").strip()
    date = input("Enter date (YYYY-MM-DD) or press Enter for today: ").strip()
    category = input("Enter category (default: Poetry): ").strip() or "Poetry"
    location = input("Enter location (default: Brooklyn, NY): ").strip() or "Brooklyn, NY"
    
    # Use today's date if none provided
    if not date:
        date = datetime.now().strftime("%Y-%m-%d")
    
    print(f"\nNow paste your poem content (press Enter twice when finished):")
    print("=" * 50)
    
    # Collect poem lines
    poem_lines = []
    empty_lines = 0
    
    while True:
        line = input()
        if line == "":
            empty_lines += 1
            if empty_lines >= 2:  # Two consecutive empty lines = end of poem
                break
            poem_lines.append("")  # Preserve empty lines in poem
        else:
            empty_lines = 0
            poem_lines.append(line)
    
    # Join lines with \n
    content = "\n".join(poem_lines)
    
    # Create JSON structure
    poem_data = {
        "id": int(poem_id),
        "title": title,
        "date": date,
        "category": category,
        "location": location,
        "content": content
    }
    
    # Create filename
    filename = f"static/poems/poem-{poem_id}.json"
    
    # Ensure directory exists
    os.makedirs("static/poems", exist_ok=True)
    
    # Write JSON file
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(poem_data, f, indent=2, ensure_ascii=False)
    
    print(f"\n✅ Poem saved as: {filename}")
    print(f"📝 Title: {title}")
    print(f"📅 Date: {date}")
    print(f"📍 Location: {location}")
    print(f"📂 Category: {category}")
    print(f"📄 Content preview:")
    print("-" * 30)
    print(content[:200] + "..." if len(content) > 200 else content)
    print("-" * 30)
    
    # Ask if user wants to create another
    another = input("\nCreate another poem? (y/n): ").strip().lower()
    if another in ['y', 'yes']:
        print("\n" + "="*60 + "\n")
        create_poem_json()

if __name__ == "__main__":
    try:
        create_poem_json()
    except KeyboardInterrupt:
        print("\n\n👋 Goodbye!")
    except Exception as e:
        print(f"\n❌ Error: {e}")
