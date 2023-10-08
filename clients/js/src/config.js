const url = require('url');
const process = require('process');

const CONFIG_PURPOSE_CLIENT = 'client';
const CONFIG_PURPOSE_SERVER = 'server';
const CONFIG_PURPOSE_ANY = 'any';



/**
 * Returns the current configuration purpose, which is either "client" or "server".
 */
function getConfigPurpose() {
    if (process.env.AWS_LAMBDA_FUNCTION_NAME) {
        return CONFIG_PURPOSE_SERVER;
    }
    return CONFIG_PURPOSE_CLIENT;
}

/**
 * Represents the configuration for a Funcie client or server application.
 * 
 * @typedef {Object} FuncieConfig
 * @property {URL} ClientBastionEndpoint
 * @property {URL} ServerBastionEndpoint
 * @property {string} ListenAddress
 * @property {string} ApplicationId
 */
class FuncieConfig {
    constructor(clientBastionEndpoint, serverBastionEndpoint, listenAddress, applicationId) {
        this.ClientBastionEndpoint = clientBastionEndpoint;
        this.ServerBastionEndpoint = serverBastionEndpoint;
        this.ListenAddress = listenAddress;
        this.ApplicationId = applicationId;
    }
}

/**
 * Returns a new FuncieConfig instance based on the following environment variables:
 *	- FUNCIE_APPLICATION_ID (required)
 *	- FUNCIE_CLIENT_BASTION_ENDPOINT (optional; defaults to http://127.0.0.1:24193)
 *	- FUNCIE_SERVER_BASTION_ENDPOINT (required for server)
 *	- FUNCIE_LISTEN_ADDRESS (optional; defaults to localhost on a random port)
 *
 * @returns {FuncieConfig}
 */
function loadConfigFromEnvironment() {
    return new FuncieConfig(
        optionalUrlEnv("FUNCIE_CLIENT_BASTION_ENDPOINT", "http://127.0.0.1:24193"),
        requireUrlEnv("FUNCIE_SERVER_BASTION_ENDPOINT", CONFIG_PURPOSE_SERVER),
        optionalUrlEnv("FUNCIE_LISTEN_ADDRESS", "http://0.0.0.0:0"),
        requiredEnv("FUNCIE_APPLICATION_ID", CONFIG_PURPOSE_ANY),
    );
}

function requireUrlEnv(name, purpose) {
    const value = requiredEnv(name, purpose);
    if (!value) {
        return null;
    }

    const parsedUrl = new url.URL(value);
    return parsedUrl;
}

function optionalUrlEnv(name, defaultValue) {
    const value = optionalEnv(name, defaultValue);
    if (!value) {
        return new url.URL(defaultValue);
    }

    return new url.URL(value);
}

function requiredEnv(name, purpose) {
    const value = process.env[name];
    if (!value) {
        const currPurpose = getConfigPurpose();
        const purposeMatches = purpose === CONFIG_PURPOSE_ANY || currPurpose === purpose;
        if (purposeMatches) {
            throw new Error(`required environment variable ${name} not set`);
        }
    }
    return value;
}

function optionalEnv(name, defaultValue) {
    const value = process.env[name];
    if (value === "") {
        return defaultValue;
    }
    return value;
}

module.exports = {
    FuncieConfig,
    loadConfigFromEnvironment,
    getConfigPurpose,
};
