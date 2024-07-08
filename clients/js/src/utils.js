const invokeLambda = (handler, event, context) => {
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

const debug = (message) => {
    if (process.env.FUNCIE_DEBUG) {
        console.log(message);
    }
}

const info = (message) => {
    if (!process.env.FUNCIE_QUIET) {
        console.log(message);
    }
}

const error = (message) => {
    console.error(message);
}

module.exports = {
    invokeLambda,
    info,
    error,
    debug,
};
