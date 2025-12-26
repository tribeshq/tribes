<br>
<p align="center">
    <img src="https://github.com/user-attachments/assets/446065f2-e029-4634-a3da-bb1a3e82fe67" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>A Linux-powered EVM rollup as a Debit Capital Market</i>
</div>
<div align="center">
<b>Tokenized debt issuance through reverse auction mechanism with collateralization</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/tribeshq/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=959CD0" alt="license">
	<img src="https://img.shields.io/github/last-commit/tribeshq/tribes?style=default&logo=git&logoColor=white&color=D1DCCB" alt="last-commit">
</p>


## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment](#environment)
  - [Running](#running)
- [Testing](#testing)
- [Development](#development)

## Getting Started

### Prerequisites

1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/);

   To install Docker RISC-V support without using Docker Desktop, run the following command:

   ```shell
   docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download);

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps on devnet. To install it, run:

   ```shell
   npm install -g @cartesi/cli
   ```

4. [Install Foundry](https://book.getfoundry.sh/getting-started/installation) for smart contract development and testing;

5. [Install Go](https://golang.org/doc/install) (version 1.24.4 or later) for backend development;

### Environment

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

#### Contracts

The contract suite features **emergency delegate call operations** for secure asset recovery in critical situations. The system also includes **asset contracts** comprising a Stablecoin and Collateral token for the debt capital market operations, and **ERC1155 Badge contracts** for managing badges.

Each deployment script saves its configuration to individual JSON files in the `./contracts/deployments/` directory for easy reference and integration.

1. Deploy all contracts:

   ```sh
   make deploy-contracts
   ```

2. Deploy individual contracts:

   ```sh
   # Deploy BadgeFactory contract (ERC1155 Badge factory)
   make deploy-badge-factory

   # Deploy Tokens (Collateral and Stablecoin)
   make deploy-tokens

   # Deploy Emergency contracts (EmergencyWithdraw)
   make deploy-emergency

   # Deploy SafeERC1155Mint contract (Safe ERC1155 minting)
   make deploy-safe-erc1155-mint
   ```

3. Simulate deployment (without broadcasting):

   ```sh
   make deploy-contracts-simulate
   ```

   This is useful for testing deployment scripts and verifying gas costs without actually deploying contracts.

#### Backend

The backend is built on [Cartesi Rollups](https://cartesi.io/), a Layer 2 scaling solution that combines the security of blockchain with the computational power of Linux. This architecture enables complex off-chain computations while maintaining cryptographic guarantees of correctness and data availability. The system runs a full Linux environment inside the blockchain, handling business logic off-chain for better performance while keeping all computations verifiable on-chain.

1. Generate bytecode and Go bindings

   ```sh
   make generate
   ```

2. Devnet

   2.1 Build application:
   ```sh
   make build
   ```

   2.2 Run application on devnet:
   ```sh
   cartesi run
   ```

## Testing

Run all tests (Contracts + Backend):

```sh
make test
```

This command will:

- Clean and test smart contracts with Foundry;
- Run mock tests with Go;
- Run integration tests with Vitest;

## Development

### Code Quality

1. Format all code (Contracts + Backend):

   ```sh
   make fmt
   ```

2. View test coverage report:

   ```sh
   make coverage
   ```

### Utility Commands

1. Check contract sizes:

   ```sh
   make size
   ```

   Shows the size of all compiled contracts to ensure they fit within deployment limits.

2. Run gas reports:

   ```sh
   make gas
   ```

   Generates detailed gas usage reports for all contract functions during testing.

### Available Make Commands

For a complete list of available commands:

```sh
make help
```

This will show all available make targets with their descriptions.
