import { bytesToHex, encodeFunctionData, padHex, stringToHex } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import {
  encodeAdvanceInput,
  encodeErc20Deposit,
  encodeVoucherOutput,
  encodeNoticeOutput,
} from "./encoder";
import {
  createMachine,
  computeBadgeAddress,
  ADMIN_ADDRESS,
  VERIFIER_ADDRESS,
  CREATOR_ADDRESS,
  TOKEN_ADDRESS,
  FACTORY_ADDRESS,
  COLLATERAL,
  APPLICATION_ADDRESS,
  ERC20_PORTAL_ADDRESS,
  INVESTOR_01_ADDRESS,
  INVESTOR_02_ADDRESS,
  INVESTOR_03_ADDRESS,
  INVESTOR_04_ADDRESS,
  INVESTOR_05_ADDRESS,
  setupTimeValues,
} from "./helpers";
import { badgeFactoryAbi } from "../../contracts";

describe("Campaign Tests", () => {
  const machine = createMachine();

  const { baseTime, closesAt, maturityAt } = setupTimeValues();

  it("should create campaign", () => {
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

    machine.advance(
      encodeAdvanceInput({
        msgSender: VERIFIER_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createSocialAccountInput).toString("hex")}`,
      }),
      { collect: true },
    );

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

    const erc20DepositPayload = encodeErc20Deposit({
      tokenAddress: COLLATERAL,
      sender: CREATOR_ADDRESS,
      amount: 10000n,
      execLayerData: `0x${Buffer.from(createCampaignInput).toString("hex")}`,
    });

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: erc20DepositPayload,
        index: metadataIndex,
      }),
      { collect: true },
    );

    expect(outputs.length).toBeGreaterThanOrEqual(2);

    // Verify voucher for badge creation
    const expectedVoucherPayload = encodeFunctionData({
      abi: badgeFactoryAbi,
      functionName: "newBadge",
      args: [APPLICATION_ADDRESS, salt],
    });
    const expectedVoucher = encodeVoucherOutput({
      destination: FACTORY_ADDRESS,
      value: 0n,
      payload: expectedVoucherPayload,
    });
    expect(bytesToHex(outputs[0])).toBe(expectedVoucher);

    // Verify notice for campaign creation
    const expectedCreateCampaignNoticeOutput = encodeNoticeOutput({
      payload: stringToHex(
        `campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"${TOKEN_ADDRESS.toLowerCase()}","creator":{"id":3,"role":"creator","address":"${CREATOR_ADDRESS}","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":${baseTime}}],"created_at":${baseTime},"updated_at":0},"collateral":"${COLLATERAL.toLowerCase()}","collateral_amount":"10000","badge_address":"${badgeAddress}","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":${baseTime},"closes_at":${closesAt},"maturity_at":${maturityAt}}`,
      ),
    });
    expect(bytesToHex(outputs[1])).toBe(expectedCreateCampaignNoticeOutput);
  });

  it("should find all campaigns", () => {
    const findAllCampaignsInput = JSON.stringify({
      path: "campaign",
    });

    const reports = machine.inspect(Buffer.from(findAllCampaignsInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    const expectedFindAllCampaignsOutput = [
      {
        id: 1,
        title: "test",
        description: "testtesttesttesttest",
        promotion: "testtesttesttesttest",
        state: "ongoing",
        creator: {
          id: 3,
          role: "creator",
          address: CREATOR_ADDRESS,
          social_accounts: [
            {
              id: 1,
              user_id: 3,
              username: "test",
              platform: "twitter",
              created_at: baseTime,
            },
          ],
          created_at: baseTime,
          updated_at: 0,
        },
        token: TOKEN_ADDRESS.toLowerCase(),
        collateral: COLLATERAL.toLowerCase(),
        collateral_amount: "10000",
        badge_address: badgeAddress,
        max_interest_rate: "10",
        debt_issued: "100000",
        total_obligation: "0",
        total_raised: "0",
        orders: [],
        closes_at: closesAt,
        maturity_at: maturityAt,
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedFindAllCampaignsOutput);
  });

  it("should find campaign by id", () => {
    const findCampaignByIdInput = JSON.stringify({
      path: "campaign/id",
      data: {
        id: 1,
      },
    });

    const reports = machine.inspect(Buffer.from(findCampaignByIdInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    const expectedFindCampaignByIdOutput = {
      id: 1,
      title: "test",
      description: "testtesttesttesttest",
      promotion: "testtesttesttesttest",
      state: "ongoing",
      creator: {
        id: 3,
        role: "creator",
        address: CREATOR_ADDRESS,
        social_accounts: [
          {
            id: 1,
            user_id: 3,
            username: "test",
            platform: "twitter",
            created_at: baseTime,
          },
        ],
        created_at: baseTime,
        updated_at: 0,
      },
      token: TOKEN_ADDRESS.toLowerCase(),
      collateral: COLLATERAL.toLowerCase(),
      collateral_amount: "10000",
      badge_address: badgeAddress,
      max_interest_rate: "10",
      debt_issued: "100000",
      total_obligation: "0",
      total_raised: "0",
      orders: [],
      closes_at: closesAt,
      maturity_at: maturityAt,
      created_at: baseTime,
      updated_at: 0,
    };
    expect(output).toEqual(expectedFindCampaignByIdOutput);
  });

  it("should find campaign by creator address", () => {
    const findCampaignByCreatorInput = JSON.stringify({
      path: "campaign/creator",
      data: {
        creator: CREATOR_ADDRESS,
      },
    });

    const reports = machine.inspect(Buffer.from(findCampaignByCreatorInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    const expectedFindCampaignByCreatorOutput = [
      {
        id: 1,
        title: "test",
        description: "testtesttesttesttest",
        promotion: "testtesttesttesttest",
        state: "ongoing",
        creator: {
          id: 3,
          role: "creator",
          address: CREATOR_ADDRESS,
          social_accounts: [
            {
              id: 1,
              user_id: 3,
              username: "test",
              platform: "twitter",
              created_at: baseTime,
            },
          ],
          created_at: baseTime,
          updated_at: 0,
        },
        token: TOKEN_ADDRESS.toLowerCase(),
        collateral: COLLATERAL.toLowerCase(),
        collateral_amount: "10000",
        badge_address: badgeAddress,
        max_interest_rate: "10",
        debt_issued: "100000",
        total_obligation: "0",
        total_raised: "0",
        orders: [],
        closes_at: closesAt,
        maturity_at: maturityAt,
        created_at: baseTime,
        updated_at: 0,
      },
    ];
    expect(output).toEqual(expectedFindCampaignByCreatorOutput);
  });

  afterAll(() => {
    machine.shutdown();
  });
});

describe("Campaign Close Tests", () => {
  const machine = createMachine();

  const { baseTime, closesAt, maturityAt } = setupTimeValues();

  it("should close campaign after time passes", async () => {
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

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: campaignErc20DepositPayload,
        index: metadataIndex,
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

    // Setup: create order
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
      amount: 70000n,
      execLayerData: `0x${Buffer.from(createOrderInput).toString("hex")}`,
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: orderErc20DepositPayload,
      }),
      { collect: true },
    );

    await new Promise((resolve) => setTimeout(resolve, 6000));

    // Close campaign
    const anyone = INVESTOR_01_ADDRESS;
    const closeCampaignInput = JSON.stringify({
      path: "campaign/close",
      data: {
        creator_address: CREATOR_ADDRESS,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: anyone,
        blockTimestamp: BigInt(baseTime + 6),
        payload: `0x${Buffer.from(closeCampaignInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBeGreaterThanOrEqual(1);

    const expectedCloseCampaignNoticeOutput = encodeNoticeOutput({
      payload: stringToHex(
        `campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"${TOKEN_ADDRESS.toLowerCase()}","creator":{"id":3,"role":"creator","address":"${CREATOR_ADDRESS}","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":${baseTime}}],"created_at":${baseTime},"updated_at":0},"collateral":"${COLLATERAL.toLowerCase()}","collateral_amount":"10000","badge_address":"${badgeAddress}","debt_issued":"100000","max_interest_rate":"10","total_obligation":"76300","total_raised":"70000","state":"closed","orders":[{"id":1,"campaign_id":1,"investor":"${INVESTOR_01_ADDRESS}","amount":"70000","interest_rate":"9","state":"accepted","created_at":${baseTime},"updated_at":1}],"created_at":${baseTime},"closes_at":${closesAt},"maturity_at":${maturityAt},"updated_at":1}`,
      ),
    });
    expect(bytesToHex(outputs[outputs.length - 1])).toBe(
      expectedCloseCampaignNoticeOutput,
    );
  }, 15000);

  afterAll(() => {
    machine.shutdown();
  });
});

describe("Campaign Execute Collateral Tests", () => {
  const machine = createMachine();

  const { baseTime, closesAt, maturityAt } = setupTimeValues();

  it("should execute campaign collateral", async () => {
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
        appContract: APPLICATION_ADDRESS,
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
        appContract: APPLICATION_ADDRESS,
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

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: campaignErc20DepositPayload,
        index: metadataIndex,
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
        appContract: APPLICATION_ADDRESS,
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createInvestorInput).toString("hex")}`,
      }),
      { collect: true },
    );

    // Setup: create order
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
      amount: 70000n,
      execLayerData: `0x${Buffer.from(createOrderInput).toString("hex")}`,
    });

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: orderErc20DepositPayload,
      }),
      { collect: true },
    );

    await new Promise((resolve) => setTimeout(resolve, 6000));

    // Close campaign
    const anyone = INVESTOR_01_ADDRESS;
    const closeCampaignInput = JSON.stringify({
      path: "campaign/close",
      data: {
        creator_address: CREATOR_ADDRESS,
      },
    });

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: anyone,
        blockTimestamp: BigInt(baseTime + 6),
        payload: `0x${Buffer.from(closeCampaignInput).toString("hex")}`,
      }),
      { collect: true },
    );

    await new Promise((resolve) => setTimeout(resolve, 6000));

    // Execute campaign collateral
    const executeCampaignCollateralInput = JSON.stringify({
      path: "campaign/execute-collateral",
      data: {
        id: 1,
      },
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: CREATOR_ADDRESS,
        blockTimestamp: BigInt(baseTime + 11),
        payload: `0x${Buffer.from(executeCampaignCollateralInput).toString("hex")}`,
      }),
      { collect: true },
    );

    expect(outputs.length).toBeGreaterThanOrEqual(1);

    const expectedExecuteCampaignCollateralNoticeOutput = encodeNoticeOutput({
      payload: stringToHex(
        `campaign collateral executed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"${TOKEN_ADDRESS.toLowerCase()}","creator":{"id":3,"role":"creator","address":"${CREATOR_ADDRESS}","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":${baseTime}}],"created_at":${baseTime},"updated_at":0},"collateral":"${COLLATERAL.toLowerCase()}","collateral_amount":"10000","badge_address":"${badgeAddress}","debt_issued":"100000","max_interest_rate":"10","total_obligation":"76300","total_raised":"70000","state":"collateral_executed","orders":[{"id":1,"campaign_id":1,"investor":"${INVESTOR_01_ADDRESS}","amount":"70000","interest_rate":"9","state":"settled_by_collateral","created_at":${baseTime},"updated_at":1}],"created_at":${baseTime},"closes_at":${closesAt},"maturity_at":${maturityAt},"updated_at":1}`,
      ),
    });
    expect(bytesToHex(outputs[outputs.length - 1])).toBe(
      expectedExecuteCampaignCollateralNoticeOutput,
    );
  }, 15000);

  afterAll(() => {
    machine.shutdown();
  });
});

describe("Campaign Settle Tests", () => {
  const machine = createMachine();

  const { baseTime, closesAt, maturityAt } = setupTimeValues();

  it("should settle campaign", async () => {
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

    const metadataIndex = 0n;
    const salt = padHex(`0x${metadataIndex.toString(16)}`, { size: 32 });
    const badgeAddress = computeBadgeAddress(
      FACTORY_ADDRESS,
      salt,
      APPLICATION_ADDRESS,
    );

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: campaignErc20DepositPayload,
        index: metadataIndex,
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

    // Setup: create order
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
      amount: 70000n,
      execLayerData: `0x${Buffer.from(createOrderInput).toString("hex")}`,
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: orderErc20DepositPayload,
      }),
      { collect: true },
    );

    await new Promise((resolve) => setTimeout(resolve, 6000));

    // Close campaign
    const anyone = INVESTOR_01_ADDRESS;
    const closeCampaignInput = JSON.stringify({
      path: "campaign/close",
      data: {
        creator_address: CREATOR_ADDRESS,
      },
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: anyone,
        blockTimestamp: BigInt(baseTime + 6),
        payload: `0x${Buffer.from(closeCampaignInput).toString("hex")}`,
      }),
      { collect: true },
    );

    // Withdraw raised amount
    const withdrawRaisedAmountInput = JSON.stringify({
      path: "user/withdraw",
      data: {
        token: TOKEN_ADDRESS,
        amount: "70000",
      },
    });

    machine.advance(
      encodeAdvanceInput({
        msgSender: CREATOR_ADDRESS,
        blockTimestamp: BigInt(baseTime + 7),
        payload: `0x${Buffer.from(withdrawRaisedAmountInput).toString("hex")}`,
      }),
      { collect: true },
    );

    await new Promise((resolve) => setTimeout(resolve, 5000));

    // Settle campaign
    const settleCampaignInput = JSON.stringify({
      path: "campaign/creator/settle",
      data: {
        id: 1,
      },
    });

    const settleCampaignPayload = encodeErc20Deposit({
      tokenAddress: TOKEN_ADDRESS,
      sender: CREATOR_ADDRESS,
      amount: 76300n,
      execLayerData: `0x${Buffer.from(settleCampaignInput).toString("hex")}`,
    });

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime + 9),
        payload: settleCampaignPayload,
      }),
      { collect: true },
    );

    expect(outputs.length).toBeGreaterThanOrEqual(1);

    const expectedSettleCampaignNoticeOutput = encodeNoticeOutput({
      payload: stringToHex(
        `campaign settled - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"${TOKEN_ADDRESS.toLowerCase()}","creator":{"id":3,"role":"creator","address":"${CREATOR_ADDRESS}","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":${baseTime}}],"created_at":${baseTime},"updated_at":0},"collateral":"${COLLATERAL.toLowerCase()}","collateral_amount":"10000","badge_address":"${badgeAddress}","debt_issued":"100000","max_interest_rate":"10","total_obligation":"76300","total_raised":"70000","state":"settled","orders":[{"id":1,"campaign_id":1,"investor":"${INVESTOR_01_ADDRESS}","amount":"70000","interest_rate":"9","state":"settled","created_at":${baseTime},"updated_at":1}],"created_at":${baseTime},"closes_at":${closesAt},"maturity_at":${maturityAt},"updated_at":1}`,
      ),
    });
    expect(bytesToHex(outputs[outputs.length - 1])).toBe(
      expectedSettleCampaignNoticeOutput,
    );
  }, 20000);

  afterAll(() => {
    machine.shutdown();
  });
});

describe("Campaign with Multiple Investors Tests", () => {
  const machine = createMachine();

  const { baseTime, closesAt, maturityAt } = setupTimeValues();

  it("should find campaigns by investor address", async () => {
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

    // Setup: create multiple investor users
    const investors = [
      INVESTOR_01_ADDRESS,
      INVESTOR_02_ADDRESS,
      INVESTOR_03_ADDRESS,
      INVESTOR_04_ADDRESS,
      INVESTOR_05_ADDRESS,
    ];

    for (let i = 0; i < investors.length; i++) {
      const createInvestorInput = JSON.stringify({
        path: "user/admin/create",
        data: {
          address: investors[i],
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
    }

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
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: campaignErc20DepositPayload,
      }),
      { collect: true },
    );

    // Create orders for each investor
    for (let i = 0; i < investors.length; i++) {
      const createOrderInput = JSON.stringify({
        path: "order/create",
        data: {
          campaign_id: 1,
          interest_rate: String(10 - i),
        },
      });

      const orderErc20DepositPayload = encodeErc20Deposit({
        tokenAddress: TOKEN_ADDRESS,
        sender: investors[i],
        amount: BigInt(20000 * (i + 1)),
        execLayerData: `0x${Buffer.from(createOrderInput).toString("hex")}`,
      });

      machine.advance(
        encodeAdvanceInput({
          msgSender: ERC20_PORTAL_ADDRESS,
          blockTimestamp: BigInt(baseTime),
          payload: orderErc20DepositPayload,
        }),
        { collect: true },
      );
    }

    // Find campaigns by investor address
    const findCampaignsByInvestorInput = JSON.stringify({
      path: "campaign/investor",
      data: {
        investor_address: INVESTOR_01_ADDRESS,
      },
    });

    const reports = machine.inspect(Buffer.from(findCampaignsByInvestorInput), {
      collect: true,
    });

    expect(reports.length).toBe(1);
    const output = JSON.parse(Buffer.from(reports[0]).toString("utf-8"));

    // This test has complex nested structures. Verify key fields instead
    expect(output.length).toBe(1);
    expect(output[0].id).toBe(1);
    expect(output[0].title).toBe("test");
    expect(output[0].state).toBe("ongoing");
    expect(output[0].creator.address).toBe(CREATOR_ADDRESS);
    expect(output[0].token).toBe(TOKEN_ADDRESS.toLowerCase());
    expect(output[0].collateral).toBe(COLLATERAL.toLowerCase());
    expect(output[0].orders.length).toBe(5);
  });

  afterAll(() => {
    machine.shutdown();
  });
});
