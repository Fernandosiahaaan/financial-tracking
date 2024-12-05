const express = require('express');
const accountsController = require('../controllers/financeAccountControllers');
const router = express.Router();

router.get('/financial-accounts/ping', accountsController.pingAccount);
router.post('/financial-accounts', accountsController.createFinanceAccount);
router.get('/financial-accounts', accountsController.getAllFinanceAccount);
router.get('/financial-accounts/:id', accountsController.getFinanceAccountByID);
router.put('/financial-accounts/:id', accountsController.updateFinanceAccountByID);
router.delete('/financial-accounts/:id', accountsController.deleteFinanceAccountByID);

module.exports = router;