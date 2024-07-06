const url = require('url');
const process = require('process');
const SSM = require('@aws-sdk/client-ssm');

const ssmClient = new SSM.SSMClient();

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
 * Returns a new FuncieConfig instance based on a combination of environment variables and SSM parameters.
 * The following variables are used:
 *	- FUNCIE_CLIENT_BASTION_ENDPOINT (optional; defaults to http://127.0.0.1:24193)
 *	- FUNCIE_SERVER_BASTION_ENDPOINT -> /funcie/<env>/bastion_host (required)
 *	- FUNCIE_LISTEN_ADDRESS (optional; defaults to localhost on a random port)
 *
 * @param {string} applicationId - The application ID.
 * @param {string?} env - The deployment environment, or "default" if not specified.
 * @returns {Promise<FuncieConfig>}
 */
async function loadConfig(applicationId, env) {
    if (!env) {
        env = "default";
    }

    let serverEndpoint = process.env.FUNCIE_SERVER_BASTION_ENDPOINT;
    if (!serverEndpoint) {
        const bastionHost = await loadSSMParameter(env, 'bastion_host');
        serverEndpoint = `http://${bastionHost}:8082/dispatch`;
    }

    return new FuncieConfig(
        optionalUrlEnv("FUNCIE_CLIENT_BASTION_ENDPOINT", "http://127.0.0.1:24193"),
        new url.URL(serverEndpoint),
        optionalUrlEnv("FUNCIE_LISTEN_ADDRESS", "http://0.0.0.0:0"),
        applicationId,
    );
}

/**
 * Loads a parameter from AWS SSM Parameter Store.
 *
 * @param {string} env - The environment name.
 * @param {string} name - The parameter name.
 * @returns {Promise<string>}
 */
async function loadSSMParameter(env, name) {
    const path = `/funcie/${env}/${name}`;
    try {
        const command = new SSM.GetParameterCommand({
            Name: path,
        });
        const result = await ssmClient.send(command);
        return result.Parameter.Value;
    } catch (error) {
        throw new Error(`failed to load SSM parameter ${path}: ${error.message}`);
    }
}

function requireUrlEnv(name, purpose) {
    const value = requiredEnv(name, purpose);
    if (!value) {
        return null;
    }

    return new url.URL(value);
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
    loadConfig,
    getConfigPurpose,
};