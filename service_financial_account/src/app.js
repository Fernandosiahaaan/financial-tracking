// Import modul Express
const express = require('express');
const accountControllers = require('./controllers/accountControllers');
const config = require('./config/config');
const logger = require('./config/logger'); 
const router = require('./routes/accountRoutes');

async function run() {
  const app = express();
  const cfg = config.createConfig('.env');

  logger.init({
    logLevel: cfg.log.level,
    logFile: cfg.log.file,
    isFile: false, 
  });

  app.use('', router);
  app.listen(cfg.port, () => {
    logger.info(`🌐 http://localhost:${cfg.port}`);
  });
}

run();
