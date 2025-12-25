import { bytesToHex, stringToHex } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import { encodeAdvanceInput, encodeNoticeOutput } from "./encoder";
import {
  createMachine,
  ADMIN_ADDRESS,
  VERIFIER_ADDRESS,
  CREATOR_ADDRESS,
} from "./helpers";

describe("Social Account Tests", () => {
  const machine = createMachine();

  const baseTime = Math.floor(Date.now() / 1000);

  it("should create social account", () => {
    const createUserInput = JSON.stringify({
      path: "user/admin/create",
      data: {
        address: CREATOR_ADDRESS,
        role: "creator",
      },
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createUserInput).toString("hex")}`,
      }),
      { collect: true },
    );

    const createSocialAccountInput = JSON.stringify({
      path: "social/verifier/create",
      data: {
        address: CREATOR_ADDRESS,
        username: "test",
        platform: "twitter",
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: VERIFIER_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createSocialAccountInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);

    const expectedNoticePayload = `social account created - {"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":${baseTime}}`;
    const expectedOutput = encodeNoticeOutput({
      payload: stringToHex(expectedNoticePayload),
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  it("should find social account by id", () => {
    const findSocialAccountInput = JSON.stringify({
      path: "social/id",
      data: {
        social_account_id: 1,
      },
    });

    const reports = machine.inspect(Buffer.from(findSocialAccountInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));
    const expectedOutput = {
      id: 1,
      user_id: 3,
      username: "test",
      platform: "twitter",
      created_at: baseTime,
      updated_at: 0,
    };
    expect(output).toEqual(expectedOutput);
  });

  it("should find social account by user id", () => {
    const findSocialAccountInput = JSON.stringify({
      path: "social/user/id",
      data: {
        user_id: 3,
      },
    });

    const reports = machine.inspect(Buffer.from(findSocialAccountInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));
    const expectedFindSocialAccountsByUserIdOutput = [
      {
        id: 1,
        user_id: 3,
        username: "test",
        platform: "twitter",
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedFindSocialAccountsByUserIdOutput);
  });

  it("should delete social account", () => {
    const deleteSocialAccountInput = JSON.stringify({
      path: "social/admin/delete",
      data: {
        social_account_id: 1,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(deleteSocialAccountInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);

    const expectedNoticePayload = `social account deleted - {"social_account_id":1}`;
    const expectedOutput = encodeNoticeOutput({
      payload: stringToHex(expectedNoticePayload),
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  afterAll(() => {
    machine.shutdown();
  });
});
