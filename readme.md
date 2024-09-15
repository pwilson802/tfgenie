# Terraform Genie

tfgenie is a curent work in progress.  It will help save time writting terraform code in use cases when the code may takea  while to get correct and you have a refernce system. 

Creating Grafana alerts and dashboards can take some playing around in the GUI, it can be nice to create them like this first the use Terraform Genie to pull them into a terraform file. 

IAM policies following least privilages principal can take some time to get correct. Terraform Genie can connect to IAM access analyzer and create terraform code for a role from a recent role or users usage.

## Getting Started

tfgenie if written in Go, this will need to be installed first, you can clone the project and build it yourself or install direct from github.


### Installing

Go is required to install tfgenie

```
go get https://github.com/pwilson802/tfgenie
```

### Usage

Example of getting a grafna alert:
```
export GRAFANA_API_KEY=XXXXXXX
tfgenie grafana --hostname grafana-server.dev --resource alert --alertId aba370f4-ba77-4de6-93f1-4a32158cb2eb
```


## Authors

* **Paul Wilson** - *Initial work* - [pwilson802](https://github.com/pwilson802)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
