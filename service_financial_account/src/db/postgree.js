const { Sequelize } = require('sequelize');
const config = require('../config/config');

const cfg = config.createConfig('.env');

const sequelize = new Sequelize({
    dialect: 'postgres',
    host: cfg.postgreeURL.hostname,
    port: cfg.postgreeURL.port,
    username: cfg.postgreeURL.username,
    password: process.env.POSTGRES_PASSWORD,
    database: process.env.POSTGRES_DATABASE,
    logging: false,
});

console.log("pass = ", process.env.POSTGRES_PASSWORD);
module.exports = sequelize;