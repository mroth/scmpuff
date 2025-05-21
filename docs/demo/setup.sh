#!/bin/bash
# 
# scmpuff demo recreation script
# This script recreates a git repository state and demonstrates scmpuff features
# as shown in the demo video at https://mroth.github.io/scmpuff/assets/scmpuffdemo-2x.mp4
#

set -e

echo "🎬 Setting up scmpuff demo repository..."

# Create a clean demo directory
DEMO_DIR="/tmp/scmpuff-demo"
mkdir -p "$DEMO_DIR"
cd "$DEMO_DIR"

echo "📁 Created demo directory: $DEMO_DIR"

# Initialize git repository
git init --quiet
git config user.name "Demo User"
git config user.email "demo@example.com"

echo "📦 Initialized git repository"

# Create some coffee-themed files (matching the playground theme)
touch README.md americano.md cappuccino.md chemex.go espresso.md macchiato.md presspot.go

echo "☕️ Created coffee-themed files"

# Make initial commit
git add chemex.go presspot.go README.md
git commit -m "Initial commit"

echo "✅ Made initial commit"
