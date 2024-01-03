# Currency Rates Fetcher

This project fetches currency rates from various APIs and stores them in a JSON file. It uses both a Go script and a Bash script to fetch and process the data.

## Project Structure

- `.github/workflows/`: Contains GitHub Actions workflows for running the scripts and committing the results.
- `currency_all/`: Directory where the JSON file with the currency rates is stored.
- `src/`: Contains the source code for the scripts.

## Scripts

- `bash-script-bi.sh`: A Bash script that fetches data from the Banco Industrial API and stores it in a JSON file.
- `main.go`: A Go script that fetches data from various currency APIs and stores the results in a JSON file.

## GitHub Actions

There are two workflows defined in this project:

- `bash.yml`: This workflow runs the Bash script every 6 hours and commits the resulting JSON file to the repository.
- `main.yml`: This workflow runs the Go script every 12 hours and commits any changes to the repository.

## Running the Scripts

To run the Bash script, use the following command:

```sh
bash ./src/bash-script-bi.sh > ./currency_all/tipo_de_cambio_bi.json
```

To run the Go script, use the following command:

```go
go run ./src/main.go
```

Please note that you need to have `bash`, `curl`, `jq`, and `go` installed on your system to run these scripts.

## Contributing
Contributions are welcome. Please open an issue or submit a pull request if you have any improvements or features to suggest. 