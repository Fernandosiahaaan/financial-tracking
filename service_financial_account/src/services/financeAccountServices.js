const FinanceAccount = require('../models/accountModels');
const logger = require('../config/logger');

async function CreateFinanceAccount(data) {
    try {
        const { name, description } = data;
        const account = await FinanceAccount.findOne({
            where: {
                name: name,
                deleted_at : null,
            }
        });
        if (account) {
            return {
                error: true,
                message: `account ${name} already exist`
            };
        }

        const newAccount = await FinanceAccount.create({
            name,
            description,
        });
        return {
            error: (!newAccount) ? true : false,
            message: (!newAccount) ? `failed create account finance account ${name}` : `success create finance account ${name}`,
            data: (!newAccount) ? null : newAccount,
        };

    } catch (error) {
        return {
            error: true,
            message: error.message
        };
    }
}

async function GetFinanceAccounts() {
    try {
        const accounts = await FinanceAccount.findAll({
            where: {
                deleted_at : null,
            }
        });
        logger.info(`accounts = ${accounts}`);
        return {
            error: (accounts.length <= 0) ? true : false,
            message: (accounts.length <= 0) ? `failed get finance accounts` : `success get all finance accounts`,
            data: (accounts.length <= 0) ? null : accounts,
        };
    } catch (error) {
        return {
            error: true,
            message: `failed get all finance accounts. err = ${error.message}`
        };
    }
}

async function GetFinanceAccountByID(id) {
    try {
        const account = await FinanceAccount.findByPk(id, {
            where: {
                deleted_at : null,
            }
        });
        return {
            error: (!account) ? true : false,
            message: (!account) ? `failed get finance account with id ${account.id}` : `success get finance account with id ${account.id}`,
            data: (!account) ? null : account,
        };
    } catch (error) {
        return {
            error: true,
            message: error.message
        };
    }
}

async function UpdateFinanceAccountByID(id, data) {
    try {
        const account = await FinanceAccount.findByPk(id);
        if (!account || account.deleted_at !== null) {
            return {
                error: true,
                message: `finance account ${name} not found`
            };
        }

        account.name = data.name;
        account.description = data.description;
        account.updated_at = new Date();
        const accountUpdate = await account.save();
        return {
            error: (accountUpdate) ? false : true,
            message: (accountUpdate) ? `success update finance account ${accountUpdate.name}` : `failed update finance account ${accountUpdate.name}`,
            data: (accountUpdate) ? accountUpdate : null,
        };
    } catch (error) {
        return {
            error: true,
            message: `failed update finance account ${id} . err ${error.message}`
        };
    }
}

async function DeleteFinanceAccountByID(id) {
    try {
        const account = await FinanceAccount.findByPk(id);
        if (!account || account.deleted_at !== null) {
            return {
                error: true,
                message: `account ${name} not found`
            };
        }
        account.deleted_at = new Date();
        const accountDelete = await account.save();
        return {
            error: (accountDelete) ? false : true,
            message: (accountDelete) ? `success delete finance account ${account.name}` : `failed delete finance account ${account.name}`,
            data: (accountDelete) ? accountDelete : null,
        };
    } catch (error) {
        return {
            error: true,
            message: error.message
        };
    }
}

module.exports = {
    CreateFinanceAccount,
    GetFinanceAccounts,
    GetFinanceAccountByID,
    UpdateFinanceAccountByID,
    DeleteFinanceAccountByID,
}