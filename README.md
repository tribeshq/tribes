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
  - [Running](#running)
  - [Testing](#testing)

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
   npm i -g @cartesi/cli
   ```

### Running

#### Frontend

#### Contracts

The contract suite features **emergency delegate call operations** for secure asset recovery in critical situations, **cross-chain NFT functionality** using Chainlink CCIP for seamless asset transfers across blockchains, and **Vlayer integration** with WebProofs that enables verification of social media profile ownership (X, Instagram, etc.) through TLSNotary and zero-knowledge proofs. The system also includes **asset contracts** comprising a Stablecoin and Collateral token for the debt capital market operations.

> [!WARNING]
> Make sure that all variables are defined in the .env file, which can be created with `make env`, before running any of the contract-related commands below.

1. Deploy all contracts:

   ```sh
   make contracts
   ```

2. Deploy Assets (Collateral and Stablecoin):

   ```sh
   make deploy-assets
   ```

3. Deploy `EmergencyWithdraw.sol` contract:

   ```sh
   make delegatecall
   ```
   
4. Deploy Vlayer contracts:

   ```sh
   make deploy-vlayer
   ```

5. Deploy cross-chain NFT contracts:

   ```sh
   make deploy-nft
   ```

6. Setup application (transfer contracts ownership):

   ```sh
   make setup
   ```

#### Backend

The backend is built on **Cartesi Rollups**, a Layer 2 scaling solution that combines the security of blockchain with the computational power of Linux. This architecture enables complex off-chain computations while maintaining cryptographic guarantees of correctness and data availability. The system runs a full Linux environment inside the blockchain, handling business logic off-chain for better performance while keeping all computations verifiable on-chain.

> [!WARNING]
> After a new deployment of Vlayer-related contracts, ensure that the `rollup.toml` is correctly defined with the correct addresses. Then run the `make generate` command to generate the latest version of the auto-generated code and also define the new addresses as environment variables that will be used in the system.

1. Devnet

   1.1 Build application:
   ```sh
   cartesi build
   ```

   1.2 Run application on testnet:
   ```sh
   cartesi run
   ```

2. Testnet < locally >

   2.1 Build application:
   ```sh
   cartesi build
   ```

   2.2 Run the Cartesi Rollups Node with the application's initial snapshot attached:
   ```sh
   docker compose --env-file .env up -d
   ```

   2.3 Deploy and register the application:
   ```sh
   docker compose --project-name cartesi-rollups-node exec advancer \
		cartesi-rollups-cli deploy application shoal /var/lib/cartesi-rollups-node/snapshot \
		--epoch-length 720 \
		--self-hosted \
		--salt <salt> \
		--json
   ```

3. Testnet < cloud >

   3.1 Create application on Fly.io:
   ```sh
   fly app create cartesi-rollups-node
   ```

   3.2 Create postgres instance:
   ```sh
   fly postgres create \
	   	--initial-cluster-size 1 \
	   	--name cartesi-rollups-node-database \
	   	--vm-size shared-cpu-1x \
	   	--volume-size 1
   ```

   3.3 Attach database to the application:
   ```sh
   fly postgres attach cartesi-rollups-node-database \
           --app cartesi-rollups-node
   ```

   3.4 Setup environment variables:
   ```sh
   fly secrets set -a cartesi-rollups-node CARTESI_BLOCKCHAIN_HTTP_ENDPOINT=<web3-provider-http-endpoint>
   fly secrets set -a cartesi-rollups-node CARTESI_BLOCKCHAIN_WS_ENDPOINT=<web3-provider-ws-endpoint>
   fly secrets set -a cartesi-rollups-node CARTESI_AUTH_MNEMONIC=`<mnemonic>`
   fly secrets set -a cartesi-rollups-node CARTESI_DATABASE_CONNECTION=<connection_string>
   ```

   3.5 Deploy Cartesi Rollups Node:
   ```sh
   fly deploy -a cartesi-rollups-node
   ```

   3.6 Access machine via SSH:
   ```sh
   fly ssh console
   ```

### Testing
```sh
make test
```
