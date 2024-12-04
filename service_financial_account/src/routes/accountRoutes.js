const express = require('express');
const accountsController = require('../controllers/accountControllers');
const router = express.Router();

router.get('/financial-accounts/ping', accountsController.pingAccount);

module.exports = router;