import { defineConfig } from '@wagmi/cli'
import { foundry } from '@wagmi/cli/plugins'

export default defineConfig({
  out: './contracts.ts',
  plugins: [
    foundry({
      project: 'node_modules/@cartesi/rollups',
      forge: {
        build: false,
      },
    }),
    foundry({
      project: './contracts',
      include: ['SafeEmergencyWithdraw.sol/**', 'BadgeFactory.sol/**'],
      forge: {
        build: false,
      },
    }),
  ],
})
