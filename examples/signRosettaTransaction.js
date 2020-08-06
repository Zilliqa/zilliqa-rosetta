const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { schnorr } = require('@zilliqa-js/crypto');
const api = 'https://dev-api.zilliqa.com';
const zilliqa = new Zilliqa(api);
const privateKey = '{{ SENDER_PRIVATE_KEY }}'

async function sign() {
    try {
        zilliqa.wallet.addByPrivateKey(privateKey);

        const rawTx = zilliqa.transactions.new({
            version: bytes.pack(333, 1),
            amount: new BN(units.toQa('2', units.Units.Zil)),
            gasLimit: Long.fromNumber(1), // normal (non-contract) transactions cost 1 gas
            gasPrice: new BN(units.toQa(1000, units.Units.Li)), // the minimum gas price is 1,000 li
            toAddr: "0x4978075dd607933122f4355B220915EFa51E84c7", // toAddr is self-explanatory
            pubKey: "02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e", // this determines which account is used to send the tx
        });
        
        // signWith uses the specified address to perform the signing of the transaction.
        // note that we provided the nonce to use when constructing the transaction.
        // if the nonce is not provided, zilliqa-js will automatically try to determine the correct nonce to use.
        // however, if there is no network connection, zilliqa-js will not be able to
        // do that, and signing will fail.
        // const signedTx = await this.zil.wallet.signWith(rawTx, address);
        const signedTx = await zilliqa.wallet.sign(rawTx);
        const signature = schnorr.toSignature(signedTx.txParams.signature);
        const signedTxHex = rawTx.bytes.toString('hex');
        const lgtm = schnorr.verify(signedTx.bytes, signature, Buffer.from("02e44ef2c5c2031386faa6cafdf5f67318cc661871b0112a27458e65f37a35655e", 'hex'));

        // console.log("schnorr signature: %o", signature);
        console.log("schnorr verify: %o", lgtm);
        console.log("signed transaction is: %o", JSON.stringify(signedTx));
        console.log("signature is: %o", signedTx.signature);
        console.log("signed transaction in hex_bytes is: %o", signedTxHex);
    } catch (err) {
        console.log(err);
    }
}

sign()