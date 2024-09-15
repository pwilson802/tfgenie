# Terraform Genie

tfgenie is a work in progress. It will help save time writing Terraform code in cases where the code may take a while to get correct, and you have a reference system. 

Creating Grafana alerts and dashboards can take some playing around in the GUI. It can be nice to create them like this first and then use Terraform Genie to pull them into a Terraform file. 

IAM policies following the least privileges principle can take some time to get correct. Terraform Genie can connect to IAM access analyzer and create Terraform code for a role based on a recent role or user's usage.

## Getting Started

As tfgenie is written in Go, it will need to be installed first. You can clone the project and build it yourself or install it directly from GitHub.


### Installing

Go is required to install tfgenie

```
go get https://github.com/pwilson802/tfgenie
```

### Usage

Example of getting a grafna alert:
```
export GRAFANA_API_KEY=XXXXXX
tfgenie grafana --hostname grafana-server.dev --resource alert --alertId aba370f4-ba77-4de6-93f1-4a32158cb2eb
```


## Authors

* **Paul Wilson** - *Initial work* - [pwilson802](https://github.com/pwilson802)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
