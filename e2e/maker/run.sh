# Run market maker bot locally. Usage: `cd e2e/maker && ./run.sh`
# First run `go run scripts/redis-repo/main.go createUser --apiKey og4lpqQUILyciacspkFESHE1qrXIxpX1` to create MM user

deactivate || echo "no venv to deactivate"
rm -rf .venv
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
# Dummy values for testing - do not use in production
export PRIVATE_KEY="0xb3f663326aae5ddd1a368f3485949362abcc54440ceb3f3f44242cfe1c08219d"
export API_KEY="og4lpqQUILyciacspkFESHE1qrXIxpX1"
python main.py