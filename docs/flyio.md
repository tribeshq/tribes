# Deploying Cartesi Rollups Node on Fly.io

This guide provides step-by-step instructions for deploying a Cartesi Rollups node on Fly.io cloud platform.

## Prerequisites

- [Fly.io CLI](https://fly.io/docs/hands-on/install-flyctl/) installed and authenticated
- A Fly.io account with billing information set up

## Step-by-Step Deployment

### 1. Create Application

Create a new Fly.io application for the Cartesi Rollups node:

```sh
fly launch --ha=false --no-deploy
```

### 2. Create PostgreSQL Database

Create a PostgreSQL instance for the application:

```sh
fly postgres create \
    --initial-cluster-size 1 \
    --name cartesi-rollups-node-database \
    --vm-size shared-cpu-1x \
    --volume-size 1
```

### 3. Attach Database to Application

Link the PostgreSQL database to your application:

```sh
fly postgres attach cartesi-rollups-node-database \
    --app cartesi-rollups-node
```

### 4. Configure Environment Variables

Set up the required environment variables for the Cartesi Rollups node:

```sh
# Blockchain endpoints
fly secrets set -a cartesi-rollups-node CARTESI_BLOCKCHAIN_HTTP_ENDPOINT=<web3-provider-http-endpoint>
fly secrets set -a cartesi-rollups-node CARTESI_BLOCKCHAIN_WS_ENDPOINT=<web3-provider-ws-endpoint>

# Authentication
fly secrets set -a cartesi-rollups-node CARTESI_AUTH_MNEMONIC=<mnemonic>

# Database connection
fly secrets set -a cartesi-rollups-node CARTESI_DATABASE_CONNECTION=<connection_string>
```

**Important Notes:**

- Replace `<web3-provider-http-endpoint>` with your blockchain RPC HTTP endpoint;
- Replace `<web3-provider-ws-endpoint>` with your blockchain RPC WebSocket endpoint;
- Replace `<mnemonic>` with your wallet mnemonic phrase (keep secure);
- Replace `<connection_string>` with the PostgreSQL connection string provided by Fly.io;

### 5. Deploy the Application

Deploy the Cartesi Rollups Node to Fly.io:

```sh
fly deploy -a cartesi-rollups-node --local-only
```

## Troubleshooting

### Check Application Status

Monitor the application status:

```sh
fly status -a cartesi-rollups-node
```

### View Logs

Check application logs for debugging:

```sh
fly logs -a cartesi-rollups-node
```

## Additional Resources

- [Fly.io Documentation](https://fly.io/docs/)
- [PostgreSQL on Fly.io](https://fly.io/docs/postgres/)