const winston = require('winston');
const { combine, timestamp, printf, label } = winston.format;

const logFormat = printf(({ timestamp, label, level, message }) => {
    return `${timestamp} [${label}] ${level} : ${message}`;
});

const logger = winston.createLogger({
    format: combine(
        label({ label: 'financial-account-service' }),
        timestamp(),
        logFormat
    ),
    transports: [],
});
module.exports = logger;

function init({ logLevel = 'info', logFile = 'app.log', isFile = false }) {
    logger.add(
        new winston.transports.Console({
            level: logLevel,
        })
    );

    if (isFile) {
        logger.add(
            new winston.transports.File({
                level: logLevel,
                filename: logFile,
            })
        );
    }
}

module.exports.init = init;
