#!/usr/bin/env python3
"""
Script to convert plain text story to JSON format for the Boma short story.
Usage: python3 boma.py [story.txt]
If no file is provided, it will prompt for input.
"""

import json
import os
import sys
from datetime import datetime

def create_boma_json(text_file=None):
    print("=== Boma Story - Text to JSON Converter ===\n")
    
    # Get story details
    title = input("Enter story title (default: Boma): ").strip() or "Boma"
    date = input("Enter date (YYYY-MM-DD) or press Enter for today: ").strip()
    location = input("Enter location (default: Brooklyn, NY): ").strip() or "Brooklyn, NY"
    
    # Use today's date if none provided
    if not date:
        date = datetime.now().strftime("%Y-%m-%d")
    
    # Read content from file or terminal
    if text_file:
        # Try to find the file - check current directory and common locations
        possible_paths = [
            text_file,  # Original path
            os.path.join("api/public/poems", text_file),  # If just filename provided
            os.path.join("api/public", text_file),
        ]
        
        file_found = None
        for path in possible_paths:
            if os.path.exists(path):
                file_found = path
                break
        
        if not file_found:
            print(f"❌ Error: File '{text_file}' not found.")
            print(f"   Checked locations:")
            for path in possible_paths:
                print(f"   - {path}")
            print(f"\n   Current working directory: {os.getcwd()}")
            return
        
        print(f"\nReading story from: {file_found}")
        try:
            with open(file_found, 'r', encoding='utf-8') as f:
                content = f.read()
            # Remove trailing newlines but preserve internal formatting
            content = content.rstrip('\n')
        except Exception as e:
            print(f"❌ Error reading file: {e}")
            return
    else:
        print(f"\nNow paste your story content (press Enter twice when finished):")
        print("=" * 50)
        
        # Collect story lines
        story_lines = []
        empty_lines = 0
        
        while True:
            try:
                line = input()
                if line == "":
                    empty_lines += 1
                    if empty_lines >= 2:  # Two consecutive empty lines = end of story
                        break
                    story_lines.append("")  # Preserve empty lines in story
                else:
                    empty_lines = 0
                    story_lines.append(line)
            except EOFError:
                # Handle paste ending without double enter
                break
        
        # Join lines with \n
        content = "\n".join(story_lines)
    
    # Create JSON structure
    story_data = {
        "title": title,
        "date": date,
        "location": location,
        "content": content
    }
    
    # Create filename
    filename = "api/public/boma.json"
    
    # Ensure directory exists
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    
    # Write JSON file
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(story_data, f, indent=2, ensure_ascii=False)
    
    print(f"\n✅ Story saved as: {filename}")
    print(f"📝 Title: {title}")
    print(f"📅 Date: {date}")
    print(f"📍 Location: {location}")
    print(f"📄 Content length: {len(content)} characters")
    print(f"📄 Content preview:")
    print("-" * 30)
    print(content[:200] + "..." if len(content) > 200 else content)
    print("-" * 30)

if __name__ == "__main__":
    try:
        # Check if a text file was provided as argument
        text_file = sys.argv[1] if len(sys.argv) > 1 else None
        create_boma_json(text_file)
    except KeyboardInterrupt:
        print("\n\n👋 Goodbye!")
    except Exception as e:
        print(f"\n❌ Error: {e}")

