const { error } = require('winston');
const logger = require('../config/logger')
const serviceFinanceAccount = require('../services/financeAccountServices');

const pingAccount = async (req, res) => {
    res.status(200).json({status: 'pong'});
}

const createFinanceAccount = async (req, res) => {
    try {
        const data = req.body;
        dataResponse = await serviceFinanceAccount.CreateFinanceAccount(data);
        if (dataResponse.error) {
            return res.status(500).json(dataResponse)
        }
        res.status(201).json(dataResponse);
    } catch (error) {
        return res.status(500).json({
            error: true,
            message: `failed create finance account ${data.name}. err = ${error.message}`,
        });
    }
};

const getAllFinanceAccount = async (req, res) => {
    try {
        dataResponse = await serviceFinanceAccount.GetFinanceAccounts();
        if (dataResponse.error) {
            return res.status(500).json(dataResponse);
        }
        return res.status(200).json(dataResponse);

    } catch (error) {
        return res.status(500).json({
            error: true,
            message: `failed get finance account. err = ${error.message}`,
        });
    }

};

const getFinanceAccountByID = async (req, res) => {
    try {
        const { id } = req.params;
        dataResponse = await serviceFinanceAccount.GetFinanceAccountByID(id);
        if (dataResponse.error) {
            return res.status(500).json(dataResponse)
        }
        return res.status(200).json(dataResponse);

    } catch (error) {
        return res.status(500).json({
            error: true,
            message: `failed get all finance account ${id}. err = ${error.message}`,
        });
    }
};

const updateFinanceAccountByID = async (req, res) => {
    try {
        const data = req.body;
        const { id } = req.params;
        dataResponse = await serviceFinanceAccount.UpdateFinanceAccountByID(id, data);
        if (dataResponse.error) {
            return res.status(500).json(dataResponse)
        }
        return res.status(200).json(dataResponse);

    } catch (error) {
        return res.status(500).json({
            error: true,
            message: `failed update finance account ${data.name}. err = ${error.message}`,
        });
    }

};

const deleteFinanceAccountByID = async (req, res) => {
    try {
        const { id } = req.params;
        dataResponse = await serviceFinanceAccount.DeleteFinanceAccountByID(id);
        if (dataResponse.error) {
            return res.status(500).json(dataResponse)
        }
        return res.status(200).json(dataResponse);

    } catch (error) {
        return res.status(500).json({
            error: true,
            message: `failed delete finance account ${id}. err = ${error.message}`,
        });
    }
};

module.exports = {
    pingAccount,
    createFinanceAccount,
    getAllFinanceAccount,
    getFinanceAccountByID,
    updateFinanceAccountByID,
    deleteFinanceAccountByID,
};