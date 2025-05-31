const Joi = require('joi');

const accountSchema = joi.object({
    id: Joi.number().integer(),
    name: Joi.string().max(50).required().error(() => 'name must be string'),
    description: Joi.string().allow(null).error(() => 'desctiption must be string'),
    createdAt: Joi.date().default(Date.now).error(() => 'created at must be date'),
    updatedAt: Joi.date().default(Date.now).error(() => 'updated at must be date'),
    deletedAt: Joi.date().allow(null).error(() => 'deleted at must be date'),
});

module.exports = {
    accountSchema,
}