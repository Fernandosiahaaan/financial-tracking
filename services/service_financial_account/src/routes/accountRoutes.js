const express = require('express');
const accountsController = require('../controllers/financeAccountControllers');
const verifyToken = require('../middlewares/jwtMiddleware');
const router = express.Router();

router.get('/financial-accounts/ping', accountsController.pingAccount);
router.post('/financial-accounts', verifyToken, accountsController.createFinanceAccount);
router.get('/financial-accounts', verifyToken, accountsController.getAllFinanceAccount);
router.get('/financial-accounts/:id', verifyToken, accountsController.getFinanceAccountByID);
router.put('/financial-accounts/:id', verifyToken, accountsController.updateFinanceAccountByID);
router.delete('/financial-accounts/:id', verifyToken, accountsController.deleteFinanceAccountByID);


module.exports = router;