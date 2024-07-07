export const invokeLambda = (handler, event, context) => {
    // Callback-based handlers are supported, but we wrap them in a promise
    if (handler.length >= 3) {
        return new Promise((resolve, reject) => {
            handler(event, context, (err, data) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(data);
                }
            });
        });
    }

    // Async handlers are just invoked directly.
    return handler(event, context);
}

export const info = (message) => {
    if (process.env.FUNCIE_DEBUG) {
        console.log(message);
    }
}

export const error = (message) => {
    console.error(message);
}
