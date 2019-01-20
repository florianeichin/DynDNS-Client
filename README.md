# Dynamic DNS Client
This client enables you to update your DynDNS based on DynDNS v2 Protocol, which is available here http://dyn.com/support/developers/api/ 
You could add this to your setup and a daily cronjob to update your DynDNS on a daily basis.

You will need some inofrmation to configure the system.
The Update Url for instance. You get those information from your Domain Registrar. Here are some examples.:

| Registrar|Update Url|
|---|---|
| DNS_O_Matic | "https://update.dnsomatic.com/nic/update" |
| DynDNS | "https://members.dyndns.org/nic/update" |
| No_IP |  "https://dynupdate.no-ip.com/nic/update" |
## System Requirements (just tested one)
- Go
- Linux, Windows
- Arm64, AMD64

## Getting Started
- Get the Source (```git clone https://github.com/florianeichin/DynDNS-Client.git```)
- Get the dependencies (```go get "github.com/jayschwa/go-dyndns" && go get "github.com/Sirupsen/logrus"```)
- Navigate into project (```cd DynDNS-Client```)
- Rename example.config.json to config.json (```mv example.config.json config.json```)
- Edit the config.json to your needs. You get the data from your registrar.
- Build the Client (```go build```)
- Run  (```./DynDNS-Client```)
- You could add this to your daily cron or whatever
