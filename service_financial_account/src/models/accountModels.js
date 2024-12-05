const logger = require('../config/logger');
const { DataTypes } = require('sequelize');
const sequelize = require('../db/postgree');

const FinanceAccount = sequelize.define('finance-account', {
    id: {
        type: DataTypes.INTEGER,
        primaryKey: true,
        autoIncrement: true,
        allowNull: false,
      },
      name: {
        type: DataTypes.STRING(50),
        unique: true,
        allowNull: false,
      },
      description: {
        type: DataTypes.TEXT,
        allowNull: true,
      },
      created_at: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
      },
      updated_at: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
      },
      deleted_at: {
        type: DataTypes.DATE,
        allowNull: true,
        defaultValue: null,
      },
}, {
    tableName: 'finance_accounts',
    timestamps: false, // Menonaktifkan otomatis pencatatan timestamps (karena kita sudah membuatnya sendiri)
});

module.exports= FinanceAccount;