# SF3500 Service
SF3500 service is a FaceID application service that serve incoming enroll and log data from Keico SF3500 device in push mode

# Running SF3500 as a service on Windows

## Install SF3500 Service

[winsw](https://github.com/kohsuke/winsw) is a wrapper to run any executable as an Windows service

- Download [WinSW-x64.exe](https://github.com/winsw/winsw/releases/download/v2.11.0/WinSW-x64.exe)
- Rename the `WinSW-x64.exe` to `sf3500-service.exe`
- Create a xml file `sf3500-service.xml` insert the configuration below
- Open a `cmd` as Administrator and execute `sf3500-service.exe install`

```xml
<service>
  <id>FaceID</id>
  <name>FaceID Service</name>
  <description>FaceID service is a application that serve incoming enroll and log data from Keico SF3500 device in push mode</description>
  <executable>%BASE%\sf3500.exe</executable>
  <!-- change dsn argument below with your database configuration -->
  <arguments>server --port 9009 --dsn "sqlserver://user:password@hostname?database=FaceID"</arguments>
  <logmode>rotate</logmode>
</service>
```

## Manual install with ENVs

[winsw](https://github.com/kohsuke/winsw) is a wrapper to run any executable as an Windows service

- Download [WinSW-x64.exe](https://github.com/winsw/winsw/releases/download/v2.11.0/WinSW-x64.exe)
- Rename the `WinSW-x64.exe` to `sf3500-service.exe`
- Create a xml file `sf3500-service.xml` insert the configuration below
- Open a `cmd` as Administrator and execute `sf3500-service.exe install`

```xml
<service>
  <id>FaceID</id>
  <name>FaceID Service</name>
  <description>FaceID service is a application that serve incoming enroll and log data from Keico SF3500 device in push mode</description>
  <executable>%BASE%\sf3500.exe</executable>
  <env name="SF3500_PORT" value="9009"/>
  <!-- change dsn argument below with your database configuration -->
  <env name="SF3500_DSN" value="sqlserver://user:password@hostname?database=FaceID"/>
  <arguments>server"</arguments>
  <logmode>rotate</logmode>
</service>
```