wget https://github.com/prysmaticlabs/prysm/releases/download/v1.0.0-alpha.29/validator-v1.0.0-alpha.29-darwin-amd64
mv validator-v1.0.0-alpha.29-darwin-amd64 validator_bin
chmod +x validator_bin

# import
./validator_bin accounts-v2 import --keys-dir="./validator/validator_keys" --wallet-dir="./validator/wallet" --wallet-password-file="./validator/passwords/wallet-password" --medalla

# list
./validator_bin accounts-v2 list --wallet-dir="./validator/wallet" --wallet-password-file="./validator/passwords/wallet-password" --show-deposit-data

rm validator_bin
