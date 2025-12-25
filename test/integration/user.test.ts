import { bytesToHex, stringToHex } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import { encodeAdvanceInput, encodeNoticeOutput } from "./encoder";
import {
  createMachine,
  ADMIN_ADDRESS,
  VERIFIER_ADDRESS,
  CREATOR_ADDRESS,
  INVESTOR_01_ADDRESS,
} from "./helpers";

describe("User Tests", () => {
  const machine = createMachine();

  it("should create creator user", () => {
    const baseTime = Math.floor(Date.now() / 1000);

    const createUserInput = JSON.stringify({
      path: "user/admin/create",
      data: {
        address: CREATOR_ADDRESS,
        role: "creator",
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createUserInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);

    const expectedNoticePayload = `user created - {"id":3,"role":"creator","address":"${CREATOR_ADDRESS}","social_accounts":[],"created_at":${baseTime}}`;
    const expectedOutput = encodeNoticeOutput({
      payload: stringToHex(expectedNoticePayload),
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  it("should create investor user", () => {
    const baseTime = Math.floor(Date.now() / 1000);

    const createUserInput = JSON.stringify({
      path: "user/admin/create",
      data: {
        address: INVESTOR_01_ADDRESS,
        role: "investor",
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createUserInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);

    const expectedNoticePayload = `user created - {"id":4,"role":"investor","address":"${INVESTOR_01_ADDRESS}","social_accounts":[],"created_at":${baseTime}}`;
    const expectedOutput = encodeNoticeOutput({
      payload: stringToHex(expectedNoticePayload),
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  it("should find all users", () => {
    const baseTime = Math.floor(Date.now() / 1000);

    const findAllUsersInput = JSON.stringify({
      path: "user",
    });

    const reports = machine.inspect(Buffer.from(findAllUsersInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));
    const expectedOutput = [
      {
        id: 1,
        role: "admin",
        address: ADMIN_ADDRESS,
        social_accounts: [],
        created_at: 0,
        updated_at: 0,
      },
      {
        id: 2,
        role: "verifier",
        address: VERIFIER_ADDRESS,
        social_accounts: [],
        created_at: 0,
        updated_at: 0,
      },
      {
        id: 3,
        role: "creator",
        address: CREATOR_ADDRESS,
        social_accounts: [],
        created_at: baseTime,
        updated_at: 0,
      },
      {
        id: 4,
        role: "investor",
        address: INVESTOR_01_ADDRESS,
        social_accounts: [],
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedOutput);
  });

  it("should find user by address", () => {
    const baseTime = Math.floor(Date.now() / 1000);

    const findUserByAddressInput = JSON.stringify({
      path: "user/address",
      data: {
        address: CREATOR_ADDRESS,
      },
    });

    const reports = machine.inspect(Buffer.from(findUserByAddressInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));
    const expectedOutput = {
      id: 3,
      role: "creator",
      address: CREATOR_ADDRESS,
      social_accounts: [],
      created_at: baseTime,
      updated_at: 0,
    };
    expect(output).toEqual(expectedOutput);
  });

  it("should delete user", () => {
    const baseTime = Math.floor(Date.now() / 1000);

    const deleteUserInput = JSON.stringify({
      path: "user/admin/delete",
      data: {
        address: INVESTOR_01_ADDRESS,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(deleteUserInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);

    const expectedNoticePayload = `user deleted - {"address":"${INVESTOR_01_ADDRESS}"}`;
    const expectedOutput = encodeNoticeOutput({
      payload: stringToHex(expectedNoticePayload),
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  afterAll(() => {
    machine.shutdown();
  });
});
