import eslint from '@eslint/js'
import tseslint from 'typescript-eslint'

export default tseslint.defineConfig(eslint.configs.recommended, ...tseslint.configs.recommended, {
  ignores: ['node_modules/', 'dist/', 'contracts/', 'out/'],
})
