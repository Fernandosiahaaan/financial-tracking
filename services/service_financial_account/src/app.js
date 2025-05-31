// Import modul Express
const express = require('express');
const accountControllers = require('./controllers/financeAccountControllers');
const config = require('./config/config');
const logger = require('./config/logger'); 
const router = require('./routes/accountRoutes');
const { error } = require('winston');

async function run() {
  const app = express();
  const cfg = config.createConfig('.env');

  logger.init({
    logLevel: cfg.log.level,
    logFile: cfg.log.file,
    isFile: false, 
  });

  app.use(express.json());
  app.use('', router);
  app.use((req,res) => { res.status(404).json({error: true, message: 'route not found'}) });
  app.listen(cfg.port, () => {
    logger.info(`🌐 http://localhost:${cfg.port}`);
  });
}

run();
