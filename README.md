<br>
<p align="center">
    <img src="https://github.com/user-attachments/assets/0fd1d05e-ce8f-4353-8a2b-7a0eafb847fd" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>An EVM Linux-powered rollup as a launchpad for crowdfundings</i>
</div>
<div align="center">
<b>Debt issuance through crowdfunding with collateralized tokenization of receivables</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/tribeshq/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=959CD0" alt="license">
	<img src="https://img.shields.io/github/last-commit/tribeshq/tribes?style=default&logo=git&logoColor=white&color=D1DCCB" alt="last-commit">
</p>

##  Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)
  - [Development](#running)
  - [Testing](#testing)

##  Overview

A crowdfunding platform designed for prominent content creators, enabling them to monetize their influence by issuing tokenized debt instruments collateralized exclusively by their tokenized future receivables. Based on Resolution No. 88 of the Brazilian Securities and Exchange Commission (CVM), the Brazilian SEC, the platform connects creators with a network of investors, offering a structured and transparent alternative to finance scalable businesses while leveraging the economic potential of their audiences, ensuring legal compliance and attractive returns for investors.
<br>

[![Docs]][Link-docs] [![Deck]][Link-deck]
	
[Docs]: https://img.shields.io/badge/Documentation-959CD0?style=for-the-badge
[Link-docs]: https://docs.google.com/document/d/1l5D6sn9DBbaJFtTCfIM1gxoH7-10fVi9t2tsNr942Rw/edit?tab=t.0#heading=h.dfmi5re7vy34

[Deck]: https://img.shields.io/badge/Pitch%20Deck-D1DCCB?style=for-the-badge
[Link-deck]: https://www.canva.com/design/DAGVvlTnNpM/GsV9c1XuhYRYCrPK5811GA/view?utm_content=DAGVvlTnNpM&utm_campaign=designshare&utm_medium=link&utm_source=editor

##  Getting Started

###  Prerequisites
1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).

    To install Docker RISC-V support without using Docker Desktop, run the following command:
    
   ```shell
   docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download).

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps. To install it, run:

   ```shell
   npm i -g @cartesi/cli
   ```

> [!IMPORTANT]
>  To run the system in development mode, it is required to install:
>
> 1. [Download and Install the latest version of Golang.](https://go.dev/doc/install)
> 2. Install development node:
>
>   ```shell
>   npm i -g nonodo
>   ```
> 3. Install air ( hot reload tool ):
>
>   ```shell
>   go install github.com/air-verse/air@latest
>   ```

###  Running

1. Production mode:

   1.1 Generate rollup filesystem:

   ```sh
   cartesi build
   ```

   1.2 Run validator node:

   ```sh
   cartesi rollups start
   ```

   1.3 Deploy application:

   ```sh
   cartesi rollups deploy
   ```

2. Unsandboxed mode:

   2.1 Generate rollup filesystem:

   ```sh
   cartesi build
   ```

   2.2 Run development node:

   ```sh
   nonodo
   ``` 

   2.3 Start the application inside a Cartesi Machine unsandboxed:

   ```sh
   cartesi-machine --network \
         --flash-drive=label:root,filename:.cartesi/image.ext2 \
         --env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004 -- /var/opt/cartesi-app/app
   ```

###  Development

1. Run development node and application w/ hot reload:

   ```sh
   nonodo -- air
   ```

> [!NOTE]
> To reach the final state of the system, run the command bellow:
>
>   ```shell
>   make state
>   ```

###  Testing

```sh
make test
```
