import { rollups } from '@tuler/node-cartesi-machine'
import { encodeAbiParameters, getContractAddress, keccak256, concat, type Address, type Hex } from 'viem'
import badgeArtifact from '../../assets/artifacts/Badge.json'

export const createMachine = () => {
  return rollups('.cartesi/image', {
    runtimeConfig: { skip_root_hash_check: true },
  })
}

export const computeBadgeAddress = (factory: Address, salt: Hex, appContract: Address): Address => {
  const bytecode = badgeArtifact.bytecode as Hex
  const constructorArgs = encodeAbiParameters([{ type: 'address' }], [appContract])
  const initCodeHash = keccak256(concat([bytecode, constructorArgs]))
  return getContractAddress({
    bytecodeHash: initCodeHash,
    from: factory,
    opcode: 'CREATE2',
    salt,
  })
}
