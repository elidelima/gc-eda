const express = require("express");
const mysql = require('mysql');

const app = express();
const port = 3003;

const dbconfig = {
    host: 'balancedb',
    port: 3306,
    user: 'root',
    password: 'root',
    database: 'balance',
};

(async () => {
    console.log('Configuring balance db: ');
    const connection = mysql.createConnection(dbconfig);
    try {
        connection.query('DROP TABLE IF EXISTS balances');
        connection.query(`
            CREATE TABLE balances(
                account_id VARCHAR(255), 
                balance INT, 
                created_at DATE, 

                primary key(account_id)
            );
        `);

        const now = new Date();
        const balances = [
            { account_id: '2161e2cf-27ba-46f2-aab7-950c00dacabf', balance: 1000, created_at: now },
            { account_id: 'b06936c1-89d4-49dc-a480-c7381e25b582', balance: 1000, created_at: now },
        ];
        balances.forEach((balance) => connection.query(`INSERT INTO balances SET ?`, balance));
    } finally {
        connection.end();
    }
})();

app.get('/balances', async (req, res) => {
    let connection;
    try {
        const selectAllBalances = 'SELECT account_id, balance, created_at FROM balances';
        connection = mysql.createConnection(dbconfig);
        connection.query(selectAllBalances, (error, result) => {
            if (error) return res.send('Internal Server Error');
            return res.json(result);
        });

    } catch (error) {
        res.send('Internal Server Error');
    } finally {
        connection?.end();
    }
});

app.get('/balances/:account_id', async (req, res) => {
    let connection;
    try {
        const accountId = req.params?.account_id;
        if (!accountId) return res.status(400).send('400 - Invalid account_id');

        const getBalanceByAccountId = 'SELECT account_id, balance FROM balances WHERE account_id = ?';

        connection = mysql.createConnection(dbconfig);
        connection.query(getBalanceByAccountId, [accountId], (error, result) => {
            if (error) return res.send('Internal Server Error');

            if (!result || result.length === 0) return res.status(404).send('404 - Not found');

            return res.json(result[0]);
        });

    } catch (error) {
        res.send('Internal Server Error');
    } finally {
        connection?.end();
    }
});

app.listen(port, () => {
    console.log(`App rodando na porta ${port}`)
});

// TOPIC CONSUMER
const { Kafka } = require('kafkajs');

const kafka = new Kafka({
    clientId: 'balanceapp',
    brokers: ['kafka:29092'],
});

async function updateBalances(payload) {
    const {
        account_id_from,
        balance_account_id_from,
        account_id_to,
        balance_account_id_to,
    } = payload;
    connection = mysql.createConnection(dbconfig);
    connection.beginTransaction(function (err) {
        if (err) { throw err; }
        //UPDATE balance account_from
        connection.query(
            'UPDATE balances SET balance = ? WHERE account_id = ?',
            [balance_account_id_from, account_id_from],
            function (error, results) {
                if (error) {
                    return connection.rollback(function () {
                        throw error;
                    });
                }

                //UPDATE balance account_to
                connection.query(
                    'UPDATE balances SET balance = ? WHERE account_id = ?',
                    [balance_account_id_to, account_id_to],
                    function (error, results) {
                        if (error) {
                            return connection.rollback(function () {
                                throw error;
                            });
                        }
                        connection.commit(function (err) {
                            if (err) {
                                return connection.rollback(function () {
                                    throw err;
                                });
                            }
                            console.log(`Balances updated successfully!`);
                        });
                    });
            });
    });
}

const subscribe = async () => {
    const consumer = kafka.consumer({ groupId: 'wallet' })
    await consumer.connect()
    await consumer.subscribe({ topic: 'balances', fromBeginning: true })

    await consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
            const prefix = `${topic}[${partition} | ${message.offset}] / ${message.timestamp}`
            const { Payload: payload } = JSON.parse(message.value);
            console.log({ prefix, payload });
            
            await updateBalances(payload);
        },
    });
}

subscribe().catch((e) => console.error(`[wallet/balances] - ${e}`));