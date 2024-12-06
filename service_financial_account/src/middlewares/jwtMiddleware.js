const jwt = require('jsonwebtoken');
const dotenv = require('dotenv');

function verifyToken(req, res, next) {
    const token = req.header('Authorization')?.replace('Bearer ', '');
    if (!token) {
        return res.status(401).json({
            error: true,
            message: 'Access Denied. No token provided.',
        });
    }

    try {
        // Verifikasi token menggunakan JWT_SECRET_KEY yang ada di .env
        jwt.verify(token, process.env.SECRET_KEY, (err, user) => {
          if (err) {
            return res.status(403).json({
              error: true,
              message: 'Invalid or expired token.',
            });
          }
          // Menyimpan informasi user dari token pada req.user untuk digunakan oleh route lainnya
          req.user = user;
          next();  // Melanjutkan ke middleware atau handler berikutnya
        });
      } catch (err) {
        return res.status(500).json({
          error: true,
          message: 'Failed to authenticate token.',
        });
    }
}

module.exports = verifyToken;