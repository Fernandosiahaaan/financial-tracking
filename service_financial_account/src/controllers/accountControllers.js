const pingAccount = async (req, res) => {
    res.status(200).json({status: 'pong'});
}

module.exports = {
    pingAccount,
};