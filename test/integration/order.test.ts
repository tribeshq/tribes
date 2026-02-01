import { bytesToHex, getAddress, stringToHex } from 'viem'
import { afterAll, describe, expect, it } from 'vitest'
import { encodeAdvanceInput, encodeERC20Deposit, encodeNoticeOutput } from './encoder'
import { createMachine } from './helpers'

const ADMIN_ADDRESS = getAddress('0xD554153658E8D466428Fa48487f5aba18dF5E628')

const VERIFIER_ADDRESS = getAddress('0xc2D8eb4a934AEc7268E414a3Fa3D20E0572d714b')

const TOKEN_ADDRESS = getAddress('0x0000000000000000000000000000000000000009')

const CREATOR_ADDRESS = getAddress('0x0000000000000000000000000000000000000007')

const COLLATERAL = getAddress('0x0000000000000000000000000000000000000008')

const APPLICATION_ADDRESS = getAddress('0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e')

const ERC20_PORTAL_ADDRESS = getAddress('0xc700D6aDd016eECd59d989C028214Eaa0fCC0051')

const INVESTOR_01_ADDRESS = getAddress('0x0000000000000000000000000000000000000001')

describe('Order Tests', () => {
  const machine = createMachine()

  const baseTime = Math.floor(Date.now() / 1000)
  const closesAt = baseTime + 5
  const maturityAt = baseTime + 10

  it('should create order', () => {
    // Setup: create creator user
    const createUserInput = JSON.stringify({
      path: 'user/admin/create',
      data: {
        address: CREATOR_ADDRESS,
        role: 'creator',
      },
    })

    machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createUserInput).toString('hex')}`,
      }),
      { collect: true }
    )

    // Setup: create social account
    const createSocialAccountInput = JSON.stringify({
      path: 'social/verifier/create',
      data: {
        address: CREATOR_ADDRESS,
        username: 'test',
        platform: 'twitter',
      },
    })

    machine.advance(
      encodeAdvanceInput({
        msgSender: VERIFIER_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createSocialAccountInput).toString('hex')}`,
      }),
      { collect: true }
    )

    // Setup: create issuance
    const createIssuanceInput = JSON.stringify({
      path: 'issuance/creator/create',
      data: {
        title: 'test',
        description: 'testtesttesttesttest',
        promotion: 'testtesttesttesttest',
        token: TOKEN_ADDRESS,
        max_interest_rate: '10',
        debt_issued: '100000',
        closes_at: closesAt,
        maturity_at: maturityAt,
      },
    })

    const issuanceErc20DepositPayload = encodeERC20Deposit({
      tokenAddress: COLLATERAL,
      sender: CREATOR_ADDRESS,
      amount: 10000n,
      execLayerData: `0x${Buffer.from(createIssuanceInput).toString('hex')}`,
    })

    machine.advance(
      encodeAdvanceInput({
        appContract: APPLICATION_ADDRESS,
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: issuanceErc20DepositPayload,
        index: 0n,
      }),
      { collect: true }
    )

    // Setup: create investor user
    const createInvestorInput = JSON.stringify({
      path: 'user/admin/create',
      data: {
        address: INVESTOR_01_ADDRESS,
        role: 'investor',
      },
    })

    machine.advance(
      encodeAdvanceInput({
        msgSender: ADMIN_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(createInvestorInput).toString('hex')}`,
      }),
      { collect: true }
    )

    // Create order
    const createOrderInput = JSON.stringify({
      path: 'order/create',
      data: {
        issuance_id: 1,
        interest_rate: '9',
      },
    })

    const orderErc20DepositPayload = encodeERC20Deposit({
      tokenAddress: TOKEN_ADDRESS,
      sender: INVESTOR_01_ADDRESS,
      amount: 10000n,
      execLayerData: `0x${Buffer.from(createOrderInput).toString('hex')}`,
    })

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: ERC20_PORTAL_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: orderErc20DepositPayload,
      }),
      { collect: true }
    )

    expect(outputs.length).toBeGreaterThanOrEqual(1)

    const expectedCreateOrderNoticeOutput = encodeNoticeOutput({
      payload: stringToHex(
        `order created - {"id":1,"issuance_id":1,"investor":{"id":4,"role":"investor","address":"${INVESTOR_01_ADDRESS}","social_accounts":[],"created_at":${baseTime},"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":${baseTime}}`
      ),
    })
    expect(bytesToHex(outputs[0])).toBe(expectedCreateOrderNoticeOutput)
  })

  it('should find all orders', () => {
    const findAllOrdersInput = JSON.stringify({
      path: 'order',
    })

    const reports = machine.inspect(Buffer.from(findAllOrdersInput), {
      collect: true,
    })

    expect(reports.length).toBe(1)
    const output = JSON.parse(Buffer.from(reports[0]).toString('utf-8'))

    const expectedFindAllOrdersOutput = [
      {
        id: 1,
        issuance_id: 1,
        state: 'pending',
        investor: {
          id: 4,
          role: 'investor',
          address: INVESTOR_01_ADDRESS,
          social_accounts: [],
          created_at: baseTime,
          updated_at: 0,
        },
        amount: '10000',
        interest_rate: '9',
        created_at: baseTime,
        updated_at: 0,
      },
    ]
    expect(output).toEqual(expectedFindAllOrdersOutput)
  })

  it('should find order by id', () => {
    const findOrderByIdInput = JSON.stringify({
      path: 'order/id',
      data: {
        id: 1,
      },
    })

    const reports = machine.inspect(Buffer.from(findOrderByIdInput), {
      collect: true,
    })

    expect(reports.length).toBe(1)
    const output = JSON.parse(Buffer.from(reports[0]).toString('utf-8'))

    const expectedFindOrderByIdOutput = {
      id: 1,
      issuance_id: 1,
      state: 'pending',
      investor: {
        id: 4,
        role: 'investor',
        address: INVESTOR_01_ADDRESS,
        social_accounts: [],
        created_at: baseTime,
        updated_at: 0,
      },
      amount: '10000',
      interest_rate: '9',
      created_at: baseTime,
      updated_at: 0,
    }
    expect(output).toEqual(expectedFindOrderByIdOutput)
  })

  it('should find orders by issuance id', () => {
    const findOrdersByIssuanceInput = JSON.stringify({
      path: 'order/issuance',
      data: {
        issuance_id: 1,
      },
    })

    const reports = machine.inspect(Buffer.from(findOrdersByIssuanceInput), {
      collect: true,
    })

    expect(reports.length).toBe(1)
    const output = JSON.parse(Buffer.from(reports[0]).toString('utf-8'))

    const expectedFindOrdersByIssuanceOutput = [
      {
        id: 1,
        issuance_id: 1,
        state: 'pending',
        investor: {
          id: 4,
          role: 'investor',
          address: INVESTOR_01_ADDRESS,
          social_accounts: [],
          created_at: baseTime,
          updated_at: 0,
        },
        amount: '10000',
        interest_rate: '9',
        created_at: baseTime,
        updated_at: 0,
      },
    ]
    expect(output).toEqual(expectedFindOrdersByIssuanceOutput)
  })

  it('should find orders by investor address', () => {
    const findOrdersByInvestorInput = JSON.stringify({
      path: 'order/investor',
      data: {
        investor_address: INVESTOR_01_ADDRESS,
      },
    })

    const reports = machine.inspect(Buffer.from(findOrdersByInvestorInput), {
      collect: true,
    })

    expect(reports.length).toBe(1)
    const output = JSON.parse(Buffer.from(reports[0]).toString('utf-8'))

    const expectedFindOrdersByInvestorOutput = [
      {
        id: 1,
        issuance_id: 1,
        state: 'pending',
        investor: {
          id: 4,
          role: 'investor',
          address: INVESTOR_01_ADDRESS,
          social_accounts: [],
          created_at: baseTime,
          updated_at: 0,
        },
        amount: '10000',
        interest_rate: '9',
        created_at: baseTime,
        updated_at: 0,
      },
    ]
    expect(output).toEqual(expectedFindOrdersByInvestorOutput)
  })

  it('should cancel order', () => {
    const cancelOrderInput = JSON.stringify({
      path: 'order/cancel',
      data: {
        id: 1,
      },
    })

    const { outputs } = machine.advance(
      encodeAdvanceInput({
        msgSender: INVESTOR_01_ADDRESS,
        blockTimestamp: BigInt(baseTime),
        payload: `0x${Buffer.from(cancelOrderInput).toString('hex')}`,
      }),
      { collect: true }
    )

    expect(outputs.length).toBe(1)
    const noticePayload = Buffer.from(outputs[0]).toString('utf-8')
    expect(noticePayload).toContain('order canceled')
  })

  afterAll(() => {
    machine.shutdown()
  })
})
