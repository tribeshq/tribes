import { bytesToHex, stringToHex } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import {
  encodeAdvanceInput,
  encodeErc20Deposit,
  encodeNoticeOutput,
} from "./encoder";
import {
  createMachine,
  ADMIN_ADDRESS,
  VERIFIER_ADDRESS,
  CREATOR_ADDRESS,
  TOKEN_ADDRESS,
  COLLATERAL,
  APPLICATION_ADDRESS,
  ERC20_PORTAL_ADDRESS,
  INVESTOR_01_ADDRESS,
  setupTimeValues,
} from "./helpers";

describe("Order Tests", () => {
  const machine = createMachine();

  const { baseTime, closesAt, maturityAt } = setupTimeValues();

  it("should create order", () => {
    // Setup: create creator user
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

    // Setup: create social account
    const createSocialAccountInput = JSON.stringify({
      path: "social/verifier/create",
      data: {
        address: CREATOR_ADDRESS,
        username: "test",
        platform: "twitter",
      },
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: VERIFIER_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createSocialAccountInput).toString("hex")}`,
      }),
      { collect: true },
    );

    // Setup: create campaign
    const createCampaignInput = JSON.stringify({
      path: "campaign/creator/create",
      data: {
        title: "test",
        description: "testtesttesttesttest",
        promotion: "testtesttesttesttest",
        token: TOKEN_ADDRESS,
        max_interest_rate: "10",
        debt_issued: "100000",
        closes_at: closesAt,
        maturity_at: maturityAt,
      },
    });

    const campaignErc20DepositPayload = encodeErc20Deposit({
      tokenAddress: COLLATERAL,
      sender: CREATOR_ADDRESS,
      amount: 10000n,
      execLayerData: `0x${Buffer.from(createCampaignInput).toString("hex")}`,
    });

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: campaignErc20DepositPayload,
        index: 0n,
      }),
      { collect: true },
    );

    // Setup: create investor user
    const createInvestorInput = JSON.stringify({
      path: "user/admin/create",
      data: {
        address: INVESTOR_01_ADDRESS,
        role: "investor",
      },
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createInvestorInput).toString("hex")}`,
      }),
      { collect: true },
    );

    // Create order
    const createOrderInput = JSON.stringify({
      path: "order/create",
      data: {
        campaign_id: 1,
        interest_rate: "9",
      },
    });

    const orderErc20DepositPayload = encodeErc20Deposit({
      tokenAddress: TOKEN_ADDRESS,
      sender: INVESTOR_01_ADDRESS,
      amount: 10000n,
      execLayerData: `0x${Buffer.from(createOrderInput).toString("hex")}`,
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: orderErc20DepositPayload,
      }),
      { collect: true },
    );

    expect(outputs.length).toBeGreaterThanOrEqual(1);

    const expectedCreateOrderNoticeOutput = encodeNoticeOutput({
      payload: stringToHex(
        `order created - {"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"${INVESTOR_01_ADDRESS}","social_accounts":[],"created_at":${baseTime},"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":${baseTime}}`,
      ),
    });
    expect(bytesToHex(outputs[0])).toBe(expectedCreateOrderNoticeOutput);
  });

  it("should find all orders", () => {
    const findAllOrdersInput = JSON.stringify({
      path: "order",
    });

    const reports = machine.inspect(Buffer.from(findAllOrdersInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const expectedFindAllOrdersOutput = [
      {
        id: 1,
        campaign_id: 1,
        state: "pending",
        investor: {
          id: 4,
          role: "investor",
          address: INVESTOR_01_ADDRESS,
          created_at: baseTime,
        },
        amount: "10000",
        interest_rate: "9",
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedFindAllOrdersOutput);
  });

  it("should find order by id", () => {
    const findOrderByIdInput = JSON.stringify({
      path: "order/id",
      data: {
        id: 1,
      },
    });

    const reports = machine.inspect(Buffer.from(findOrderByIdInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const expectedFindOrderByIdOutput = {
      id: 1,
      campaign_id: 1,
      state: "pending",
      investor: {
        id: 4,
        role: "investor",
        address: INVESTOR_01_ADDRESS,
        created_at: baseTime,
      },
      amount: "10000",
      interest_rate: "9",
      created_at: baseTime,
      updated_at: 0,
    };
    expect(output).toEqual(expectedFindOrderByIdOutput);
  });

  it("should find orders by campaign id", () => {
    const findOrdersByCampaignInput = JSON.stringify({
      path: "order/campaign",
      data: {
        campaign_id: 1,
      },
    });

    const reports = machine.inspect(Buffer.from(findOrdersByCampaignInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const expectedFindOrdersByCampaignOutput = [
      {
        id: 1,
        campaign_id: 1,
        state: "pending",
        investor: {
          id: 4,
          role: "investor",
          address: INVESTOR_01_ADDRESS,
          created_at: baseTime,
        },
        amount: "10000",
        interest_rate: "9",
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedFindOrdersByCampaignOutput);
  });

  it("should find orders by investor address", () => {
    const findOrdersByInvestorInput = JSON.stringify({
      path: "order/investor",
      data: {
        investor_address: INVESTOR_01_ADDRESS,
      },
    });

    const reports = machine.inspect(Buffer.from(findOrdersByInvestorInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const expectedFindOrdersByInvestorOutput = [
      {
        id: 1,
        campaign_id: 1,
        state: "pending",
        investor: {
          id: 4,
          role: "investor",
          address: INVESTOR_01_ADDRESS,
          created_at: baseTime,
        },
        amount: "10000",
        interest_rate: "9",
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedFindOrdersByInvestorOutput);
  });

  it("should cancel order", () => {
    const cancelOrderInput = JSON.stringify({
      path: "order/cancel",
      data: {
        id: 1,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: INVESTOR_01_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(cancelOrderInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBe(1);
    const noticePayload = Buffer.from(outputs[0]).toString("utf-8");
    expect(noticePayload).toContain("order canceled");
  });

  afterAll(() => {
    machine.shutdown();
  });
});
