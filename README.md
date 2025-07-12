# Scrabbled - AI-Assisted Multiplayer Scrabble

## Project Overview & Features

Scrabbled is a multiplayer [Scrabble](docs/rules_of_scrabble.md) game implementation built entirely through AI-assisted development as an exercise in effective human-AI collaboration in software engineering. This project demonstrates how AI can be leveraged to design, architect, and implement a complete multiplayer game system.

**Key Features:**
- **Multiplayer gameplay** with real-time updates via WebSockets
- **Web-based client** - play directly in your browser, no downloads required
- **Persistent game state** - games survive server restarts and player disconnections
- **Player reconnection** - seamlessly rejoin games after network interruptions
- **Session management** - secure player authentication and game access
- **Automatic cleanup** - inactive games expire after 1 week
- **Complete Scrabble rules** - official board layout, tile distribution, and scoring
- **Dictionary validation** - comprehensive word checking for fair play

**AI Development Approach:**
This project showcases systematic AI-assisted development including requirements gathering, architecture design, test-driven development planning, and iterative implementation. The entire codebase, documentation, and project structure are generated through human-AI collaboration.

## Architecture Overview

- **Backend**: Go server with WebSocket communication and database persistence
- **Frontend**: HTML/CSS/JavaScript web client with drag-and-drop interface
- **Data Storage**: SQLite/PostgreSQL for game state, text files for dictionary (performance optimized)
- **Communication**: Real-time WebSocket protocol for game updates
- **Deployment**: Single binary server with embedded web assets

## Prerequisites

- Go 1.21 or later
- SQLite3 (for database storage)
- Modern web browser with WebSocket support

## Quick Start

*Note: This project is currently in development. The quick start steps will be updated as implementation progresses.*

```bash
# Clone the repository
git clone <repository-url>
cd scrabbled

# Build the server
go build -o scrabbled cmd/server/main.go

# Run the server
./scrabbled

# Open your browser to http://localhost:8080
```

## Project Status

**Current Phase**: Project Setup & Foundation
- ✅ Project structure and documentation
- ✅ Requirements and architecture design
- ✅ Implementation roadmap and testing strategy
- 🔄 **In Progress**: Core game logic implementation
- ⏳ **Upcoming**: WebSocket server, web client, database integration

See [`docs/implementation_checklist.md`](docs/implementation_checklist.md) for detailed progress tracking.

## Documentation

- **[Design Requirements](docs/design_requirements.md)** - System architecture and technical specifications
- **[Implementation Checklist](docs/implementation_checklist.md)** - Detailed development roadmap and progress tracking
- **[Rules of Scrabble](docs/rules_of_scrabble.md)** - Complete game rules and specifications

## License

This project is licensed under the MIT License - see the LICENSE file for details. 