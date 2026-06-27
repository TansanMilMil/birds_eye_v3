#!/bin/sh
set -e

# Taskfile
sudo curl -sL https://taskfile.dev/install.sh | sudo sh -s -- -b /usr/local/bin

# Go dependencies
cd /workspaces/birds_eye_v3/backend && go mod download

# Node dependencies
cd /workspaces/birds_eye_v3/frontend && npm ci

# Claude Code CLI
curl -fsSL https://claude.ai/install.sh | bash
