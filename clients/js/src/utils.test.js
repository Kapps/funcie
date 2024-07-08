const { invokeLambda, debug, info} = require('./utils');

describe('invokeLambda', () => {

    it('should handle async handlers correctly', async () => {
        const mockHandler = jest.fn().mockResolvedValue('some data');
        const event = { some: 'event' };
        const context = { some: 'context' };

        const result = await invokeLambda(mockHandler, event, context);

        expect(result).toBe('some data');
        expect(mockHandler).toHaveBeenCalledWith(event, context);
    });

    it('should handle callback-based handlers correctly', async () => {
        const mockHandler = jest.fn((event, context, callback) => {
            callback(null, 'some data');
        });

        const event = { some: 'event' };
        const context = { some: 'context' };

        const result = await invokeLambda(mockHandler, event, context);

        expect(result).toBe('some data');
        expect(mockHandler).toHaveBeenCalledWith(event, context, expect.any(Function));
    });

    it('should reject for callback-based handlers if error', async () => {
        const mockHandler = jest.fn((event, context, callback) => {
            callback(new Error('some error'));
        });

        const event = { some: 'event' };
        const context = { some: 'context' };

        await expect(invokeLambda(mockHandler, event, context)).rejects.toThrow('some error');
    });
});

describe('debug', () => {
    let consoleLogSpy;

    beforeEach(() => {
        consoleLogSpy = jest.spyOn(console, 'log').mockImplementation();
    });

    afterEach(() => {
        consoleLogSpy.mockRestore();
    });

    it('should log if FUNCIE_DEBUG is set', () => {
        process.env.FUNCIE_DEBUG = 'true';

        const { debug } = require('./utils');
        debug('some message', 'some arg1', 'some arg2');

        expect(consoleLogSpy).toHaveBeenCalledWith('some message', 'some arg1', 'some arg2');
    });

    it('should not log if FUNCIE_DEBUG is not set', () => {
        process.env.FUNCIE_DEBUG = '';

        const { debug } = require('./utils');
        debug('some message', 'some arg1', 'some arg2');

        expect(consoleLogSpy).not.toHaveBeenCalled();
    });
});
