const dotEnv =  require('dotenv');
const joi = require('joi');

const envVarsSchema = joi.object().keys({
    PORT_HTTP: joi.number().default(8082),
    POSTGRES_USERNAME: joi.string(),
    POSTGRES_PASSWORD: joi.string(),
    POSTGRES_LOCALHOST: joi.string().default('localhost'),
    POSTGRES_PORT: joi.number().default(5432),
    POSTGRES_DATABASE: joi.string(),
    POSTGRES_URI: joi.string().required().description('postgree url db'),
    LOG_LEVEL: joi.string().valid('error', 'warn', 'info', 'debug').default('info'),
    LOG_FILE: joi.string().default('app.log'),
    REDIS_HOST: joi.string().default('localhost'),
    REDIS_PORT: joi.string().default('6379'),
}
).unknown();

function createConfig(path){
    dotEnv.config({ path });

    const {value:envVars, error } = envVarsSchema.prefs({errors : { label: 'key'}}).validate(process.env);
    if (error) {
        throw new Error(`Config validation error: ${error.message}`);
    }

    return {
        port: envVars.PORT_HTTP,
        postgreeURL: {
            username: envVars.POSTGRES_USERNAME,
            password: envVars.POSTGRES_PASSWORD,
            hostname: envVars.POSTGRES_HOSTNAME,
            port: envVars.POSTGRES_PORT,
            db: envVars.POSTGRES_DATABASE,
            url : envVars.POSTGRES_URI
        },
        log: {
            level: envVars.LOG_LEVEL,
            file: envVars.LOG_FILE,
        },
        logLevel: envVars.LOG_LEVEL,
        redis: {
            hostname : envVars.REDIS_HOST,
            port: envVars.REDIS_PORT,
        },
    };
};

module.exports = {
    createConfig,
};