import { bytesToHex, encodeFunctionData, getAddress } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import { encodeAdvanceInput, encodeDelegateCallVoucherOutput } from "./encoder";
import {
  createMachine,
  ADMIN_ADDRESS,
  EMERGENCY_WITHDRAW_ADDRESS,
} from "./helpers";
import { emergencyWithdrawAbi } from "../../contracts";

describe("Emergency Tests", () => {
  const machine = createMachine();

  it("should emergency erc20 withdraw", () => {
    const to = getAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955");
    const token = getAddress("0xfafafafafafafafafafafafafafafafafafafafa");

    const emergencyERC20WithdrawInput = JSON.stringify({
      path: "user/admin/emergency-erc20-withdraw",
      data: {
        to: to,
        token: token,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        payload: `0x${Buffer.from(emergencyERC20WithdrawInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);
    const delegateCallPayload = encodeFunctionData({
      abi: emergencyWithdrawAbi,
      functionName: "emergencyERC20Withdraw",
      args: [token, to],
    });
    const expectedOutput = encodeDelegateCallVoucherOutput({
      destination: EMERGENCY_WITHDRAW_ADDRESS,
      payload: delegateCallPayload,
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  it("should emergency ether withdraw", () => {
    const to = getAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955");

    const emergencyEtherWithdrawInput = JSON.stringify({
      path: "user/admin/emergency-ether-withdraw",
      data: {
        to: to,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        payload: `0x${Buffer.from(emergencyEtherWithdrawInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);
    const delegateCallPayload = encodeFunctionData({
      abi: emergencyWithdrawAbi,
      functionName: "emergencyETHWithdraw",
      args: [to],
    });
    const expectedOutput = encodeDelegateCallVoucherOutput({
      destination: EMERGENCY_WITHDRAW_ADDRESS,
      payload: delegateCallPayload,
    });
    expect(bytesToHex(outputs[0])).toBe(expectedOutput);
  });

  afterAll(() => {
    machine.shutdown();
  });
});
