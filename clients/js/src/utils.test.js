const { invokeLambda } = require('./utils');

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
