module.exports = {
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
    ecmaFeatures: {
      jsx: true,
    },
  },
  settings: {
    react: {
      version: 'detect',
    },
  },
  extends: [
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
    'plugin:react/recommended',
    'prettier',
    'plugin:jsx-a11y/recommended',
  ],
  rules: {
    'indent': ['error', 2],
    '@typescript-eslint/quotes': [
      2,
      'single',
      {
        'avoidEscape': true
      }
    ],
  },
};
