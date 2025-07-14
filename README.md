<br>
<p align="center">
    <img src="https://github.com/user-attachments/assets/465b7615-842a-4f92-9f68-f3ffb8670fda" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>A Linux-powered EVM rollup serving as a Debt Capital Market for the creator economy</i>
</div>
<div align="center">
<b>Tokenized debt issuance through reverse campaign mechanism with collateralization</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/tribeshq/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=959CD0" alt="license">
	<img src="https://img.shields.io/github/last-commit/tribeshq/tribes?style=default&logo=git&logoColor=white&color=D1DCCB" alt="last-commit">
</p>

## Table of Contents
- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Setup](#environment-setup)
  - [Running](#running)
  - [Testing](#testing)
  - [Development](#development)

## Overview
A debt capital market platform designed for the creator economy, enabling content creators to monetize their influence by issuing tokenized debt instruments collateralized. Through a reverse campaign mechanism, the platform connects creators with investors, offering a structured and transparent alternative to finance scalable businesses while leveraging the economic potential of their audiences, ensuring legal compliance and attractive returns for investors.
<br>

[![Docs]][Link-docs] [![Deck]][Link-deck]
	
[Docs]: https://img.shields.io/badge/Documentation-959CD0?style=for-the-badge
[Link-docs]: https://docs.google.com/document/d/1l5D6sn9DBbaJFtTCfIM1gxoH7-10fVi9t2tsNr942Rw/edit?tab=t.0#heading=h.dfmi5re7vy34

[Deck]: https://img.shields.io/badge/Pitch%20Deck-D1DCCB?style=for-the-badge
[Link-deck]: https://www.canva.com/design/DAGVvlTnNpM/GsV9c1XuhYRYCrPK5811GA/view?utm_content=DAGVvlTnNpM&utm_campaign=designshare&utm_medium=link&utm_source=editor

## Getting Started

### Prerequisites
1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).

   To install Docker RISC-V support without using Docker Desktop, run the following command:
    
   ```shell
   docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download).

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps on devnet. To install it, run:

   ```shell
   npm install -g @cartesi/cli@2.0.0-alpha.15
   ```

4. [Install Foundry](https://book.getfoundry.sh/getting-started/installation) for smart contract development and testing.

5. [Install Go](https://golang.org/doc/install) (version 1.21 or later) for backend development.

### Environment Setup

1. Create the environment variables file:

   ```sh
   make env
   ```

2. Edit the `.env` file with your configuration values.

3. Import your private key for contract deployment:

   ```sh
   cast wallet import defaultKey --interactive
   ```
   
This will prompt you to enter your private key securely for contract deployment operations.

### Running

#### Frontend

> [!NOTE]
> Frontend documentation will be added as the project evolves.

#### Contracts

The contract suite features **emergency delegate call operations** for secure asset recovery in critical situations, and [Vlayer WebProofs](https://book.vlayer.xyz/features/web.html) that enables **verification of social media profile ownership (X, Instagram, etc.)** through TLSNotary and zero-knowledge proofs. The system also includes **asset contracts** comprising a Stablecoin and Collateral token for the debt capital market operations.

> [!WARNING]
> Make sure that all variables are defined in the .env file, which can be created with `make env`, before running any of the contract-related commands below.

1. Deploy all contracts:

   ```sh
   make contracts
   ```

2. Deploy individual contracts:

   ```sh
   # Deploy Assets (Collateral and Stablecoin)
   make deploy-assets
   
   # Deploy Badge contract
   make deploy-badge
   
   # Deploy Vlayer contracts
   make deploy-vlayer
   
   # Deploy EmergencyWithdraw.sol contract
   make deploy-delegatecall
   
   # Deploy CREAT deployer proxy
   make deploy-creat-deployer-proxy
   
   # Deploy CREAT2 deployer proxy
   make deploy-creat2-deployer-proxy
   ```

#### Backend

The backend is built on [Cartesi Rollups](https://cartesi.io/), a Layer 2 scaling solution that combines the security of blockchain with the computational power of Linux. This architecture enables complex off-chain computations while maintaining cryptographic guarantees of correctness and data availability. The system runs a full Linux environment inside the blockchain, handling business logic off-chain for better performance while keeping all computations verifiable on-chain.

> [!WARNING]
> After a new deployment of Vlayer-related contracts, ensure that the `rollup.toml` is correctly defined with the correct addresses. Then run the `make generate` command to generate the latest version of the auto-generated code and also define the new addresses as environment variables that will be used in the system.

1. Generate bytecode and Go bindings:

   ```sh
   make generate
   ```

2. Devnet

   2.1 Build application:
   ```sh
   cartesi build
   ```

   2.2 Run application on testnet:
   ```sh
   cartesi run
   ```

3. Testnet < locally >

   3.1 Build application:
   ```sh
   cartesi build
   ```

   3.2 Run the Cartesi Rollups Node with the application's initial snapshot attached:
   ```sh
   docker compose --env-file .env up -d
   ```

   3.3 Deploy and register the application:
   ```sh
   docker compose --project-name cartesi-rollups-node exec advancer \
		cartesi-rollups-cli deploy application tribes /var/lib/cartesi-rollups-node/snapshot \
		--epoch-length 720 \
		--self-hosted \
		--salt <salt> \ # cast keccak256 "your-unique-string"
		--json
   ```

4. Testnet < cloud >

   For detailed Cartesi Rollups Node deployment instructions on Fly.io, see [docs/flyio.md](docs/flyio.md).

   4.1 Access machine via SSH:
   ```sh
   fly ssh console
   ```

   4.2 Create directory to store snapshot:
   ```sh
   mkdir -p /var/lib/cartesi-rollups-node/snapshots/tribes
   ```

   4.3 Download and extract initial snapshot:
   ```sh
   curl -L <initial-snapshot-arctifac> | tar -xz -C /var/lib/cartesi-rollups-node/snapshots/tribes
   ```

   4.4 Deploy and register the application:
   ```sh
   cartesi-rollups-cli deploy application tribes /var/lib/cartesi-rollups-node/snapshot \
		--epoch-length 720 \
		--self-hosted \
		--salt <salt> \ # cast keccak256 "your-unique-string"
		--json
   ```

### Testing

Run all tests (Contracts + Backend):

```sh
make test
```

This command will:
- Clean and test smart contracts with Foundry
- Generate Go bindings
- Run Go tests with coverage

### Development

#### Code Quality

1. Run linting and formatting checks:

   ```sh
   make lint
   ```

2. Format all code (Contracts + Backend):

   ```sh
   make fmt
   ```

3. View test coverage report:

   ```sh
   make coverage
   ```

#### Available Make Commands

For a complete list of available commands:

```sh
make help
```

This will show all available make targets with their descriptions.
