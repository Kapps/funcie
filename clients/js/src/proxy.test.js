const { lambdaProxy, lambdaProxyWithConfig } = require('./proxy');
const { sendMessage } = require('./bastionClient');
const { loadConfig } = require('./config');
const { invokeLambda } = require('./utils');

jest.mock('./config');
jest.mock('./bastionClient');
jest.mock('./utils');

describe('lambdaProxy[WithConfig]', () => {
  const mockEvent = { some: 'event' };
  const mockContext = { some: 'context' };
  const mockHandler = jest.fn();
  const mockConfig = {
    ApplicationId: 'app-id',
    ServerBastionEndpoint: 'some-endpoint',
  };

  beforeEach(() => {
    jest.clearAllMocks();
    loadConfig.mockReturnValue(mockConfig);
  });

  it('should forward from lambdaProxy to lambdaProxyWithConfig', async () => {
    const mockResponse = { data: { body: 'some-data' } };
    sendMessage.mockResolvedValue(mockResponse);

    const result = await lambdaProxy('app-id', mockHandler)(mockEvent, mockContext);
    expect(loadConfig).toHaveBeenCalledWith('app-id');
    expect(result).toBe('some-data');
  });

  it('should forward response data body if no error occurs', async () => {
    const mockResponse = { data: { body: 'some-data' } };
    sendMessage.mockResolvedValue(mockResponse);

    const result = await lambdaProxyWithConfig(mockConfig, mockHandler)(mockEvent, mockContext);
    expect(result).toBe('some-data');
  });

  it('should invokeLambda directly if sendMessage throws an error', async () => {
    sendMessage.mockRejectedValue(new Error('some error'));
    invokeLambda.mockResolvedValue('direct-response');

    const result = await lambdaProxyWithConfig(mockConfig, mockHandler)(mockEvent, mockContext);
    expect(result).toBe('direct-response');
  });

  it('should invokeLambda directly if no consumer is active', async () => {
    const mockResponse = { error: { message: 'no consumer is active on this tunnel' } };
    sendMessage.mockResolvedValue(mockResponse);
    invokeLambda.mockResolvedValue('direct-response');

    const result = await lambdaProxyWithConfig(mockConfig, mockHandler)(mockEvent, mockContext);
    expect(result).toBe('direct-response');
  });

  it('should invokeLambda directly if application is not found', async () => {
    const mockResponse = { error: 'application not found' };
    sendMessage.mockResolvedValue(mockResponse);
    invokeLambda.mockResolvedValue('direct-response');

    const result = await lambdaProxyWithConfig(mockConfig, mockHandler)(mockEvent, mockContext);
    expect(result).toBe('direct-response');
  });

  it('should throw an error if any other error message is received from bastion', async () => {
    const mockResponse = { error: { message: 'some-other-error' } };
    sendMessage.mockResolvedValue(mockResponse);

    await expect(lambdaProxyWithConfig(mockConfig, mockHandler)(mockEvent, mockContext))
        .rejects.toThrow('some-other-error');
  });
});
