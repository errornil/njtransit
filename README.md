# njtransit

The library that provides an interface to NJTransit public Bus API.

## [Register as 3rd Party Developer](https://datasource.njtransit.com/SignUp.aspx)

Register on NJTransit website to get username and password.

## Usage example

```
client := njt.NewBusDataClient(
    os.Getenv("BUSDATA_USERNAME"),
    os.Getenv("BUSDATA_PASSWORD"),
    njt.BusDataProdURL,
)

log.Println("Calling GetBusVehicleData...")
resp, err := client.GetBusVehicleData()
if err != nil {
    log.Fatalf("Failed to call GetBusDV: %v", err)
}
log.Printf("%#v", *resp)
```

## Limits

> User access limits are set to allow reasonable daily usage. These limits will not allow further accesses to the web service for the remainder of the day. After midnight these will be reset to zero. There will be a 40,000 limit per day for the current and vehicle data and 10 accesses per day for the full schedule. (Users needing more than 40,000 per day must demonstrate their user base of more than this limit and provide monthly user reports.)

## Legal

If you are using NJTransit data, you have to provide (and comply with) this disclaimer:

> Data provided by NJ TRANSIT, which is the sole owner of the Data. This “App” is not endorsed by, directly affiliated with, maintained, authorized, or sponsored by NJ TRANSIT. All product and company names are the registered trademarks of their original owners. The use of any trade name or trademark is for identification and reference purposes only and does not imply any association with the trademark owner.

*By the way*, this repository is not endorsed by, directly affiliated with, maintained, authorized, or sponsored by NJ TRANSIT. All product and company names are the registered trademarks of their original owners. The use of any trade name or trademark is for identification and reference purposes only and does not imply any association with the trademark owner.
