const { Pool } = require('pg');
const logger = require('../../config/logger');  

const pool = new Pool({
    connectionString: process.env.POSTGREE_URL,
});

// Event listener ketika koneksi berhasil
pool.on('connect', () => {
    logger.info('Berhasil terhubung ke database!');
});

// Event listener ketika terjadi error dalam koneksi
pool.on('error', (error) => {
    logger.error(`Terjadi kesalahan saat terhubung ke database: ${error.stack}`);
});


module.exports = {
    pool,
};
