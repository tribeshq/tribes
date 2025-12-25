import { rollups, RollupsMachine } from "@tuler/node-cartesi-machine";
import {
  encodeAbiParameters,
  getAddress,
  getContractAddress,
  keccak256,
  concat,
  type Address,
  type Hex,
} from "viem";
import badgeArtifact from "../../assets/artifacts/Badge.json";

export const BADGE_BYTECODE = badgeArtifact.bytecode as Hex;

export const computeBadgeAddress = (
  factory: Address,
  salt: Hex,
  appContract: Address,
): Address => {
  const constructorArgs = encodeAbiParameters(
    [{ type: "address" }],
    [appContract],
  );
  const initCodeHash = keccak256(concat([BADGE_BYTECODE, constructorArgs]));
  return getContractAddress({
    bytecodeHash: initCodeHash,
    from: factory,
    opcode: "CREATE2",
    salt,
  });
};

export const MACHINE_STORED_DIR = ".cartesi/image";
export const MACHINE_RUNTIME_CONFIG = { skip_root_hash_check: true };

// These addresses are automatically created in the database by sqlite.go
export const ADMIN_ADDRESS = getAddress(
  "0xD554153658E8D466428Fa48487f5aba18dF5E628",
);
export const VERIFIER_ADDRESS = getAddress(
  "0xc2D8eb4a934AEc7268E414a3Fa3D20E0572d714b",
);

export const TOKEN_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000009",
);
export const CREATOR_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000007",
);
export const FACTORY_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000013",
);
export const COLLATERAL = getAddress(
  "0x0000000000000000000000000000000000000008",
);
export const APPLICATION_ADDRESS = getAddress(
  "0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e",
);
export const EMERGENCY_WITHDRAW_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000006",
);

export const INVESTOR_01_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000001",
);
export const INVESTOR_02_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000002",
);
export const INVESTOR_03_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000003",
);
export const INVESTOR_04_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000004",
);
export const INVESTOR_05_ADDRESS = getAddress(
  "0x0000000000000000000000000000000000000005",
);

// Cartesi Portal addresses
export const ERC20_PORTAL_ADDRESS = getAddress(
  "0xc700D6aDd016eECd59d989C028214Eaa0fCC0051",
);

export const createMachine = (): RollupsMachine => {
  return rollups(MACHINE_STORED_DIR, {
    runtimeConfig: MACHINE_RUNTIME_CONFIG,
  });
};

export const setupTimeValues = (): {
  baseTime: number;
  closesAt: number;
  maturityAt: number;
} => {
  const baseTime = Math.floor(Date.now() / 1000);
  const closesAt = baseTime + 5;
  const maturityAt = baseTime + 10;
  return { baseTime, closesAt, maturityAt };
};
