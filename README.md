# Blockchain System

A Go-based blockchain implementation with core blockchain functionality, networking, and virtual machine capabilities.

## Features

### Core Blockchain Components
- ✨ Block creation and validation
- 📝 Transaction handling with digital signatures
- 🔄 Custom virtual machine (VM) for smart contracts
- 📊 State management system

### Networking
- 🌐 Peer-to-peer communication
- 📡 Message broadcasting
- 🔧 Local transport implementation for testing

### Cryptography
- 🔑 Public/private key pairs using ECDSA
- ✍️ Transaction and block signing
- 🔗 Hash-based block linking

## Installation

```bash
git clone https://github.com/yourusername/blockchain-system.git
cd blockchain-system
go mod download
```

## Quick Start

```bash
# Build the project
make build

# Run the node
make run
## Testing
make test

## Architecture

```
.
├── core/           # Core blockchain implementation
├── crypto/         # Cryptographic utilities
├── network/        # P2P networking components
├── vm/            # Virtual machine implementation
└── types/         # Shared types and interfaces
