const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { schnorr, fromBech32Address } = require('@zilliqa-js/crypto');
const api = 'https://dev-api.zilliqa.com';
const zilliqa = new Zilliqa(api);
const privateKey = '{{ SENDER_PRIVATE_KEY }}'

async function sign() {
    try {
        zilliqa.wallet.addByPrivateKey(privateKey);

        const recipientAddr = fromBech32Address("zil1f9uqwhwkq7fnzgh5x4djyzg4a7j3apx8dsnnc0");

        const rawTx = zilliqa.transactions.new({
            version: bytes.pack(333, 1),
            amount: new BN(units.toQa('2', units.Units.Zil)),
            gasLimit: Long.fromNumber(50), // normal (non-contract) transactions cost 1 gas
            gasPrice: new BN(units.toQa(2000, units.Units.Li)), // the minimum gas price is 1,000 li
            toAddr: recipientAddr, // recipient address must be converted to checksum address, 
            pubKey: "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e", // this determines which account is used to send the tx
        });
        

        // note that we provided the nonce to use when constructing the transaction.
        // if the nonce is not provided, zilliqa-js will automatically try to determine the correct nonce to use.
        // however, if there is no network connection, zilliqa-js will not be able to
        // do that, and signing will fail.
        const signedTx = await zilliqa.wallet.sign(rawTx);
        const signature = schnorr.toSignature(signedTx.txParams.signature);
        const lgtm = schnorr.verify(signedTx.bytes, signature, Buffer.from("02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e", 'hex'));

        // console.log("schnorr signature: %o", signature);
        console.log("schnorr verify: %o", lgtm);
        console.log("signed transaction is: %o", JSON.stringify(signedTx));
        console.log("signature is: %o", signedTx.signature);
        console.log("signed transaction in hex_bytes is: %o", signedTx.bytes.toString('hex'));
    } catch (err) {
        console.log(err);
    }
}

sign()