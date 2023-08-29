module.exports = {
  // Tells Jest to look for test files with any name ending in .test.js
  testMatch: ['**/*.test.js'],

  // The root directory that Jest should scan for tests and modules within
  testPathIgnorePatterns: ['/node_modules/'],

  // Indicates whether each individual test should be reported during the run
  verbose: true,
};
